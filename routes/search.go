package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	ginoidc "git.solutions.im/Solutions.IM/ginOidc"
	agriInject "git.solutions.im/XeroxAgriCensus/AgriInject/goPg"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/model"
	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

func (srv *Server) search(footer string) {
	srv.router.GET("/adms/search.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)

		var q searchQuery
		err := c.ShouldBind(&q)
		if err != nil {
			log.Error(err)
			srv.searchWithError(
				c,
				header,
				footer,
				fmt.Sprintf("unparsable request : %s", err.Error()),
				q)
			return
		}
		srv.searchOkWithData(c, header, footer, &q)
	})

	srv.router.GET("/adms/division", func(context *gin.Context) {
		term := context.Query("query")
		s := struct {
			Query       string   `json:"query"`
			Suggestions []string `json:"suggestions"`
		}{}
		if term != "" {
			s.Query = term
			err := srv.Db.Conn.Model((*model.GeoCodes)(nil)).
				ColumnExpr("distinct (division || ' - ' || name_division)").
				Where("(division || ' - ' || name_division like ?)", fmt.Sprintf("%%%s%%", strings.ReplaceAll(strings.ToUpper(s.Query), " ", "%"))).
				Select(&s.Suggestions)
			if err != nil {
				log.Error(err)
			}
		}
		context.JSON(http.StatusOK, s)
	})

	srv.router.GET("/adms/district", func(context *gin.Context) {
		s := struct {
			Query       string   `json:"query"`
			Suggestions []string `json:"suggestions"`
		}{}
		term := context.Query("query")
		division := context.Query("division")
		if term != "" {
			divisionNumber := strings.Split(division, "-")[0]
			s.Query = term
			query := srv.Db.Conn.Model((*model.GeoCodes)(nil)).
				ColumnExpr("distinct (district || ' - ' || name_district)").
				Where("(district || ' - ' || name_district like ?)", fmt.Sprintf("%%%s%%", strings.ReplaceAll(strings.ToUpper(s.Query), " ", "%")))
			if divisionNumber != "" {
				query.Where("division = ?", divisionNumber)
			}
			err := query.Select(&s.Suggestions)
			if err != nil {
				log.Error(err)
			}
		}
		context.JSON(http.StatusOK, s)
	})

	srv.router.GET("/adms/upazilla", func(context *gin.Context) {
		term := context.Query("query")
		s := struct {
			Query       string   `json:"query"`
			Suggestions []string `json:"suggestions"`
		}{}
		division := context.Query("division")
		district := context.Query("district")
		if term != "" {
			divisionNumber := strings.Split(division, "-")[0]
			districtNumber := strings.Split(district, "-")[0]
			s.Query = term
			query := srv.Db.Conn.Model((*model.GeoCodes)(nil)).
				ColumnExpr("distinct (upazilla || ' - ' || name_upazilla)").
				Where("(upazilla || ' - ' || name_upazilla like ?)", fmt.Sprintf("%%%s%%", strings.ReplaceAll(strings.ToUpper(s.Query), " ", "%")))
			if divisionNumber != "" {
				query.Where("division = ?", divisionNumber)
			}
			if districtNumber != "" {
				query.Where("district = ?", districtNumber)
			}
			err := query.Select(&s.Suggestions)
			if err != nil {
				log.Error(err)
			}
		}
		context.JSON(http.StatusOK, s)
	})

	srv.router.GET("/adms/union", func(context *gin.Context) {
		term := context.Query("query")
		s := struct {
			Query       string   `json:"query"`
			Suggestions []string `json:"suggestions"`
		}{}
		division := context.Query("division")
		district := context.Query("district")
		upazila := context.Query("upazila")
		if term != "" {
			divisionNumber := strings.Split(division, "-")[0]
			districtNumber := strings.Split(district, "-")[0]
			upazilaNumber := strings.Split(upazila, "-")[0]
			s.Query = term
			query := srv.Db.Conn.Model((*model.GeoCodes)(nil)).
				ColumnExpr("distinct (\"union\" || ' - ' || name_union)").
				Where("(\"union\" || ' - ' || name_union) like ?", fmt.Sprintf("%%%s%%", strings.ReplaceAll(strings.ToUpper(s.Query), " ", "%")))
			if divisionNumber != "" {
				query.Where("division = ?", divisionNumber)
			}
			if districtNumber != "" {
				query.Where("district = ?", districtNumber)
			}
			if upazilaNumber != "" {
				query.Where("upazilla = ?", upazilaNumber)
			}
			err := query.Select(&s.Suggestions)
			if err != nil {
				log.Error(err)
			}
		}
		context.JSON(http.StatusOK, s)
	})

	srv.router.GET("/adms/mouza", func(context *gin.Context) {
		term := context.Query("query")
		s := struct {
			Query       string   `json:"query"`
			Suggestions []string `json:"suggestions"`
		}{}
		division := context.Query("division")
		district := context.Query("district")
		upazila := context.Query("upazila")
		union := context.Query("union")
		if term != "" {
			divisionNumber := strings.Split(division, "-")[0]
			districtNumber := strings.Split(district, "-")[0]
			upazilaNumber := strings.Split(upazila, "-")[0]
			unionNumber := strings.Split(union, "-")[0]
			s.Query = term
			query := srv.Db.Conn.Model((*model.GeoCodes)(nil)).
				ColumnExpr("distinct(mouza || ' - ' || name_mouza)").
				Where("(mouza || ' - ' || name_mouza) like ?", fmt.Sprintf("%%%s%%", strings.ReplaceAll(strings.ToUpper(s.Query), " ", "%")))
			if divisionNumber != "" {
				query.Where("division = ?", divisionNumber)
			}
			if districtNumber != "" {
				query.Where("district = ?", districtNumber)
			}
			if upazilaNumber != "" {
				query.Where("upazilla = ?", upazilaNumber)
			}
			if unionNumber != "" {
				query.Where("\"union\" = ?", unionNumber)
			}
			err := query.Select(&s.Suggestions)
			if err != nil {
				log.Error(err)
			}
		}
		context.JSON(http.StatusOK, s)
	})

	srv.router.POST("/adms/search.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		var q searchQuery
		err := c.ShouldBind(&q)
		if err != nil {
			log.Error(err)
			srv.searchWithError(
				c,
				header,
				footer,
				fmt.Sprintf("unparsable request : %s", err.Error()),
				q)
			return
		}
		srv.searchOkWithData(c, header, footer, &q)

	})
}

