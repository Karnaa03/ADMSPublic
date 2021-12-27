package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	ginoidc "git.solutions.im/Solutions.IM/ginOidc"
	"git.solutions.im/XeroxAgriCensus/AgriInject/goPg"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/model"
	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

type updateQuery struct {
	BookletNumber  string `form:"BookletNumber"`
	CurrentGeoCode string `form:"CurrentGeocode"`
	NewGeocode     string `form:"NewGeocode"`
}

type updateResult struct {
	nbUpdatedQuestionnaires int
	nbUpdatedTallySheets    int
	nbUpdatedBooklet        int
}

func (srv *Server) updateGeoCode(footer string) {
	srv.router.GET("/geocodeUpdate/updateGeoCode.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)

		q := updateQuery{}
		srv.updateGeocodeOkWithData(c, header, footer, q, nil)

	})

	srv.router.POST("/geocodeUpdate/updateGeoCode.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)

		var q updateQuery
		err := c.ShouldBind(&q)
		if err != nil {
			log.Error(err)
			srv.updateGeocodeWithError(
				c,
				header,
				footer,
				fmt.Sprintf("unparsable request : %s", err.Error()),
				q)
			return
		}
		result, err := srv.doUpdateGeoCode(q)
		if err != nil {
			srv.updateGeocodeWithError(c, header, footer, fmt.Sprintf("Error when trying to update GoeCode : %s", err), q)
		} else {
			srv.updateGeocodeOkWithData(c, header, footer, q, result)
		}
	})
}

func (srv *Server) updateGeocodeOkWithData(c *gin.Context, header, footer string, q updateQuery, result *updateResult) {
	name := ginoidc.GetValue(c, "name")
	if result != nil {
		c.HTML(http.StatusOK, "updateGeoCode.html", gin.H{
			"Name":           name,
			"BookletNumber":  q.BookletNumber,
			"CurrentGeocode": q.CurrentGeoCode,
			"NewGeocode":     q.NewGeocode,
			"Header":         template.HTML(header),
			"Footer":         template.HTML(footer),
			"UpdateResult": template.HTML(fmt.Sprintf(`
								<div class="form-group remove-empty-values">
                                    <label class="control-label col-md-4 col-sm-4 col-xs-4">Nb of updated questionnaire(s) :</label>
                                    <div class="col-md-8 col-sm-8 col-xs-12">
                                        <h4 style="text-align:left" class="text-success">%d</h4>
                                    </div>
                                </div>
                                <div class="form-group remove-empty-values">
                                    <label class="control-label col-md-4 col-sm-4 col-xs-4">Nb of updated tally sheet(s) :</label>
                                    <div class="col-md-8 col-sm-8 col-xs-12">
                                        <h4 style="text-align:left" class="text-success">%d</h4>
                                    </div>
                                </div>
								<div class="form-group remove-empty-values">
                                    <label class="control-label col-md-4 col-sm-4 col-xs-4">Nb of updated booklet(s) :</label>
                                    <div class="col-md-8 col-sm-8 col-xs-12">
                                        <h4 style="text-align:left" class="text-success">%d</h4>
                                    </div>
                                </div>`, result.nbUpdatedQuestionnaires, result.nbUpdatedTallySheets, result.nbUpdatedBooklet)),
		})
	} else {
		c.HTML(http.StatusOK, "updateGeoCode.html", gin.H{
			"Name":           name,
			"BookletNumber":  q.BookletNumber,
			"CurrentGeocode": q.CurrentGeoCode,
			"NewGeocode":     q.NewGeocode,
			"Header":         template.HTML(header),
			"Footer":         template.HTML(footer),
		})
	}
}

func (srv *Server) updateGeocodeWithError(c *gin.Context, header, footer, alertMsg string, q updateQuery) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := ginoidc.GetValue(c, "name")
	c.HTML(http.StatusOK, "updateGeoCode.html", gin.H{
		"Name":           name,
		"BookletNumber":  q.BookletNumber,
		"CurrentGeocode": q.CurrentGeoCode,
		"NewGeocode":     q.NewGeocode,
		"Header":         template.HTML(header),
		"Footer":         template.HTML(footer),
		"Alert":          template.HTML(alert),
	})
}

func (srv Server) doUpdateGeoCode(q updateQuery) (result *updateResult, err error) {
	tx, err := srv.Db.Conn.Begin()
	if err != nil {
		return result, err
	}

	// Update questionnaires
	rq, err := tx.Model((*goPg.Questionnaire)(nil)).Exec(`
	update questionnaires
set geocode_id = ?
where tally_sheet_no = ?
  and geocode_id = ?;`, q.NewGeocode, strings.Replace(q.BookletNumber, ".", "", 1), q.CurrentGeoCode)
	if err != nil {
		errMsg := ""
		func() {
			errRollback := tx.Rollback()
			if errRollback != nil {
				errMsg = fmt.Sprintf("error when trying to rollback : %s, due to error ", errRollback)
			}
		}()
		return result, fmt.Errorf("%s%s", errMsg, err)
	}

	// Update Tally sheet
	rt, err := tx.Model((*goPg.TallySheet)(nil)).Exec(`
	update tally_sheets set geocode_id = ?
where tally_sheet_no = ?
  and geocode_id = ?;`, q.NewGeocode, strings.Replace(q.BookletNumber, ".", "", 1), q.CurrentGeoCode)
	if err != nil {
		errMsg := ""
		func() {
			errRollback := tx.Rollback()
			if errRollback != nil {
				errMsg = fmt.Sprintf("error when trying to rollback : %s, due to error ", errRollback)
			}
		}()
		return result, fmt.Errorf("%s%s", errMsg, err)
	}

	// Update booklet
	rb, err := tx.Model((*model.Booklet)(nil)).Exec(`
	update booklets set geo_code_id = ?
	where number = ?
  	and geo_code_id = ?;`, q.NewGeocode, q.BookletNumber, q.CurrentGeoCode)
	if err != nil {
		errMsg := ""
		func() {
			errRollback := tx.Rollback()
			if errRollback != nil {
				errMsg = fmt.Sprintf("error when trying to rollback : %s, due to error ", errRollback)
			}
		}()
		return result, fmt.Errorf("%s%s", errMsg, err)
	}

	err = tx.Commit()
	if err != nil {
		return result, fmt.Errorf("%s", err)
	}

	result = &updateResult{
		nbUpdatedQuestionnaires: rq.RowsAffected(),
		nbUpdatedTallySheets:    rt.RowsAffected(),
		nbUpdatedBooklet:        rb.RowsAffected(),
	}
	return
}
