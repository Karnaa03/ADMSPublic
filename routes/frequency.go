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
			srv.frequencyWithError(
				c,
				header,
				footer,
				fmt.Sprintf("unparsable request : %s", err.Error()),
				q)
			return
		}
		srv.frequencyOkWithData(c, header, footer, &q, "")
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
			srv.frequencyWithError(
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
			tableAndDonut, err = srv.FormatOccupationOfTheHouseHold(division, district, upazilla, union, mouza, &q)
		case "2":
			tableAndDonut, err = srv.FormatEducationHouseHoldHead(division, district, upazilla, union, mouza, &q)
		case "3":
			tableAndDonut, err = srv.FormatGenderOfTheHouseholdHead(division, district, upazilla, union, mouza, &q)
		case "4":
			tableAndDonut, err = srv.FormatFisheryHolding(division, district, upazilla, union, mouza, &q)
		case "5":
			tableAndDonut, err = srv.FormatAgriculuralLaborHolding(division, district, upazilla, union, mouza, &q)
		}

		if err != nil {
			log.Error(err)
			srv.frequencyWithError(
				c,
				header,
				footer,
				fmt.Sprintf("unparsable request : %s", err.Error()),
				q)
			return
		}
		srv.frequencyOkWithData(c, header, footer, &q, tableAndDonut)

	})
}

func (srv *Server) frequencyWithError(c *gin.Context, header, footer, alertMsg string, q searchQuery) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	c.HTML(http.StatusOK, "frequency.html", gin.H{
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

func (srv *Server) frequencyOkWithData(c *gin.Context, header, footer string, q *searchQuery, tableAndDonut string) {

	c.HTML(http.StatusOK, "frequency.html", gin.H{
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

func (srv *Server) FormatOccupationOfTheHouseHold(division, district, upazilla, union, mouza string, q *searchQuery) (tableAndDonut string, err error) {
	p := message.NewPrinter(language.English)
	data, err := srv.Db.GetOccupationOfHouseHold(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	tableData := fmt.Sprintf(`
	<tr>
		<td>Agriculture</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>
	<tr>
		<td>Industry</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>
	<tr>
		<td>Service</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>
	<tr>
		<td>Business</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>
	<tr>
		<td>Other</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>		
	<tr>
		<td>Total</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>		
	`,
		p.Sprintf("%d", data.Occ),
		(float64(data.Occ)/float64(data.Total))*100,
		p.Sprintf("%d", data.Occ2),
		(float64(data.Occ2)/float64(data.Total))*100,
		p.Sprintf("%d", data.Occ3),
		(float64(data.Occ3)/float64(data.Total))*100,
		p.Sprintf("%d", data.Occ4),
		(float64(data.Occ4)/float64(data.Total))*100,
		p.Sprintf("%d", data.Occ5),
		(float64(data.Occ5)/float64(data.Total))*100,
		p.Sprintf("%d", data.Total),
		float64(100),
	)

	donutData := fmt.Sprintf(`
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
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
				]
			}
		]
	};

	option && myChart.setOption(option);

</script>
`,
		data.Occ, "Agriculture",
		data.Occ2, "Industry",
		data.Occ3, "Service",
		data.Occ4, "Business",
		data.Occ5, "Other")

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	<h4>Result<small> ফলাফল</small></h4>
	<h5>Data for table number : %s</h5>
	<table class="table">
		<thead>
			<tr>
				<th>Household Head Occupation</th>
				<th>Number of household</th>
				<th>Percentage</th>
			</tr>
		</thead>
		<tbody>
			%s
		</tbody>
	</table>
	</div>
	<div class="form-group">
		<div class="col-md-2 col-sm-2 col-xs-12 col-md-offset-3">
			%s
		</div>
	</div>
	`,
		q.TableNumber,
		tableData,
		donutData)

	return
}

func (srv *Server) FormatEducationHouseHoldHead(division, district, upazilla, union, mouza string, q *searchQuery) (tableAndDonut string, err error) {
	p := message.NewPrinter(language.English)
	data, err := srv.Db.GetEducationOfTheHouseholdHead(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	tableData := fmt.Sprintf(`
	<tr>
		<td>No Education</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>
	<tr>
		<td>Class–I</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>
	<tr>
		<td>Class–II</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>
	<tr>
		<td>Class–III</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>
	<tr>
		<td>Class–IV</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>	
	<tr>
		<td>Class–V</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>	
	<tr>
		<td>Class–VI</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>	
	<tr>
		<td>Class–VII</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>	
	<tr>
		<td>Class–VIII</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>					
	<tr>
		<td>Class–IX</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>	
	<tr>
		<td>SCC</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>				
	<tr>
		<td>HSC</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>	
	<tr>
		<td>Bachelor/Equivalent</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>	
	<tr>
		<td>Masters/Equivalent or Higher</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>				
	<tr>
		<td>Total</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>		
	`,
		p.Sprintf("%d", data.NoEducation),
		(float64(data.NoEducation)/float64(data.Total))*100,
		p.Sprintf("%d", data.Class1),
		(float64(data.Class1)/float64(data.Total))*100,
		p.Sprintf("%d", data.Class2),
		(float64(data.Class2)/float64(data.Total))*100,
		p.Sprintf("%d", data.Class3),
		(float64(data.Class3)/float64(data.Total))*100,
		p.Sprintf("%d", data.Class4),
		(float64(data.Class4)/float64(data.Total))*100,
		p.Sprintf("%d", data.Class5),
		(float64(data.Class5)/float64(data.Total))*100,
		p.Sprintf("%d", data.Class6),
		(float64(data.Class6)/float64(data.Total))*100,
		p.Sprintf("%d", data.Class7),
		(float64(data.Class7)/float64(data.Total))*100,
		p.Sprintf("%d", data.Class8),
		(float64(data.Class8)/float64(data.Total))*100,
		p.Sprintf("%d", data.Class9),
		(float64(data.Class9)/float64(data.Total))*100,
		p.Sprintf("%d", data.Ssc),
		(float64(data.Ssc)/float64(data.Total))*100,
		p.Sprintf("%d", data.Hsc),
		(float64(data.Hsc)/float64(data.Total))*100,
		p.Sprintf("%d", data.BachelorEquivalent),
		(float64(data.BachelorEquivalent)/float64(data.Total))*100,
		p.Sprintf("%d", data.MastersEquivalentOrHigher),
		(float64(data.MastersEquivalentOrHigher)/float64(data.Total))*100,
		p.Sprintf("%d", data.Total),
		float64(100),
	)

	donutData := fmt.Sprintf(`
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
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },									
				]
			}
		]
	};

	option && myChart.setOption(option);

</script>
`,
		data.NoEducation, "No Education",
		data.Class1, "Class–I",
		data.Class2, "Class–II",
		data.Class3, "Class–III",
		data.Class4, "Class–IV",
		data.Class5, "Class–V",
		data.Class6, "Class-VI",
		data.Class7, "Class–VII",
		data.Class8, "Class–VIII",
		data.Class9, "Class–IX",
		data.Ssc, "SSC",
		data.Hsc, "HSC",
		data.BachelorEquivalent, "Bachelor/Equivalent",
		data.MastersEquivalentOrHigher, "Masters/Equivalent or Higher")

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	<h4>Result<small> ফলাফল</small></h4>
	<h5>Data for table number : %s</h5>
	<table class="table">
		<thead>
			<tr>
				<th>Household Head Occupation</th>
				<th>Number of household</th>
				<th>Percentage</th>
			</tr>
		</thead>
		<tbody>
			%s
		</tbody>
	</table>
	</div>
	<div class="form-group">
		<div class="col-md-2 col-sm-2 col-xs-12 col-md-offset-3">
			%s
		</div>
	</div>
	`,
		q.TableNumber,
		tableData,
		donutData)

	return
}

