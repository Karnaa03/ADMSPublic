package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/ADMSPublic/model"
	"git.solutions.im/XeroxAgriCensus/ADMSPublic/templates"
)

func (srv *Server) tableGeneration(footer string) {
	srv.router.GET("/production/table_generation.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)

		var q searchQuery
		err := c.ShouldBind(&q)
		if err != nil {
			log.Error(err)
			srv.tableGenerationWithError(
				c,
				header,
				footer,
				fmt.Sprintf("unparsable request : %s", err.Error()),
				q)
			return
		}
		srv.tableGenerationOkWithData(c, header, footer, &q, "")
	})

	srv.router.POST("/production/table_generation.html", func(c *gin.Context) {
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
			srv.tableGenerationWithError(
				c,
				header,
				footer,
				fmt.Sprintf("unparsable request : %s", err.Error()),
				q)
			return
		}
		var tableAndDonut string
		switch tableNumber {
		case "1":
			tableAndDonut, err = srv.FormatHouseholdHeadInformation(division, district, upazilla, union, mouza, &q)
		case "2":
			tableAndDonut, err = srv.FormatHouseholdLandInformation(division, district, upazilla, union, mouza, &q)
		case "3":
			tableAndDonut, err = srv.FormatHouseholdLandFisheryInformation(division, district, upazilla, union, mouza, &q)
		case "4":
			tableAndDonut, err = srv.FormatHouseholdPoultryInformation(division, district, upazilla, union, mouza, &q)
		}

		if err != nil {
			log.Error(err)
			srv.tableGenerationWithError(
				c,
				header,
				footer,
				fmt.Sprintf("unparsable request : %s", err.Error()),
				q)
			return
		}
		srv.tableGenerationOkWithData(c, header, footer, &q, tableAndDonut)

	})
}

func (srv *Server) tableGenerationWithError(c *gin.Context, header, footer, alertMsg string, q searchQuery) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	c.HTML(http.StatusOK, "table_generation.html", gin.H{
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

func (srv *Server) tableGenerationOkWithData(c *gin.Context, header, footer string, q *searchQuery, tableAndDonut string) {

	c.HTML(http.StatusOK, "table_generation.html", gin.H{
		"Header":         template.HTML(header),
		"Footer":         template.HTML(footer),
		"DivisionNumber": q.DivisionNumber,
		"DistrictNumber": q.DistrictNumber,
		"UpazilaNumber":  q.UpazilaNumber,
		"UnionNumber":    q.UnionNumber,
		"MouzaNumber":    q.MouzaNumber,
		"TableAndDonut":  template.HTML(tableAndDonut),
	})
}

func FormattableGenerationDonuts(data []model.RawTableData) (donuts string) {
	if len(data) > 0 {
		var urban, rural float64
		for _, line := range data {
			if line.Rmo == 2 {
				urban += line.Data
			} else {
				rural += line.Data
			}
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
					name: 'Access From',
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
	`, urban, "Urban", rural, "Rural")
	}
	return
}

func getTableGenerationName(tableNumber string) string {
	tableName := make(map[string]string)

	tableName["1"] = "Information of the household head"
	tableName["2"] = "Land information"
	tableName["3"] = "Household Fishery Land"
	tableName["4"] = "Household Poultry"
	tableName["5"] = "Household cattle"
	tableName["6"] = "Temporary Crops"
	tableName["7"] = "Permanent Crops"

	return tableName[tableNumber]
}