func getNumber(numberAndName string) string {
	parts := strings.Split(numberAndName, "-")
	if len(parts) > 1 {
		return parts[0]
	} else {
		return ""
	}
}

func (srv *Server) searchOkWithData(c *gin.Context, header, footer string, q *searchQuery) {
	name := ginoidc.GetValue(c, "name")
	tl, err := srv.GetTallySheet(*q)
	if err != nil {
		log.Error(err)
	}
	hf, err := srv.GetHouseWithFisheries(*q)
	if err != nil {
		log.Error(err)
	}
	g, err := srv.GetGeoCodeNames(*q)
	if err != nil {
		log.Error(err)
	}
	if g.Division != 0 {
		q.DivisionNumber = fmt.Sprintf("%d - %s", g.Division, g.NameDivision)
	}

	formater := message.NewPrinter(language.English)
	c.HTML(http.StatusOK, "search.html", gin.H{
		"Name":                   name,
		"Header":                 template.HTML(header),
		"Footer":                 template.HTML(footer),
		"BookletNumber":          q.BookletNumber,
		"DivisionNumber":         q.DivisionNumber,
		"DistrictNumber":         q.DistrictNumber,
		"UpazilaNumber":          q.UpazilaNumber,
		"UnionNumber":            q.UnionNumber,
		"MouzaNumber":            q.MouzaNumber,
		"EA":                     q.EA,
		"EaName":                 g.NameCountingArea,
		"RMONumber":              q.RMONumber,
		"RMOName":                g.NameRMO,
		"TableData":              template.HTML(FormatTable(tl)),
		"TotalHouseHolds":        formater.Sprintf("%d", tl.TotalHouseHolds()),
		"TotalCockHen":           formater.Sprintf("%d", tl.TotalCock()),
		"FishingHouse":           formater.Sprintf("%d", tl.FishingHouse()),
		"TotalDuck":              formater.Sprintf("%d", tl.TotalDuck()),
		"ALH":                    formater.Sprintf("%d", tl.AgriLaborHouse()),
		"TotalTurkey":            formater.Sprintf("%d", tl.TotalTurkeys()),
		"HouseWithNoLand":        formater.Sprintf("%d", tl.HouseWithNoLand()),
		"TotalCow":               formater.Sprintf("%d", tl.TotalCow()),
		"HouseWithLandFromOther": formater.Sprintf("%d", tl.HouseWithLandFromOther()),
		"TotalBuffalo":           formater.Sprintf("%d", tl.TotalBuffalo()),
		"House5More":             formater.Sprintf("%d", tl.House5more()),
		"TotalGoat":              formater.Sprintf("%d", tl.TotalGoat()),
		"HouseWithFisheries":     formater.Sprintf("%d", hf),
		"TotalSheep":             formater.Sprintf("%d", tl.TotalSheep()),
		"InitDivision":           template.HTML(srv.autoCompleteGeoCode),
	})
}

