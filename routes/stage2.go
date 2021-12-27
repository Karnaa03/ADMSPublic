package routes

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	gin_oidc "git.solutions.im/Solutions.IM/ginOidc"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/model"
	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

type batchStats struct {
	BookletInCrate int
	SheetInCrate   int
	CrateInShelf   int
	SheetInShelf   int
}

func (srv *Server) stage2(footer string) {
	srv.router.GET("/production/stage2.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		type s2query struct {
			Shelf  string `form:"shelf"`
			Crate  string `form:"crate"`
			Delete string `form:"delete"`
		}
		var q s2query
		err := c.ShouldBind(&q)
		if err != nil {
			log.Error(err)
			srv.stage2WithError(
				c,
				header,
				footer,
				fmt.Sprintf("unparsable request : %s", err.Error()),
				"",
				"",
				"")
			return
		}

		switch {
		case q.Delete == "":
			srv.stage2OkWithData(c, header, footer, q.Crate, q.Shelf)
		case q.Delete != "":
			log.Infof("remove %s from batch in shelf %s and crate %s", q.Delete, q.Shelf, q.Crate)
			header, _ := templates.RenderHeader(c)
			id, err := gin_oidc.GetIdentity(c)
			if err != nil {
				log.Error(err)
				srv.stage2WithError(
					c,
					header,
					footer,
					"Who are you ???",
					q.Shelf,
					q.Crate,
					"")
				return
			}

			booklet, err := srv.Db.GetBooklet(q.Delete)
			if err != nil {
				log.Error(err)
				srv.stage2WithError(
					c,
					header,
					footer,
					err.Error(),
					q.Shelf,
					q.Crate,
					"")
				return
			}
			err = booklet.RemoveFromBatch(&srv.Db, id)
			if err != nil {
				log.Error(err)
				srv.stage2WithError(
					c,
					header,
					footer,
					err.Error(),
					q.Shelf,
					q.Crate,
					"")
				return
			}
			c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("/production/stage2.html?shelf=%s&crate=%s", q.Shelf, q.Crate))

		}
	})

	srv.router.POST("/production/stage2.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		postedShelfNumber := c.PostForm("ShelfNumber")
		postedCrateNumber := c.PostForm("CrateNumber")
		postedBookletNumber := c.PostForm("BookletNumber")
		id, err := gin_oidc.GetIdentity(c)
		if err != nil {
			log.Error(err)
			srv.stage2WithError(
				c,
				header,
				footer,
				"Who are you ???",
				postedShelfNumber,
				postedCrateNumber,
				postedBookletNumber)
			return
		}

		switch {
		case postedShelfNumber != "" && postedCrateNumber != "" && postedBookletNumber != "":
			log.Infof("%s %s %s", postedShelfNumber, postedCrateNumber, postedBookletNumber)
			// get if booklet exist
			booklet, err := srv.Db.GetBooklet(postedBookletNumber)
			if err != nil {
				log.Error(err)
				srv.stage2WithError(
					c,
					header,
					footer,
					"Unknown booklet, please check number or booklet registration",
					postedShelfNumber,
					postedCrateNumber,
					postedBookletNumber)
				return
			}

			err = booklet.RegisterInBatch(postedCrateNumber, postedShelfNumber, id, &srv.Db)
			if err != nil {
				log.Error(err)
				srv.stage2WithError(
					c,
					header,
					footer,
					err.Error(),
					postedShelfNumber,
					postedCrateNumber,
					postedBookletNumber)
				return
			}
			srv.stage2OkWithData(c, header, footer, postedCrateNumber, postedShelfNumber)
		default:
			srv.stage2OkWithData(c, header, footer, postedCrateNumber, postedShelfNumber)
		}

	})

}

