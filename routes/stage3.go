package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	gin_oidc "git.solutions.im/Solutions.IM/ginOidc"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

type simpleBooklets struct {
	ShelfNumber      string
	CrateNumber      string
	BookletNumber    string
	BookletSize      int
	BookletStatus    string
	ArchiveBoxNumber string
	ArchivedOn       time.Time
	CutOperator      string
}

type cuttingStat struct {
	BookletInCrate int
	SheetInCrate   int
	CrateInShelf   int
	SheetInShelf   int
}

func (srv *Server) stage3(footer string) {
	srv.router.GET("/production/stage3.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		srv.stage3OkWithData(c, header, footer, "", "")
	})

	srv.router.POST("/production/stage3.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		shelfNumber := c.PostForm("ShelfNumber")
		crateNumber := c.PostForm("CrateNumber")
		bookletNumber := c.PostForm("BookletNumber")

		switch {
		case shelfNumber != "" && crateNumber != "" && bookletNumber != "":
			srv.stage3BookletCut(c, header, footer, shelfNumber, crateNumber, bookletNumber)
		default:
			srv.stage3OkWithData(c, header, footer, shelfNumber, crateNumber)
		}
	})
}

func (srv *Server) stage3BookletCut(c *gin.Context, header, footer, shelfNumber, crateNumber, bookletNumber string) {
	id, err := gin_oidc.GetIdentity(c)
	if err != nil {
		log.Error(err)
		srv.stage3WithError(
			c,
			header,
			footer,
			"Who are you ???",
			shelfNumber,
			crateNumber,
			bookletNumber)
		return
	}

	// get if booklet exist
	booklet, err := srv.Db.GetBooklet(bookletNumber)
	if err != nil {
		log.Error(err)
		srv.stage3WithError(
			c,
			header,
			footer,
			"Unknown booklet, please check number or booklet registration",
			shelfNumber,
			crateNumber,
			bookletNumber)
		return
	}

	err = booklet.CutBooklet(crateNumber, shelfNumber, id, &srv.Db)
	if err != nil {
		log.Error(err)
		srv.stage3WithError(
			c,
			header,
			footer,
			err.Error(),
			shelfNumber,
			crateNumber,
			bookletNumber)
		return
	}
	srv.stage3OkWithData(c, header, footer, shelfNumber, crateNumber)
}

func (srv *Server) stage3WithError(c *gin.Context, header, footer, alertMsg string, shelfNumber, crateNumber, bookletNumber string) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := gin_oidc.GetValue(c, "name")
	c.HTML(http.StatusOK, "stage3.html", gin.H{
		"Name":          name,
		"Header":        template.HTML(header),
		"Footer":        template.HTML(footer),
		"Alert":         template.HTML(alert),
		"ShelfNumber":   shelfNumber,
		"CrateNumber":   crateNumber,
		"BookletNumber": bookletNumber,
	})
}

func (srv *Server) stage3OkWithData(c *gin.Context, header, footer, shelfNumber, crateNumber string) {
	name := gin_oidc.GetValue(c, "name")

	tableData, lastCrateNumber, lastShelfNumber := srv.getCuttingTask(shelfNumber, crateNumber)
	stat, err := srv.getCuttingStats(lastCrateNumber, lastShelfNumber)
	if err != nil {
		log.Error(err)
	}

	c.HTML(http.StatusOK, "stage3.html", gin.H{
		"Name":                       name,
		"Header":                     template.HTML(header),
		"Footer":                     template.HTML(footer),
		"ShelfNumber":                lastShelfNumber,
		"CrateNumber":                lastCrateNumber,
		"TableData":                  template.HTML(tableData),
		"NumberOfCrateInShelf":       stat.CrateInShelf,
		"PercentageOfCrateInShelf":   float64(stat.CrateInShelf) / 21 * 100,
		"NumberOfBookletInCrate":     stat.BookletInCrate,
		"PercentageOfBookletInCrate": float64(stat.BookletInCrate) / 10 * 100,
		"NumberOfSheetInShelf":       stat.SheetInShelf,
		"PercentageOfSheetInShelf":   float64(stat.SheetInShelf) / 21_000 * 100,
		"NumberOfSheetInCrate":       stat.SheetInCrate,
		"PercentageOfSheetInCrate":   float64(stat.SheetInCrate) / 1_000 * 100,
	})
}

func (srv *Server) getCuttingTask(shelfNumber, crateNumber string) (data, lastCrateNumber, lastShelfNumber string) {
	var booklets []simpleBooklets
	if shelfNumber != "" {
		switch {
		case crateNumber != "":
			_, err := srv.Db.Conn.Query(&booklets, `
	select
       s.number as shelf_number,
       c.number as crate_number,
       b.number as booklet_number,
       b.size as booklet_size,
		b.status as booklet_status
from booklets b,
     crates c,
     shelves s
where
      b.crate_number = c.number
  and c.shelf_number = s.number
  and c.number = ?
  and s.number = ?;`, crateNumber, shelfNumber)
			if err != nil {
				log.Error(err)
				return
			}
		case crateNumber == "":
			_, err := srv.Db.Conn.Query(&booklets, `
	select
       s.number as shelf_number,
       c.number as crate_number,
       b.number as booklet_number,
       b.size as booklet_size,
		b.status as booklet_status
from booklets b,
     crates c,
     shelves s
where
      b.crate_number = c.number
  and c.shelf_number = s.number
  and s.number = ?;`, shelfNumber)
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
			<td>%s</td>
			<td>%s</td>
			<td>%d</td>
			<td>%s</td>
		</tr>
	`, booklet.ShelfNumber,
				booklet.CrateNumber,
				booklet.BookletNumber,
				booklet.BookletSize,
				booklet.BookletStatus,
			)
		}
	}

	if len(booklets) >= 1 {
		lastCrateNumber = booklets[0].CrateNumber
		lastShelfNumber = booklets[0].ShelfNumber
	}
	return
}

func (srv *Server) getCuttingStats(crateNumber, shelfNumber string) (stats cuttingStat, err error) {
	stats = cuttingStat{}
	_, err = srv.Db.Conn.QueryOne(&stats, `
	select *
	from (select count(distinct (b)) as booklet_in_crate, sum(b.size) as sheet_in_crate
      from booklets b,
           crates c
      where b.crate_number = c.number
        and c.number = ?
		and b.status != 'inCuttingStation') as crate
        ,
     (select count(distinct (c.number)) as crate_in_shelf, sum(b.size) as sheet_in_shelf
      from booklets b,
           crates c,
           shelves s
      where b.crate_number = c.number
        and c.shelf_number = s.number
        and s.number = ?
		and b.status != 'inCuttingStation') as shelf;`, crateNumber, shelfNumber)
	if err != nil {
		return
	}
	return
}