func (srv *Server) searchWithError(c *gin.Context, header, footer, alertMsg string, q searchQuery) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := ginoidc.GetValue(c, "name")
	c.HTML(http.StatusOK, "freeze.html", gin.H{
		"Name":           name,
		"Header":         template.HTML(header),
		"Footer":         template.HTML(footer),
		"Alert":          template.HTML(alert),
		"BookletNumber":  q.BookletNumber,
		"DivisionNumber": q.DivisionNumber,
		"DistrictNumber": q.DistrictNumber,
		"UpazilaNumber":  q.UpazilaNumber,
		"UnionNumber":    q.UnionNumber,
		"MouzaNumber":    q.MouzaNumber,
		"VillageNumber":  q.VillageNumber,
		"EA":             q.EA,
		"RMONumber":      q.RMONumber,
	})
}

type searchQuery struct {
	BookletNumber  string
	DivisionNumber string
	DistrictNumber string
	UpazilaNumber  string
	UnionNumber    string
	MouzaNumber    string
	VillageNumber  string
	EA             string
	RMONumber      string
}

func (s searchQuery) IsEmpty() bool {
	if s.BookletNumber == "" &&
		s.DivisionNumber == "" &&
		s.DistrictNumber == "" &&
		s.UpazilaNumber == "" &&
		s.UnionNumber == "" &&
		s.MouzaNumber == "" &&
		s.VillageNumber == "" &&
		s.EA == "" &&
		s.RMONumber == "" {
		return true
	}
	return false
}

type TallySheets []agriInject.TallySheet

func (srv *Server) GetGeoCodeNames(q searchQuery) (g model.GeoCodes, err error) {
	var geocodes []model.GeoCodes
	if !q.IsEmpty() {
		req := srv.Db.Conn.Model(&geocodes)
		if q.DistrictNumber != "" {
			req.Where("District = ?", getNumber(q.DistrictNumber))
		}
		if q.DivisionNumber != "" {
			req.Where("Division = ?", getNumber(q.DivisionNumber))
		}
		if q.UpazilaNumber != "" {
			req.Where("Upazilla = ?", getNumber(q.UpazilaNumber))
		}
		if q.UnionNumber != "" {
			req.Where("\"union\" = ?", getNumber(q.UnionNumber))
		}
		if q.MouzaNumber != "" {
			req.Where("Mouza = ?", getNumber(q.MouzaNumber))
		}
		if q.EA != "" {
			req.Where("CA = ?", q.EA)
		}
		if q.RMONumber != "" {
			req.Where("Rmo = ?", q.RMONumber)
		}
		err = req.Select()
		if err != nil {
			return model.GeoCodes{}, err
		}
		if len(geocodes) > 0 {
			if q.DistrictNumber != "" {
				g.NameDistrict = geocodes[0].NameDistrict
				g.NameDivision = geocodes[0].NameDivision
				g.Division = geocodes[0].Division
			}
			if q.DivisionNumber != "" {
				g.NameDivision = geocodes[0].NameDivision
			}
			if q.UpazilaNumber != "" {
				g.NameUpazilla = geocodes[0].NameUpazilla
			}
			if q.UnionNumber != "" {
				g.NameUnion = geocodes[0].NameUnion
			}
			if q.MouzaNumber != "" {
				g.NameMouza = geocodes[0].NameMouza
			}
			if q.EA != "" {
				g.NameCountingArea = geocodes[0].NameCountingArea
			}
			if q.RMONumber != "" {
				g.NameRMO = geocodes[0].NameRMO
			}
		}
	}
	return
}

