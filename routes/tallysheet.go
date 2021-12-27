package routes

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"

	ginoidc "git.solutions.im/Solutions.IM/ginOidc"
	"git.solutions.im/XeroxAgriCensus/AgriInject/goPg"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

func (srv *Server) tallySheet() {
	srv.router.GET("/adms/tallySheet.html", func(c *gin.Context) {
		tlno := c.Query("no")
		srv.tallySheetOkWithData(c, tlno)
	})

	srv.router.GET("/adms/getTallySheet", func(c *gin.Context) {
		tlno := c.Query("no")
		content, err := srv.S3.Get(srv.Config.S3Config.Bucket, fmt.Sprintf("%s/%s_001_A.pdf", tlno, tlno))
		if err != nil && err.Error() != io.EOF.Error() {
			log.Error(err)
		}
		c.Data(http.StatusOK, "application/pdf", content)
	})
}

func (srv *Server) tallySheetOkWithData(c *gin.Context, tlno string) {
	user, err := ginoidc.GetIdentity(c)
	if err != nil {
		log.Error(err)
	}
	tl, err := srv.GetTallySheetWithQuestionnaire(struct {
		BookletNumber  string
		DivisionNumber string
		DistrictNumber string
		UpazilaNumber  string
		UnionNumber    string
		MouzaNumber    string
		VillageNumber  string
		EA             string
		RMONumber      string
	}{BookletNumber: tlno})

	if err != nil || len(tl) != 1 {
		log.Errorf("wrong number of tallysheet or error : %s", err)
	}
	tls := TallySheets{tl[0]}
	c.HTML(http.StatusOK, "tallySheet.html", gin.H{
		"User": user,

		"FullGeoCode":            tl[0].GeoCode.GeocodeID,
		"TlNumber":               tl[0].TallySheetNo,
		"BookletNumber":          fmt.Sprintf("%s.%s", tl[0].TallySheetNo[0:3], tl[0].TallySheetNo[3:6]),
		"DivisionNumber":         fmt.Sprintf("%02d - %s", tl[0].GeoCode.Division, tl[0].GeoCode.NameDivision),
		"District":               fmt.Sprintf("%02d - %s", tl[0].GeoCode.District, tl[0].GeoCode.NameDistrict),
		"Upazila":                fmt.Sprintf("%02d - %s", tl[0].GeoCode.Upazilla, tl[0].GeoCode.NameUpazilla),
		"Union":                  fmt.Sprintf("%03d - %s", tl[0].GeoCode.Union, tl[0].GeoCode.NameUnion),
		"Mouza":                  fmt.Sprintf("%03d - %s", tl[0].GeoCode.Mouza, tl[0].GeoCode.NameMouza),
		"EA":                     fmt.Sprintf("%03d - %s", tl[0].GeoCode.CA, tl[0].GeoCode.NameCountingArea),
		"RMO":                    fmt.Sprintf("%d - %s", tl[0].GeoCode.Rmo, tl[0].GeoCode.NameRMO),
		"TotalHouseHolds":        tls.TotalHouseHolds(),
		"FishingHouse":           tls.FishingHouse(),
		"AgricultureLaborCode":   tls.AgriLaborHouse(),
		"HouseWithNoLand":        tls.HouseWithNoLand(),
		"HouseWithLandFromOther": tls.HouseWithLandFromOther(),
		"House5More":             tls.House5more(),
		"HousesWithFisheries":    tls.HouseWithFisheries(),
		"TotalCock":              tls.TotalCock(),
		"TotalDuck":              tls.TotalDuck(),
		"TotalTurkey":            tls.TotalTurkeys(),
		"TotalCow":               tls.TotalCow(),
		"TotalBuffalo":           tls.TotalBuffalo(),
		"TotalGoat":              tls.TotalGoat(),
		"TotalSheep":             tls.TotalSheep(),
		"TableData":              template.HTML(getTableData(tls)),
	})
}

func (srv *Server) GetTallySheetWithQuestionnaire(q searchQuery) (tl TallySheets, err error) {
	if !q.IsEmpty() {
		req := srv.Db.Conn.Model(&tl).
			Relation("GeoCode").
			Relation("Questionnaires")
		if q.BookletNumber != "" {
			req.Where("tally_sheet.tally_sheet_no = ?", strings.Replace(q.BookletNumber, ".", "", 1))
		}
		if q.DistrictNumber != "" {
			req.Where("geo_code.District = ?", q.DistrictNumber)
		}
		if q.DivisionNumber != "" {
			req.Where("geo_code.Division = ?", q.DivisionNumber)
		}
		if q.UpazilaNumber != "" {
			req.Where("geo_code.Upazilla = ?", q.UpazilaNumber)
		}
		if q.UnionNumber != "" {
			req.Where("geo_code.\"union\" = ?", q.UnionNumber)
		}
		if q.MouzaNumber != "" {
			req.Where("geo_code.Mouza = ?", q.MouzaNumber)
		}
		if q.EA != "" {
			req.Where("geo_code.CA = ?", q.EA)
		}
		if q.RMONumber != "" {
			req.Where("geo_code.Rmo = ?", q.RMONumber)
		}
		err = req.Select()
	}
	return
}

func (srv *Server) tallySheetWithError(c *gin.Context, alertMsg string) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := ginoidc.GetValue(c, "name")
	c.HTML(http.StatusOK, "tallySheet.html", gin.H{
		"Name":  name,
		"Alert": template.HTML(alert),
	})
}

func getTableData(tls TallySheets) (tableData string) {
	for _, tl := range tls {
		for _, q := range tl.Questionnaires {
			tableData += fmt.Sprintf(`
		  <tr>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		  	<td>
				<a href="/adms/questionnaire.html?TlNumber=%s&QNumber=%s" target="_blank">
					<center>
						<i class="fa fa-search"></i>
					</center>
				</a>
			</td>
		  </tr>`,
				getQnumStr(q),
				fmt.Sprintf("%03d", q.Village),
				fmt.Sprintf("%03d", q.HouseSerial),
				q.OwnerName,
				q.TallySheetNo,
				q.QuestionnaireNum,
			)
		}
	}
	return
}

func getQnumStr(q *goPg.Questionnaire) string {
	if q.QuestionnaireEmpty {
		return fmt.Sprintf("%s <i class=\"fa fa-file-o\"></i>", q.QuestionnaireNum)
	} else {
		return fmt.Sprintf("<b>%s</b> <i class=\"fa fa-file-text\"></i>", q.QuestionnaireNum)
	}
}