func (srv *Server) stage2WithError(c *gin.Context, header, footer, alertMsg string, shelfNumber, crateNumber, bookletNumber string) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := gin_oidc.GetValue(c, "name")

	tableData, lastCrateNumber, lastShelfNumber := srv.getBatch(crateNumber, shelfNumber)
	stat, err := srv.getBatchStats(lastCrateNumber, lastShelfNumber)
	if err != nil {
		log.Error(err)
	}

	c.HTML(http.StatusOK, "stage2.html", gin.H{
		"Name":                       name,
		"Header":                     template.HTML(header),
		"Footer":                     template.HTML(footer),
		"Alert":                      template.HTML(alert),
		"ShelfNumber":                shelfNumber,
		"CrateNumber":                crateNumber,
		"BookletNumber":              bookletNumber,
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

func (srv *Server) stage2OkWithData(c *gin.Context, header, footer, crateNumber, shelfNumber string) {
	name := gin_oidc.GetValue(c, "name")
	tableData, lastCrateNumber, lastShelfNumber := srv.getBatch(crateNumber, shelfNumber)
	stat, err := srv.getBatchStats(lastCrateNumber, lastShelfNumber)
	if err != nil {
		log.Error(err)
	}

	c.HTML(http.StatusOK, "stage2.html", gin.H{
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

func (srv *Server) getBatch(crateNumber, shelfNumber string) (data, lastCrateNumber, lastShelfNumber string) {
	var batches []model.Booklet
	var err error
	switch {
	case crateNumber == "" && shelfNumber == "":
		// err = srv.Db.Conn.Model(&batches).Relation("Crate").Relation("Crate.Shelf").
		// 	Limit(50).
		// 	Where("Booklet.status IN (?)", pg.In([]string{"inBatch", "inCuttingStation"})).
		// 	Order("added_in_batch_on DESC").
		// 	Distinct().
		// 	Select()
	case crateNumber == "" && shelfNumber != "":
		err = srv.Db.Conn.Model(&batches).Relation("Crate").Relation("Crate.Shelf").
			// Where("Booklet.status IN (?)", pg.In([]string{"inBatch", "inCuttingStation"})).
			Where("Crate.shelf_number = ?", shelfNumber).
			Order("added_in_batch_on DESC").
			Distinct().
			Select()
	default:
		err = srv.Db.Conn.Model(&batches).Relation("Crate").Relation("Crate.Shelf").
			// Where("Booklet.status IN (?)", pg.In([]string{"inBatch", "inCuttingStation"})).
			Where("Booklet.crate_number = ?", crateNumber).
			Where("Crate.shelf_number = ?", shelfNumber).
			Order("added_in_batch_on DESC").
			Distinct().
			Select()
	}

	if err != nil {
		log.Error(err)
		return
	}
	for _, batch := range batches {
		switch {
		case batch.Status == "registered":
			data += "<tr>"
		case batch.Status == "inBatch":
			data += "<tr class=\"green\">"
		default:
			data += "<tr class=\"aero\">"
		}
		data += fmt.Sprintf(`
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%d</td>
			<td>%s</ts>
			<td>%s</td>
			<td>
				<a href="#" data-href="/production/stage2.html?shelf=%s&crate=%s&delete=%s" data-toggle="modal" data-target="#confirm-delete-batch">
					<center>
						<i class="fa fa-trash"></i>
					</center>
				</a>
			</td>
		</tr>
	`, batch.GetShelfNumber(),
			batch.GetCrateNumber(),
			batch.Number,
			batch.Size,
			batch.Status,
			batch.AddedInBatchBy,
			batch.Crate.ShelfNumber,
			batch.CrateNumber,
			batch.Number,
		)
	}
	if len(batches) >= 1 {
		lastCrateNumber = batches[0].CrateNumber
		lastShelfNumber = batches[0].Crate.ShelfNumber
	}
	return
}

func (srv *Server) getBatchStats(crateNumber, shelfNumber string) (stats batchStats, err error) {
	stats = batchStats{}
	_, err = srv.Db.Conn.QueryOne(&stats, `
	select *
	from (select count(distinct (b)) as booklet_in_crate, sum(b.size) as sheet_in_crate
      from booklets b,
           crates c
      where b.crate_number = c.number
        and c.number = ?) as crate
        ,
     (select count(distinct (c.number)) as crate_in_shelf, sum(b.size) as sheet_in_shelf
      from booklets b,
           crates c,
           shelves s
      where b.crate_number = c.number
        and c.shelf_number = s.number
        and s.number = ?) as shelf;`, crateNumber, shelfNumber)
	if err != nil {
		return
	}
	return
}