func (srv *Server) GetHouseWithFisheries(q searchQuery) (sum int, err error) {
	if !q.IsEmpty() {

		baseReq := `
		SELECT count(*)
		FROM questionnaires q, geo_codes g
		WHERE q.geocode_id = g.geocode_id
		  AND (GREATEST(0,q.pond_land)
          + GREATEST(0,q.fish_cultivation_land)
          + GREATEST(0,q.paddy_cultivation_land)
          + GREATEST(0,q.mixed_cultivation_land)
          + GREATEST(0,q.fish_salt_cultive_land)
          + GREATEST(0,q.fish_cage_cultive_land)
          + GREATEST(0,q.creek_land) > 30)`
		if q.BookletNumber != "" {
			baseReq += fmt.Sprintf(" AND q.tally_sheet_no = '%s'", strings.Replace(q.BookletNumber, ".", "", 1))
		}
		if q.DistrictNumber != "" {
			baseReq += fmt.Sprintf(" AND g.District = %s", getNumber(q.DistrictNumber))
		}
		if q.DivisionNumber != "" {
			baseReq += fmt.Sprintf(" AND g.Division = %s", getNumber(q.DivisionNumber))
		}
		if q.UpazilaNumber != "" {
			baseReq += fmt.Sprintf(" AND g.Upazilla = %s", getNumber(q.UpazilaNumber))
		}
		if q.UnionNumber != "" {
			baseReq += fmt.Sprintf(" AND g.\"union\" = %s", getNumber(q.UnionNumber))
		}
		if q.MouzaNumber != "" {
			baseReq += fmt.Sprintf(" AND g.Mouza = %s", getNumber(q.MouzaNumber))
		}
		if q.EA != "" {
			baseReq += fmt.Sprintf(" AND g.CA = %s", q.EA)
		}
		if q.RMONumber != "" {
			baseReq += fmt.Sprintf(" AND g.Rmo = %s", q.RMONumber)
		}

		_, err = srv.Db.Conn.Model(&sum).QueryOne(pg.Scan(&sum), baseReq)
		if err != nil {
			return
		}
	}
	return
}

