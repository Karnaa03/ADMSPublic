package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	agriInject "git.solutions.im/XeroxAgriCensus/AgriInject/goPg"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

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
		srv.indicatorOkWithData(c, header, footer, &q, []model.RawTableData{}, "")
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
				ColumnExpr("distinct (division || ' - ' || name_division) as di").
				Where("(division || ' - ' || name_division like ?)", fmt.Sprintf("%%%s%%", strings.ReplaceAll(strings.ToUpper(s.Query), " ", "%"))).
				Order("di").
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
				ColumnExpr("distinct (district || ' - ' || name_district) as di").
				Where("(district || ' - ' || name_district like ?)", fmt.Sprintf("%%%s%%", strings.ReplaceAll(strings.ToUpper(s.Query), " ", "%"))).
				Order("di")
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
				ColumnExpr("distinct (upazilla || ' - ' || name_upazilla) as up").
				Where("(upazilla || ' - ' || name_upazilla like ?)", fmt.Sprintf("%%%s%%", strings.ReplaceAll(strings.ToUpper(s.Query), " ", "%"))).
				Order("up")
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
				ColumnExpr("distinct (\"union\" || ' - ' || name_union) as un").
				Where("(\"union\" || ' - ' || name_union) like ?", fmt.Sprintf("%%%s%%", strings.ReplaceAll(strings.ToUpper(s.Query), " ", "%"))).
				Order("un")
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
				ColumnExpr("distinct(mouza || ' - ' || name_mouza) as mo").
				Where("(mouza || ' - ' || name_mouza) like ?", fmt.Sprintf("%%%s%%", strings.ReplaceAll(strings.ToUpper(s.Query), " ", "%"))).
				Order("mo")
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
		geoLocation := formatGeoSelection(q.DivisionNumber, q.DistrictNumber, q.UpazilaNumber, q.UnionNumber, q.MouzaNumber)
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
		srv.indicatorOkWithData(c, header, footer, &q, data, geoLocation)

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

func (srv *Server) indicatorOkWithData(c *gin.Context, header, footer string, q *searchQuery, data []model.RawTableData, geoLocation string) {
	tableName := make(map[string]string)
	tableName["1"] = "Total number of holdings"
	tableName["2"] = "Total number of farm holdings & percentage"
	tableName["3"] = "Total number of small farm & percentage"
	tableName["4"] = "Total number of medium farm & percentage"
	tableName["5"] = "Total number of large farm & percentage"
	tableName["6"] = "Total number of landless (owning no land) & percentage"
	tableName["7"] = "Total number of agriculture labor household & percentage"
	tableName["8"] = "Total number of fishery household & percentage"
	tableName["9"] = "Agriculture worker per household"
	tableName["10"] = "Agriculture worker per farm household"
	tableName["11"] = "Operated area (in acre) per holding"
	tableName["12"] = "Operated area (in acre) per farm holding"
	tableName["13"] = "Homestead area (in acre) per holding"
	tableName["14"] = "Homestead area (in acre) per farm holding"
	tableName["15"] = "Total number of male headed households & percentage"
	tableName["16"] = "Total Number of female headed households & percentage"
	tableName["17"] = "Total number of holding reporting irrigation & percentage"
	tableName["18"] = "Total number of holding area (in acre) irrigated & percentage"
	tableName["19"] = "Total number of farms holding reporting irrigation & percentage"
	tableName["20"] = "Total number of farms holding area (in acre) irrigated & percentage "
	tableName["21"] = "Total number of cow"
	tableName["22"] = "Total number of buffalo"
	tableName["23"] = "Total number of goat"
	tableName["24"] = "Total number of sheep "
	tableName["25"] = "Total number of cock/hen"
	tableName["26"] = "Total number of duck"
	tableName["27"] = "Cropping Intensity"
	tableName["28"] = "Net cultivated area (in acre)"

	c.HTML(http.StatusOK, "indicator.html", gin.H{
		// "Name":                   name,
		"Header":         template.HTML(header),
		"Footer":         template.HTML(footer),
		"DivisionNumber": q.DivisionNumber,
		"DistrictNumber": q.DistrictNumber,
		"UpazilaNumber":  q.UpazilaNumber,
		"UnionNumber":    q.UnionNumber,
		"MouzaNumber":    q.MouzaNumber,
		"QueryType":      fmt.Sprintf("%s<br>%s", tableName[q.TableNumber], geoLocation),
		"TableData":      template.HTML(FormatTable(data)),
		"Donuts":         template.HTML(FormatDonuts(data)),
	})
}