func (srv *Server) FormatGenderOfTheHouseholdHead(division, district, upazilla, union, mouza string, q *searchQuery) (tableAndDonut string, err error) {
	p := message.NewPrinter(language.English)
	data, err := srv.Db.GetGenderOfTheHouseholdHead(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	tableData := fmt.Sprintf(`
	<tr>
		<td>Male</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>
	<tr>
		<td>Female</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>
	<tr>
		<td>Hijra</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td>Total</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>		
	`,
		p.Sprintf("%d", data.Male),
		(float64(data.Male)/float64(data.Total))*100,
		p.Sprintf("%d", data.Female),
		(float64(data.Female)/float64(data.Total))*100,
		p.Sprintf("%d", data.Hijra),
		(float64(data.Hijra)/float64(data.Total))*100,
		p.Sprintf("%d", data.Total),
		float64(100),
	)

	donutData := fmt.Sprintf(`
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
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
					{ value: %d, name: '%s' },
				]
			}
		]
	};

	option && myChart.setOption(option);

</script>
`,
		data.Male, "Male",
		data.Female, "Female",
		data.Hijra, "Hijra")

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	<h4>Result<small> ফলাফল</small></h4>
	<h5>Data for table number : %s</h5>
	<table class="table">
		<thead>
			<tr>
				<th>Household Head Occupation</th>
				<th>Number of household</th>
				<th>Percentage</th>
			</tr>
		</thead>
		<tbody>
			%s
		</tbody>
	</table>
	</div>
	<div class="form-group">
		<div class="col-md-2 col-sm-2 col-xs-12 col-md-offset-3">
			%s
		</div>
	</div>
	`,
		q.TableNumber,
		tableData,
		donutData)

	return
}

func (srv *Server) FormatFisheryHolding(division, district, upazilla, union, mouza string, q *searchQuery) (tableAndDonut string, err error) {
	p := message.NewPrinter(language.English)
	data, err := srv.Db.GetFisheryHolding(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	tableData := fmt.Sprintf(`
	<tr>
		<td>%s</td>
		<td>%s%%</td>
	</tr>	
	`,
		p.Sprintf("%d", data.NumberOfFisheryHousehold),
		p.Sprintf("%f", data.Percentage),
	)

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	<h4>Result<small> ফলাফল</small></h4>
	<h5>Data for table number : %s</h5>
	<table class="table">
		<thead>
			<tr>
				<th>Number of Fishery Household</th>
				<th>Percentage</th>
			</tr>
		</thead>
		<tbody>
			%s
		</tbody>
	</table>
	</div>
	`,
		q.TableNumber,
		tableData)

	return
}

func (srv *Server) FormatAgriculuralLaborHolding(division, district, upazilla, union, mouza string, q *searchQuery) (tableAndDonut string, err error) {
	p := message.NewPrinter(language.English)
	data, err := srv.Db.GetAgriculuralLaborHolding(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	tableData := fmt.Sprintf(`
	<tr>
		<td>%s</td>
		<td>%s%%</td>
	</tr>	
	`,
		p.Sprintf("%d", data.NumberOfAgriLaborHouseHold),
		p.Sprintf("%f", data.Percentage),
	)

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	<h4>Result<small> ফলাফল</small></h4>
	<h5>Data for table number : %s</h5>
	<table class="table">
		<thead>
			<tr>
				<th>Number of Agri labor household</th>
				<th>Percentage</th>
			</tr>
		</thead>
		<tbody>
			%s
		</tbody>
	</table>
	</div>
	`,
		q.TableNumber,
		tableData)

	return
}
