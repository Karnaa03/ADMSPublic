package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	agriInject "git.solutions.im/XeroxAgriCensus/AgriInject/goPg"
	"github.com/gin-gonic/gin"
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
		srv.indicatorOkWithData(c, header, footer, &q, []model.TableData{})
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

func (srv *Server) indicatorOkWithData(c *gin.Context, header, footer string, q *searchQuery, data []model.TableData) {

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

func FormatTable(data []model.TableData) (tableData string) {
	if len(data) > 0 {
		var urban, rural uint
		for _, line := range data {
			if line.Rmo == 2 {
				urban = line.Count
			} else {
				rural = line.Count
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
	}
	return
}
