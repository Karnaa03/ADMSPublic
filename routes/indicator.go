package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	agriInject "git.solutions.im/XeroxAgriCensus/AgriInject/goPg"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/ADMSPublic/model"
	"git.solutions.im/XeroxAgriCensus/ADMSPublic/templates"
)

func (srv *Server) indicator(footer string) {
	srv.router.GET("/production/indicator.html", func(c *gin.Context) {
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
		srv.indicatorOkWithData(c, header, footer, &q, []model.Agregated{})
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

	srv.router.POST("/production/indicator.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		var q searchQuery
		err := c.ShouldBind(&q)
		division := strings.Trim(strings.Split(q.DivisionNumber, "-")[0], " ")
		district := strings.Trim(strings.Split(q.DistrictNumber, "-")[0], " ")
		upazilla := strings.Trim(strings.Split(q.UpazilaNumber, "-")[0], " ")
		union := strings.Trim(strings.Split(q.UnionNumber, "-")[0], " ")
		mouza := strings.Trim(strings.Split(q.MouzaNumber, "-")[0], " ")
		tableNumber := q.TableNumber
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
		println(division, district, upazilla, union, mouza, tableNumber)
		data, err := srv.Db.GetAgregate(division, district, upazilla, union, mouza, tableNumber)
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

		println(data)
		srv.indicatorOkWithData(c, header, footer, &q, data)

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

func (srv *Server) indicatorOkWithData(c *gin.Context, header, footer string, q *searchQuery, data []model.Agregated) {

	c.HTML(http.StatusOK, "indicator.html", gin.H{
		// "Name":                   name,
		"Header":         template.HTML(header),
		"Footer":         template.HTML(footer),
		"DivisionNumber": q.DivisionNumber,
		"DistrictNumber": q.DistrictNumber,
		"UpazilaNumber":  q.UpazilaNumber,
		"UnionNumber":    q.UnionNumber,
		"MouzaNumber":    q.MouzaNumber,
		"TableData":      template.HTML(FormatTable(data)),
	})
}

func (srv *Server) searchWithError(c *gin.Context, header, footer, alertMsg string, q searchQuery) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	c.HTML(http.StatusOK, "freeze.html", gin.H{
		"Header":         template.HTML(header),
		"Footer":         template.HTML(footer),
		"Alert":          template.HTML(alert),
		"DivisionNumber": q.DivisionNumber,
		"DistrictNumber": q.DistrictNumber,
		"UpazilaNumber":  q.UpazilaNumber,
		"UnionNumber":    q.UnionNumber,
		"MouzaNumber":    q.MouzaNumber,
	})
}

type searchQuery struct {
	DivisionNumber string
	DistrictNumber string
	UpazilaNumber  string
	UnionNumber    string
	MouzaNumber    string
	TableNumber    string
}

func (s searchQuery) IsEmpty() bool {
	if s.DivisionNumber == "" &&
		s.DistrictNumber == "" &&
		s.UpazilaNumber == "" &&
		s.UnionNumber == "" &&
		s.MouzaNumber == "" {
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
		err = req.Select()
	}
	return
}

func FormatTable(data []model.Agregated) (tableData string) {
	var urban, rural uint
	for _, line := range data {
		if line.Rmo == 2 {
			urban = line.HhSno
		} else {
			rural = line.HhSno
		}
	}
	total := urban + rural
	tableData += fmt.Sprintf(`
		<tr>
			<td>%s</td>
			<td>%d</td>
			<td>%d</td>
			<td>%d</td>
		</tr>
		<tr>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		</tr>`,
		"Holdings",
		total,
		urban,
		rural,
		"Percentage",
		"100%",
		fmt.Sprintf("%d%%", (urban/total)*100),
		fmt.Sprintf("%d%%", (rural/total)*100),
	)

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
