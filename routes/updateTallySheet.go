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

type updateTallySheetQuery struct {
	CurrentTallySheetNumber string `form:"CurrentTallySheetNumber"`
	NewTallySheetNumber     string `form:"NewTallySheetNumber"`
}

type updateTallySheetResult struct {
	nbUpdatedQuestionnaires int
	nbUpdatedTallySheets    int
	nbUpdatedBooklet        int
}

func (srv *Server) updateTallySheet(footer string) {
	srv.router.GET("/geocodeUpdate/updateTallySheetNumber.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)

		q := updateTallySheetQuery{}
		srv.updateTallySheetOkWithData(c, header, footer, q, nil)

	})

	srv.router.POST("/geocodeUpdate/updateTallySheetNumber.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)

		var q updateTallySheetQuery
		err := c.ShouldBind(&q)
		if err != nil {
			log.Error(err)
			srv.updateTallySheetWithError(
				c,
				header,
				footer,
				fmt.Sprintf("unparsable request : %s", err.Error()),
				q)
			return
		}
		result, err := srv.doUpdateTallySheet(q)
		if err != nil {
			srv.updateTallySheetWithError(c, header, footer, fmt.Sprintf("Error when trying to update GoeCode : %s", err), q)
		} else {
			srv.updateTallySheetOkWithData(c, header, footer, q, result)
		}
	})
}

func (srv *Server) updateTallySheetOkWithData(c *gin.Context, header, footer string, q updateTallySheetQuery, result *updateTallySheetResult) {
	name := ginoidc.GetValue(c, "name")
	if result != nil {
		c.HTML(http.StatusOK, "updateTallySheetNumber.html", gin.H{
			"Name":                    name,
			"CurrentTallySheetNumber": q.CurrentTallySheetNumber,
			"NewTallySheetNumber":     q.NewTallySheetNumber,
			"Header":                  template.HTML(header),
			"Footer":                  template.HTML(footer),
			"UpdateResult": template.HTML(fmt.Sprintf(`
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
                                </div>
								<div class="form-group remove-empty-values">
                                    <label class="control-label col-md-4 col-sm-4 col-xs-4">Nb of updated questionnaire(s) :</label>
                                    <div class="col-md-8 col-sm-8 col-xs-12">
                                        <h4 style="text-align:left" class="text-success">%d</h4>
                                    </div>
                                </div>`, result.nbUpdatedTallySheets, result.nbUpdatedBooklet, result.nbUpdatedQuestionnaires)),
		})
	} else {
		c.HTML(http.StatusOK, "updateTallySheetNumber.html", gin.H{
			"Name":                    name,
			"CurrentTallySheetNumber": q.CurrentTallySheetNumber,
			"NewTallySheetNumber":     q.NewTallySheetNumber,
			"Header":                  template.HTML(header),
			"Footer":                  template.HTML(footer),
		})
	}
}

func (srv *Server) updateTallySheetWithError(c *gin.Context, header, footer, alertMsg string, q updateTallySheetQuery) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := ginoidc.GetValue(c, "name")
	c.HTML(http.StatusOK, "updateTallySheetNumber.html", gin.H{
		"Name":                    name,
		"CurrentTallySheetNumber": q.CurrentTallySheetNumber,
		"NewTallySheetNumber":     q.NewTallySheetNumber,
		"Header":                  template.HTML(header),
		"Footer":                  template.HTML(footer),
		"Alert":                   template.HTML(alert),
	})
}

func (srv Server) doUpdateTallySheet(q updateTallySheetQuery) (result *updateTallySheetResult, err error) {
	tx, err := srv.Db.Conn.Begin()
	if err != nil {
		return result, err
	}

	// Update Booklet
	rb, err := tx.Model((*model.Booklet)(nil)).Exec(`
	update booklets set 
	number = ? 
	where number = ?;`,
		q.NewTallySheetNumber,
		q.CurrentTallySheetNumber)
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
	update tally_sheets set 
	tally_sheet_barcode = ?, 
	tally_sheet_no = ? 
	where tally_sheet_no = ?;`,
		strings.Replace(q.NewTallySheetNumber, ".", "", 1),
		strings.Replace(q.NewTallySheetNumber, ".", "", 1),
		strings.Replace(q.CurrentTallySheetNumber, ".", "", 1))
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

	// Update Questionnaires
	_, err = tx.Model((*goPg.Questionnaire)(nil)).Exec(`
	update questionnaires set 
    questionnaire_num = replace(questionnaire_num, ?, ?),
	tally_sheet_no    = ?,
    booklet_number    = ?
	where tally_sheet_no = ?;`,
		strings.Replace(q.CurrentTallySheetNumber, ".", "", 1),
		strings.Replace(q.NewTallySheetNumber, ".", "", 1),
		strings.Replace(q.NewTallySheetNumber, ".", "", 1),
		q.NewTallySheetNumber,
		strings.Replace(q.CurrentTallySheetNumber, ".", "", 1))
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

	// Update Questionnaire_num
	rqq, err := tx.Model((*goPg.Questionnaire)(nil)).Exec(`
	update questionnaires set
	questionnaire_num = replace(questionnaire_num, ?, ?)
	where tally_sheet_no = ?;`,
		strings.Replace(q.CurrentTallySheetNumber, ".", "", 1),
		strings.Replace(q.NewTallySheetNumber, ".", "", 1),
		strings.Replace(q.NewTallySheetNumber, ".", "", 1))
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

	result = &updateTallySheetResult{
		nbUpdatedQuestionnaires: rqq.RowsAffected(),
		nbUpdatedTallySheets:    rt.RowsAffected(),
		nbUpdatedBooklet:        rb.RowsAffected(),
	}
	return
}
