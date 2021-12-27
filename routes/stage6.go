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

type archivingStat struct {
	BookletInCrate              int
	SheetInCrate                int
	CrateInShelf                int
	SheetInShelf                int
	NumberOfBookletInArchiveBox int
	NumberOfSheetInArchiveBox   int
}

func (srv *Server) stage6(footer string) {
	srv.router.GET("/production/stage6.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		srv.stage6OkWithData(c, header, footer, "", "", "")
	})

	srv.router.POST("/production/stage6.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		shelfNumber := c.PostForm("ShelfNumber")
		crateNumber := c.PostForm("CrateNumber")
		bookletNumber := c.PostForm("BookletNumber")
		archiveBoxNumber := c.PostForm("ArchiveBoxNumber")

		switch {
		case shelfNumber != "" && crateNumber != "" && bookletNumber != "":
			srv.stage6Archiving(c, header, footer, shelfNumber, crateNumber, bookletNumber, archiveBoxNumber)
		default:
			srv.stage6OkWithData(c, header, footer, shelfNumber, crateNumber, archiveBoxNumber)
		}
	})
}

func (srv *Server) stage6Archiving(c *gin.Context, header, footer, shelfNumber, crateNumber, bookletNumber, archiveBoxNumber string) {
	id, err := gin_oidc.GetIdentity(c)
	if err != nil {
		log.Error(err)
		srv.stage6WithError(
			c,
			header,
			footer,
			err.Error(),
			shelfNumber,
			crateNumber,
			bookletNumber)
		return
	}

	// get if booklet exist
	booklet, err := srv.Db.GetBooklet(bookletNumber)
	if err != nil {
		log.Error(err)
		srv.stage6WithError(
			c,
			header,
			footer,
			"Unknown booklet, please check number or booklet registration",
			shelfNumber,
			crateNumber,
			bookletNumber)
		return
	}

	err = booklet.Archive(id, crateNumber, shelfNumber, archiveBoxNumber, &srv.Db)
	if err != nil {
		log.Error(err)
		srv.stage6WithError(
			c,
			header,
			footer,
			err.Error(),
			shelfNumber,
			crateNumber,
			bookletNumber)
		return
	}
	srv.stage6OkWithData(c, header, footer, shelfNumber, crateNumber, archiveBoxNumber)
}

func (srv *Server) stage6WithError(c *gin.Context, header, footer, alertMsg string, shelfNumber, crateNumber, bookletNumber string) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := gin_oidc.GetValue(c, "name")
	c.HTML(http.StatusOK, "stage6.html", gin.H{
		"Name":          name,
		"Header":        template.HTML(header),
		"Footer":        template.HTML(footer),
		"Alert":         template.HTML(alert),
		"ShelfNumber":   shelfNumber,
		"CrateNumber":   crateNumber,
		"BookletNumber": bookletNumber,
	})
}

func (srv *Server) stage6OkWithData(c *gin.Context, header, footer, shelfNumber, crateNumber, archiveBoxNumber string) {
	name := gin_oidc.GetValue(c, "name")

	tableData, lastCrateNumber, lastShelfNumber := srv.getArchivingTask(shelfNumber, crateNumber, archiveBoxNumber)
	stat, err := srv.getArchivingStats(lastCrateNumber, lastShelfNumber, archiveBoxNumber)
	if err != nil {
		log.Error(err)
	}

	c.HTML(http.StatusOK, "stage6.html", gin.H{
		"Name":                        name,
		"Header":                      template.HTML(header),
		"Footer":                      template.HTML(footer),
		"ShelfNumber":                 lastShelfNumber,
		"CrateNumber":                 lastCrateNumber,
		"ArchiveBoxNumber":            archiveBoxNumber,
		"TableData":                   template.HTML(tableData),
		"NumberOfCrateInShelf":        stat.CrateInShelf,
		"PercentageOfCrateInShelf":    float64(stat.CrateInShelf) / 21 * 100,
		"NumberOfBookletInCrate":      stat.BookletInCrate,
		"PercentageOfBookletInCrate":  float64(stat.BookletInCrate) / 10 * 100,
		"NumberOfSheetInShelf":        stat.SheetInShelf,
		"PercentageOfSheetInShelf":    float64(stat.SheetInShelf) / 21_000 * 100,
		"NumberOfSheetInCrate":        stat.SheetInCrate,
		"PercentageOfSheetInCrate":    float64(stat.SheetInCrate) / 1_000 * 100,
		"NumberOfBookletInArchiveBox": stat.NumberOfBookletInArchiveBox,
		"NumberOfSheetInArchiveBox":   stat.NumberOfSheetInArchiveBox,
	})
}

