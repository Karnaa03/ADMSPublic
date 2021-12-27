package routes

import (
	"fmt"
	"html/template"
	"net/http"

	ginoidc "git.solutions.im/Solutions.IM/ginOidc"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

func (srv *Server) archiveFrozen(footer string) {
	srv.router.GET("/admin/archiveFrozen.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		srv.archiveOkWithData(c, header, footer, "", "")

	})

	srv.router.POST("/admin/archiveFrozen.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		bookletNumber := c.PostForm("BookletNumber")
		archiveBoxNumber := c.PostForm("ArchiveBoxNumber")
		if bookletNumber != "" && archiveBoxNumber != "" {
			srv.archiveBooklet(c, header, footer, bookletNumber, archiveBoxNumber)
		} else {
			srv.archiveOkWithData(c, header, footer, bookletNumber, archiveBoxNumber)
		}
	})
}

func (srv *Server) archiveOkWithData(c *gin.Context, header, footer, bookletNumber, archiveBoxNumber string) {
	name := ginoidc.GetValue(c, "name")

	tableData, lastBookletNumber := srv.getFrozenBooklet(bookletNumber)

	c.HTML(http.StatusOK, "archiveFrozen.html", gin.H{
		"Name":          name,
		"Header":        template.HTML(header),
		"Footer":        template.HTML(footer),
		"BookletNumber": lastBookletNumber,
		"TableData":     template.HTML(tableData),
	})
}

func (srv *Server) archiveWithError(c *gin.Context, header, footer, alertMsg string, bookletNumber, archiveBoxNumber string) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := ginoidc.GetValue(c, "name")
	c.HTML(http.StatusOK, "archiveFrozen.html", gin.H{
		"Name":             name,
		"Header":           template.HTML(header),
		"Footer":           template.HTML(footer),
		"Alert":            template.HTML(alert),
		"BookletNumber":    bookletNumber,
		"ArchiveBoxNumber": archiveBoxNumber,
	})
}

func (srv *Server) getFrozenBooklet(bookletNumber string) (data, lastBookletNumber string) {
	var booklets []simpleBooklets

	if bookletNumber != "" {
		_, err := srv.Db.Conn.Query(&booklets, `
	select b.number as booklet_number,
       	b.size   as booklet_size,
       	b.status as booklet_status,
		b.cut_by as cut_operator
from booklets b
where b.status = 'inIceBox' and b.number = ?;`, bookletNumber)
		if err != nil {
			log.Error(err)
			return
		}
	} else {
		_, err := srv.Db.Conn.Query(&booklets, `
	select b.number as booklet_number,
    	b.size   as booklet_size,
       	b.status as booklet_status,
		b.cut_by as cut_operator
from booklets b
where b.status = 'inIceBox';`)
		if err != nil {
			log.Error(err)
			return
		}
	}

	for _, booklet := range booklets {
		switch {
		case booklet.BookletStatus == "inCuttingStation":
			data += "<tr class=\"green\">"
		case booklet.BookletStatus == "inBatch":
			data += "<tr>"
		default:
			data += "<tr class=\"aero\">"
		}
		data += fmt.Sprintf(`
			<td>%s</td>
			<td>%d</td>
			<td>%s</td>
			<td>%s</td>
		</tr>`,
			booklet.BookletNumber,
			booklet.BookletSize,
			booklet.BookletStatus,
			booklet.CutOperator)
	}

	if len(booklets) == 1 {
		lastBookletNumber = booklets[0].BookletNumber
	}
	return
}

func (srv *Server) archiveBooklet(c *gin.Context, header, footer, bookletNumber, archiveBoxNumber string) {
	id, err := ginoidc.GetIdentity(c)
	if err != nil {
		log.Error(err)
		srv.archiveWithError(
			c,
			header,
			footer,
			"Who are you ???",
			bookletNumber,
			archiveBoxNumber)
		return
	}

	// get if booklet exist
	booklet, err := srv.Db.GetBooklet(bookletNumber)
	if err != nil {
		log.Error(err)
		srv.archiveWithError(
			c,
			header,
			footer,
			"Unknown booklet, please check number or booklet registration",
			bookletNumber,
			archiveBoxNumber)
		return
	}

	err = booklet.Archive(id, "", "", archiveBoxNumber, &srv.Db)
	if err != nil {
		log.Error(err)
		srv.archiveWithError(
			c,
			header,
			footer,
			err.Error(),
			bookletNumber,
			archiveBoxNumber)
		return
	}
	srv.archiveOkWithData(c, header, footer, "", "")
}
