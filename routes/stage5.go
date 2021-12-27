package routes

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	gin_oidc "git.solutions.im/Solutions.IM/ginOidc"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

type scanningStat struct {
	BookletInCrate int
	SheetInCrate   int
	CrateInShelf   int
	SheetInShelf   int
}

func (srv *Server) stage5(footer string) {
	srv.router.GET("/production/stage5.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		srv.stage5OkWithData(c, header, footer, "", "")
	})

	srv.router.POST("/production/stage5.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		shelfNumber := c.PostForm("ShelfNumber")
		crateNumber := c.PostForm("CrateNumber")

		switch {
		case shelfNumber != "" && crateNumber != "":
			srv.stage5Scanning(c, header, footer, shelfNumber, crateNumber)
		default:
			srv.stage5OkWithData(c, header, footer, shelfNumber, crateNumber)
		}
	})
}

func (srv *Server) stage5Scanning(c *gin.Context, header, footer, shelfNumber, crateNumber string) {
	id, err := gin_oidc.GetIdentity(c)
	if err != nil {
		log.Error(err)
		srv.stage5WithError(
			c,
			header,
			footer,
			err.Error(),
			shelfNumber,
			crateNumber)
		return
	}

	// get if crate exist
	crate, err := srv.Db.GetCrate(crateNumber)
	if err != nil {
		log.Error(err)
		srv.stage5WithError(
			c,
			header,
			footer,
			"Unknown crate, please check number or crate registration",
			shelfNumber,
			crateNumber)
		return
	}

	err = crate.Scann(shelfNumber, id, &srv.Db)
	if err != nil {
		log.Error(err)
		srv.stage5WithError(
			c,
			header,
			footer,
			err.Error(),
			shelfNumber,
			crateNumber)
		return
	}
	srv.stage5OkWithData(c, header, footer, shelfNumber, crateNumber)
}

func (srv *Server) stage5WithError(c *gin.Context, header, footer, alertMsg string, shelfNumber, crateNumber string) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := gin_oidc.GetValue(c, "name")
	c.HTML(http.StatusOK, "stage5.html", gin.H{
		"Name":        name,
		"Header":      template.HTML(header),
		"Footer":      template.HTML(footer),
		"Alert":       template.HTML(alert),
		"ShelfNumber": shelfNumber,
		"CrateNumber": crateNumber,
	})
}

func (srv *Server) stage5OkWithData(c *gin.Context, header, footer, shelfNumber, crateNumber string) {
	name := gin_oidc.GetValue(c, "name")

	tableData, lastCrateNumber, lastShelfNumber := srv.getScanningTask(shelfNumber, crateNumber)
	stat, err := srv.getScanningStats(lastCrateNumber, lastShelfNumber)
	if err != nil {
		log.Error(err)
	}

	c.HTML(http.StatusOK, "stage5.html", gin.H{
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

func (srv *Server) getScanningTask(shelfNumber, crateNumber string) (data, lastCrateNumber, lastShelfNumber string) {
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
			case booklet.BookletStatus == "inScanningStation":
				data += "<tr class=\"green\">"
			case booklet.BookletStatus == "inPreScanning":
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

func (srv *Server) getScanningStats(crateNumber, shelfNumber string) (stats scanningStat, err error) {
	stats = scanningStat{}
	_, err = srv.Db.Conn.QueryOne(&stats, `
	select *
	from (select count(distinct (b)) as booklet_in_crate, sum(b.size) as sheet_in_crate
      from booklets b,
           crates c
      where b.crate_number = c.number
        and c.number = ?
		and b.status != 'inScanningStation') as crate
        ,
     (select count(distinct (c.number)) as crate_in_shelf, sum(b.size) as sheet_in_shelf
      from booklets b,
           crates c,
           shelves s
      where b.crate_number = c.number
        and c.shelf_number = s.number
        and s.number = ?
		and b.status != 'inScanningStation') as shelf;`, crateNumber, shelfNumber)
	if err != nil {
		return
	}
	return
}