func (srv *Server) searchWithError(c *gin.Context, header, footer, alertMsg string, q searchQuery) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	c.HTML(http.StatusOK, "indicator.html", gin.H{
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
func FormatFloat(num float64, prc int) string {
	p := message.NewPrinter(language.English)
	var (
		zero, dot = "0", "."
		str       = p.Sprintf("%."+strconv.Itoa(prc)+"f", num)
	)
	return strings.TrimRight(strings.TrimRight(str, zero), dot)
}

func FormatDonuts(data []model.RawTableData) (donuts string) {
	if len(data) > 0 {
		var urban, rural, total float64
		with_percentage := true
		for _, line := range data {
			switch line.Rmo {
			case 2:
				urban += line.Data
			case -1:
				total = line.Data
				with_percentage = false
			default:
				rural += line.Data
			}
		}
		if total == 0.0 {
			total = urban + rural
		}
		donuts = fmt.Sprintf(`
		<div id="main" style="width: 600px;height:400px; align:center" class="x_content"></div>
		<script type="text/javascript">	
		var chartDom = document.getElementById('main');
		var myChart = echarts.init(chartDom);
		var option;
	
		option = {
			tooltip: {
				trigger: 'item'
			},
			legend: {
				top: '5%%',
				left: 'center'
			},
			series: [
				{
					name: '',
					type: 'pie',
					radius: ['40%%', '70%%'],
					avoidLabelOverlap: false,
					itemStyle: {
						borderRadius: 10,
						borderColor: '#fff',
						borderWidth: 2
					},
					label: {
						show: false,
						position: 'center'
					},
					emphasis: {
						label: {
							show: true,
							fontSize: '40',
							fontWeight: 'bold'
						}
					},
					labelLine: {
						show: false
					},
					data: [
						{ value: %f, name: '%s' },
						{ value: %f, name: '%s' },
					]
				}
			]
		};
		option && myChart.setOption(option);
	</script>
	`,
			urban,
			GetUrban(with_percentage, urban, total),
			rural,
			GetRural(with_percentage, rural, total),
		)
	}
	return
}

func GetUrban(with_percentage bool, urban, total float64) string {
	if with_percentage {
		return fmt.Sprintf("Urban : %.2f%%", (float64(urban)/float64(total))*100)
	} else {
		return fmt.Sprintf("Urban : %.2f", urban)
	}
}

func GetRural(with_percentage bool, rural, total float64) string {
	if with_percentage {
		return fmt.Sprintf("Rural: %.2f%%", (float64(rural)/float64(total))*100)
	} else {
		return fmt.Sprintf("Rural: %.2f", rural)
	}
}

func FormatTable(data []model.RawTableData) (tableData string) {
	if len(data) > 0 {
		var urban, rural, total float64
		with_percentage := true
		for _, line := range data {
			switch line.Rmo {
			case 2:
				urban += line.Data
			case -1:
				total = line.Data
				with_percentage = false
			default:
				rural += line.Data
			}
		}
		if total == 0.0 {
			total = urban + rural
		}
		tableData += fmt.Sprintf(`
		<tr>
		<th>Total</th>
        <th>Urban</th>
        <th>Rural</th>
		</tr>
		<tr>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
		</tr>
			`,
			FormatFloat(total, 2),
			FormatFloat(urban, 2),
			FormatFloat(rural, 2),
		)

		if with_percentage {
			tableData += fmt.Sprintf(`
			<tr>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
			</tr>
			`,
				"100%",
				fmt.Sprintf("%.2f%%", (float64(urban)/float64(total))*100),
				fmt.Sprintf("%.2f%%", (float64(rural)/float64(total))*100),
			)
		}
		// tableData += fmt.Sprintf(`
		// <tr>
		// 	<td>%s</td>
		// 	<td></td>
		// 	<td></td>
		// </tr>`, "test source")

	}
	return
}
