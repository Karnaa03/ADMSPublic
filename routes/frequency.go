package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"git.solutions.im/XeroxAgriCensus/ADMSPublic/model"
	"git.solutions.im/XeroxAgriCensus/ADMSPublic/templates"
)

func (srv *Server) frequency(footer string) {
	srv.router.GET("/production/frequency.html", func(c *gin.Context) {
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
		srv.frequencyOkWithData(c, header, footer, &q, []model.RawTableData{})
	})

	srv.router.POST("/production/frequency.html", func(c *gin.Context) {
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
		data, err := srv.Db.GetOccupationOfHouseHold(division, district, upazilla, union, mouza, tableNumber)
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
		srv.frequencyOkWithData(c, header, footer, &q, data)

	})
}

func (srv *Server) frequencyOkWithData(c *gin.Context, header, footer string, q *searchQuery, data []model.RawTableData) {

	c.HTML(http.StatusOK, "frequency.html", gin.H{
		// "Name":                   name,
		"Header":         template.HTML(header),
		"Footer":         template.HTML(footer),
		"DivisionNumber": q.DivisionNumber,
		"DistrictNumber": q.DistrictNumber,
		"UpazilaNumber":  q.UpazilaNumber,
		"UnionNumber":    q.UnionNumber,
		"MouzaNumber":    q.MouzaNumber,
		"QueryType":      q.TableNumber,
		"TableData":      template.HTML(FormatFrequencyTable(data)),
		"Donuts":         template.HTML(FormatFrequencyDonuts(data)),
	})
}

func FormatFrequencyDonuts(data []model.RawTableData) (donuts string) {
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

func FormatFrequencyTable(data model.FrequencyListData) (tableData string) {
	p := message.NewPrinter(language.English)

	tableData += fmt.Sprintf(`
			<tr>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
			</tr>
			<tr>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
			</tr>`,
		"Holdings",
		p.Sprintf("%.2f", total),
		p.Sprintf("%.2f", urban),
		p.Sprintf("%.2f", rural),
		"Percentage",
		"100%",
		fmt.Sprintf("%.2f%%", (float64(urban)/float64(total))*100),
		fmt.Sprintf("%.2f%%", (float64(rural)/float64(total))*100),
	)

	return
}