func (srv *Server) getArchivingTask(shelfNumber, crateNumber, archiveBoxNumber string) (data, lastCrateNumber, lastShelfNumber string) {
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

	}
	var archivedBooklets []simpleBooklets

	if archiveBoxNumber != "" {
		_, err := srv.Db.Conn.Query(&archivedBooklets, `
select b.number             as booklet_number,
       b.size               as booklet_size,
       b.status             as booklet_status,
       b.archive_box_number as archive_box_number,
	   b.archived_on as archived_on	
from booklets b
where b.archive_box_number = ?;
		`, archiveBoxNumber)
		if err != nil {
			log.Error(err)
			return
		}
	}

	booklets = append(booklets, archivedBooklets...)
	for _, booklet := range booklets {
		switch {
		case booklet.BookletStatus == "archived":
			data += "<tr class=\"green\">"
		case booklet.BookletStatus == "inScanningStation":
			data += "<tr>"
		default:
			data += "<tr class=\"aero\">"
		}
		var archiveDate string
		if !booklet.ArchivedOn.IsZero() {
			archiveDate = booklet.ArchivedOn.String()
		}
		data += fmt.Sprintf(`
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%d</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		</tr>
	`, booklet.ShelfNumber,
			booklet.CrateNumber,
			booklet.BookletNumber,
			booklet.BookletSize,
			booklet.BookletStatus,
			booklet.ArchiveBoxNumber,
			archiveDate,
		)
	}

	if len(booklets) >= 1 {
		lastCrateNumber = booklets[0].CrateNumber
		lastShelfNumber = booklets[0].ShelfNumber
	}
	return
}

func (srv *Server) getArchivingStats(crateNumber, shelfNumber, arcvhiveBoxNumber string) (stats archivingStat, err error) {
	stats = archivingStat{}
	_, err = srv.Db.Conn.QueryOne(&stats, `
	select *
	from (select count(distinct (b)) as booklet_in_crate, sum(b.size) as sheet_in_crate
      from booklets b,
           crates c
      where b.crate_number = c.number
        and c.number = ?
		and b.status != 'inArchiveStation') as crate
        ,
     (select count(distinct (c.number)) as crate_in_shelf, sum(b.size) as sheet_in_shelf
      from booklets b,
           crates c,
           shelves s
      where b.crate_number = c.number
        and c.shelf_number = s.number
        and s.number = ?
		and b.status != 'archived') as shelf;`, crateNumber, shelfNumber)
	if err != nil {
		return
	}

	archivedStats := archivingStat{}
	if arcvhiveBoxNumber != "" {
		_, err = srv.Db.Conn.QueryOne(&archivedStats, `
select count(b.number) as number_of_booklet_in_archive_box,
       sum(b.size)     as number_of_sheet_in_archive_box
from booklets b
where b.archive_box_number = ?;		
		`, arcvhiveBoxNumber)
	}
	if err != nil {
		return
	}

	stats.NumberOfSheetInArchiveBox = archivedStats.NumberOfSheetInArchiveBox
	stats.NumberOfBookletInArchiveBox = archivedStats.NumberOfBookletInArchiveBox
	return
}