func (srv *Server) GetTallySheet(q searchQuery) (tl TallySheets, err error) {
	if !q.IsEmpty() {
		req := srv.Db.Conn.Model(&tl).
			Relation("GeoCode")
			// Relation("Questionnaires")
		if q.BookletNumber != "" {
			req.Where("tally_sheet.tally_sheet_no = ?", strings.Replace(q.BookletNumber, ".", "", 1))
		}
		if q.DistrictNumber != "" {
			req.Where("geo_code.District = ?", getNumber(q.DistrictNumber))
		}
		if q.DivisionNumber != "" {
			req.Where("geo_code.Division = ?", getNumber(q.DivisionNumber))
		}
		if q.UpazilaNumber != "" {
			req.Where("geo_code.Upazilla = ?", getNumber(q.UpazilaNumber))
		}
		if q.UnionNumber != "" {
			req.Where("geo_code.\"union\" = ?", getNumber(q.UnionNumber))
		}
		if q.MouzaNumber != "" {
			req.Where("geo_code.Mouza = ?", getNumber(q.MouzaNumber))
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

func FormatTable(tl TallySheets) (tableData string) {
	for _, sheet := range tl {
		tableData += fmt.Sprintf(`
		<tr>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%d</td>
			<td>
				<a href="/adms/tallySheet.html?no=%s" target="_blank">
					<center>
						<i class="fa fa-search"></i>
					</center>
				</a>
			</td>
		</tr>`,
			sheet.TallySheetNo,
			sheet.GeocodeID,
			fmt.Sprintf("%d - %s", sheet.GeoCode.District, sheet.GeoCode.NameDistrict),
			fmt.Sprintf("%d - %s", sheet.GeoCode.Upazilla, sheet.GeoCode.NameUpazilla),
			fmt.Sprintf("%d - %s", sheet.GeoCode.Union, sheet.GeoCode.NameUnion),
			fmt.Sprintf("%d - %s", sheet.GeoCode.Mouza, sheet.GeoCode.NameMouza),
			*sheet.UpdatedTotalHouse,
			sheet.TallySheetNo,
		)
	}
	return
}

func (t TallySheets) TotalHouseHolds() (th int) {
	for _, sheet := range t {
		if sheet.UpdatedTotalHouse != nil {
			th += *sheet.UpdatedTotalHouse
		} else {
			th += sheet.OriginalTotalHouse
		}
	}
	return
}

func (t TallySheets) FishingHouse() (fh int) {
	for _, sheet := range t {
		if sheet.UpdatedHouseFisheries != nil {
			fh += *sheet.UpdatedHouseFisheries
		} else {
			fh += sheet.OriginalHouseFisheries
		}
	}
	return
}

func (t TallySheets) AgriLaborHouse() (al int) {
	for _, sheet := range t {
		if sheet.UpdatedAgriProfessionals != nil {
			al += *sheet.UpdatedAgriProfessionals
		} else {
			al += sheet.OriginalAgriProfessionals
		}
	}
	return
}

func (t TallySheets) HouseWithNoLand() (nl int) {
	for _, sheet := range t {
		if sheet.UpdatedHouseNoLand != nil {
			nl += *sheet.UpdatedHouseNoLand
		} else {
			nl += sheet.OriginalHouseNoLand
		}
	}
	return
}

func (t TallySheets) HouseWithLandFromOther() (ho int) {
	for _, sheet := range t {
		if sheet.UpdatedHouseReceivedLand != nil {
			ho += *sheet.UpdatedHouseReceivedLand
		} else {
			ho += sheet.OriginalHouseReceivedLand
		}
	}
	return
}

func (t TallySheets) House5more() (hm int) {
	for _, sheet := range t {
		if sheet.UpdatedHouse5More != nil {
			hm += *sheet.UpdatedHouse5More
		} else {
			hm += sheet.OriginalHouse5More
		}
	}
	return
}

func (t TallySheets) HouseWithFisheries() (hf int) {
	for _, sheet := range t {
		var totalSurface float64
		for _, questionnaire := range sheet.Questionnaires {
			totalSurface +=
				questionnaire.PondLand +
					questionnaire.FishCultivationLand +
					questionnaire.PaddyCultivationLand +
					questionnaire.MixedCultivationLand +
					questionnaire.FishSaltCultiveLand +
					questionnaire.FishCageCultiveLand +
					questionnaire.CreekLand
		}
		if totalSurface > 30 {
			hf += 1
		}
	}
	return
}

func (t TallySheets) TotalCock() (tc int) {
	for _, sheet := range t {
		if sheet.UpdatedCockCount != nil {
			tc += *sheet.UpdatedCockCount
		} else {
			tc += sheet.OriginalCockCount
		}
	}
	return
}

func (t TallySheets) TotalDuck() (td int) {
	for _, sheet := range t {
		if sheet.UpdatedDuckCount != nil {
			td += *sheet.UpdatedDuckCount
		} else {
			td += sheet.OriginalDuckCount
		}
	}
	return
}

func (t TallySheets) TotalTurkeys() (tt int) {
	for _, sheet := range t {
		if sheet.UpdatedTurkeyCount != nil {
			tt += *sheet.UpdatedTurkeyCount
		} else {
			tt += sheet.OriginalTurkeyCount
		}
	}
	return
}

func (t TallySheets) TotalCow() (tc int) {
	for _, sheet := range t {
		if sheet.UpdatedCowCount != nil {
			tc += *sheet.UpdatedCowCount
		} else {
			tc += sheet.OriginalCowCount
		}
	}
	return
}

func (t TallySheets) TotalBuffalo() (tb int) {
	for _, sheet := range t {
		if sheet.UpdatedBuffaloCount != nil {
			tb += *sheet.UpdatedBuffaloCount
		} else {
			tb += sheet.OriginalBuffaloCount
		}
	}
	return
}

func (t TallySheets) TotalGoat() (tg int) {
	for _, sheet := range t {
		if sheet.UpdatedGoatCount != nil {
			tg += *sheet.UpdatedGoatCount
		} else {
			tg += sheet.OriginalGoatCount
		}
	}
	return
}

func (t TallySheets) TotalSheep() (ts int) {
	for _, sheet := range t {
		if sheet.UpdatedSheepCount != nil {
			ts += *sheet.UpdatedSheepCount
		} else {
			ts += sheet.OriginalSheepCount
		}
	}
	return
}
