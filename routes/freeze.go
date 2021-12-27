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

func (srv *Server) freeze(footer string) {
	srv.router.GET("/admin/freeze.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)

		type freezeQuery struct {
			Shelf   string `form:"shelf"`
			Crate   string `form:"crate"`
			Booklet string `form:"booklet"`
			Freeze  bool   `form:"freeze"`
		}
		var q freezeQuery
		err := c.ShouldBind(&q)
		if err != nil {
			log.Error(err)
			srv.freezeWithError(
				c,
				header,
				footer,
				fmt.Sprintf("unparsable request : %s", err.Error()),
				"",
				"",
				"")
			return
		}
		if q.Booklet == "" {
			srv.freezeOkWithData(c, header, footer, "", "", "")
		} else {
			srv.freezeBooklet(c, header, footer, q.Shelf, q.Crate, q.Booklet, q.Freeze)
		}
	})

	srv.router.POST("/admin/freeze.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		shelfNumber := c.PostForm("ShelfNumber")
		crateNumber := c.PostForm("CrateNumber")
		bookletNumber := c.PostForm("BookletNumber")
		srv.freezeOkWithData(c, header, footer, shelfNumber, crateNumber, bookletNumber)

	})
}

func (srv *Server) freezeBooklet(c *gin.Context, header, footer, shelfNumber, crateNumber, bookletNumber string, freeze bool) {
	id, err := ginoidc.GetIdentity(c)
	if err != nil {
		log.Error(err)
		srv.freezeWithError(
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
		srv.freezeWithError(
			c,
			header,
			footer,
			"Unknown booklet, please check number or booklet registration",
			shelfNumber,
			crateNumber,
			bookletNumber)
		return
	}

	if freeze {
		err = booklet.FreezeBooklet(crateNumber, shelfNumber, id, &srv.Db)
	}

	if err != nil {
		log.Error(err)
		srv.freezeWithError(
			c,
			header,
			footer,
			err.Error(),
			shelfNumber,
			crateNumber,
			bookletNumber)
		return
	}
	srv.freezeOkWithData(c, header, footer, shelfNumber, crateNumber, "")
}

func (srv *Server) freezeWithError(c *gin.Context, header, footer, alertMsg string, shelfNumber, crateNumber, bookletNumber string) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := ginoidc.GetValue(c, "name")
	c.HTML(http.StatusOK, "freeze.html", gin.H{
		"Name":          name,
		"Header":        template.HTML(header),
		"Footer":        template.HTML(footer),
		"Alert":         template.HTML(alert),
		"ShelfNumber":   shelfNumber,
		"CrateNumber":   crateNumber,
		"BookletNumber": bookletNumber,
	})
}

func (srv *Server) freezeOkWithData(c *gin.Context, header, footer, shelfNumber, crateNumber, bookletNumber string) {
	name := ginoidc.GetValue(c, "name")

	tableData, lastCrateNumber, lastShelfNumber, lastBookletNumber := srv.getFreezableBooklet(shelfNumber, crateNumber, bookletNumber)
	stat, err := srv.getFreezeStats(lastCrateNumber, lastShelfNumber)
	if err != nil {
		log.Error(err)
	}

	c.HTML(http.StatusOK, "freeze.html", gin.H{
		"Name":                       name,
		"Header":                     template.HTML(header),
		"Footer":                     template.HTML(footer),
		"ShelfNumber":                lastShelfNumber,
		"CrateNumber":                lastCrateNumber,
		"BookletNumber":              lastBookletNumber,
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

func (srv *Server) getFreezableBooklet(shelfNumber, crateNumber, bookletNumber string) (data, lastCrateNumber, lastShelfNumber, lastBookletNumber string) {
	var booklets []simpleBooklets
	if shelfNumber != "" {
		switch {
		case crateNumber != "" && bookletNumber != "":
			_, err := srv.Db.Conn.Query(&booklets, `
	select
       s.number as shelf_number,
       c.number as crate_number,
       b.number as booklet_number,
       b.size as booklet_size,
	   b.status as booklet_status,
	   b.cut_by as cut_operator		
from booklets b,
     crates c,
     shelves s
where
      b.crate_number = c.number
  and c.shelf_number = s.number
  and (b.status = 'inCuttingStation' or b.status = 'inIceBox') 				
  and c.number = ?
  and s.number = ?
  and b.number = ?;`, crateNumber, shelfNumber, bookletNumber)
			if err != nil {
				log.Error(err)
				return
			}
		case crateNumber != "":
			_, err := srv.Db.Conn.Query(&booklets, `
	select
       s.number as shelf_number,
       c.number as crate_number,
       b.number as booklet_number,
       b.size as booklet_size,
	   b.status as booklet_status,
	   b.cut_by as cut_operator		
from booklets b,
     crates c,
     shelves s
where
      b.crate_number = c.number
  and c.shelf_number = s.number
  and (b.status = 'inCuttingStation' or b.status = 'inIceBox') 				
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
	   b.status as booklet_status,
	   b.cut_by as cut_operator	
from booklets b,
     crates c,
     shelves s
where
      b.crate_number = c.number
  and c.shelf_number = s.number
  and (b.status = 'inCuttingStation' or b.status = 'inIceBox') 			
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
			<td>%s</td>
			<td>
				<a href="#" data-href="/admin/freeze.html?shelf=%s&crate=%s&booklet=%s&freeze=true" data-toggle="modal" data-target="#confirm-freeze">
					<center>
						<i class="fa fa-archive"></i>
					</center>
				</a>
			</td>
		</tr>`, booklet.ShelfNumber,
				booklet.CrateNumber,
				booklet.BookletNumber,
				booklet.BookletSize,
				booklet.BookletStatus,
				booklet.CutOperator,
				booklet.ShelfNumber,
				booklet.CrateNumber,
				booklet.BookletNumber)
		}
	}

	if len(booklets) == 1 {
		lastCrateNumber = booklets[0].CrateNumber
		lastShelfNumber = booklets[0].ShelfNumber
		lastBookletNumber = booklets[0].BookletNumber
	} else if len(booklets) > 1 {
		lastCrateNumber = booklets[0].CrateNumber
		lastShelfNumber = booklets[0].ShelfNumber
	}
	return
}

func (srv *Server) getFreezeStats(crateNumber, shelfNumber string) (stats cuttingStat, err error) {
	stats = cuttingStat{}
	_, err = srv.Db.Conn.QueryOne(&stats, `
	select *
	from (select count(distinct (b)) as booklet_in_crate, sum(b.size) as sheet_in_crate
      from booklets b,
           crates c
      where b.crate_number = c.number
        and c.number = ?
		and b.status = 'inCuttingStation') as crate
        ,
     (select count(distinct (c.number)) as crate_in_shelf, sum(b.size) as sheet_in_shelf
      from booklets b,
           crates c,
           shelves s
      where b.crate_number = c.number
        and c.shelf_number = s.number
        and s.number = ?
		and b.status = 'inCuttingStation') as shelf;`, crateNumber, shelfNumber)
	if err != nil {
		return
	}
	return
}
