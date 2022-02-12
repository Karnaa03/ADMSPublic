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

func (srv *Server) FormatHouseholdHeadInformation(division, district, upazilla, union, mouza string, q *searchQuery) (tableAndDonut string, err error) {
	p := message.NewPrinter(language.English)
	educationOftheHouseholdHead, err := srv.Db.GetHouseholdHeadInformation(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	occupationOfHouseholdHead, err := srv.Db.GetOccupationOfHouseholdHead(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	totalNumberOfHouseholdMembers, err := srv.Db.GetTotalNumberOfHouseholdMembers(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	totalNumberOfHouseholdWorkers, err := srv.Db.GetTotalNumberOfHouseholdWorkers(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	totalNumberOfHouseholdWorkers1014, err := srv.Db.GetTotalNumberOfHouseholdWorkers1014(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	totalNumberOfHouseholdWorkers15plus, err := srv.Db.GetTotalNumberOfHouseholdWorkers15plus(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	donutEducationOftheHouseholdHead := fmt.Sprintf(`
		<div id="donutEducationOftheHouseholdHead" style="width: 600px;height:400px; align:center" class="x_content"></div>
		<script type="text/javascript">
		var chartDom = document.getElementById('donutEducationOftheHouseholdHead');
		var donutEducationOftheHouseholdHead = echarts.init(chartDom);
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
					name: 'Education of the household head',
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
					]
				}
			]
		};

		option && donutEducationOftheHouseholdHead.setOption(option);

	</script>
	`,
		educationOftheHouseholdHead.NoEducation, "No education",
		educationOftheHouseholdHead.Class_I_V, "Class-I-V",
		educationOftheHouseholdHead.Class_VI_IX, "Class VI-IX",
		educationOftheHouseholdHead.SccPassed, "SSC Passed",
		educationOftheHouseholdHead.HscPassed, "HSC Passed",
		educationOftheHouseholdHead.DegreePassed, "Degree Passed",
		educationOftheHouseholdHead.MasterPassed, "Master Passed",
	)

	donutOccupationOfHouseholdHead := fmt.Sprintf(`
	<div id="donutOccupationOfHouseholdHead" style="width: 600px;height:400px; align:center" class="x_content"></div>
	<script type="text/javascript">
	var chartDom = document.getElementById('donutOccupationOfHouseholdHead');
	var OccupationOfHouseholdHeadChart = echarts.init(chartDom);
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
				name: 'Occupation of household head',
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

	option && OccupationOfHouseholdHeadChart.setOption(option);

</script>
`,
		occupationOfHouseholdHead.Agriculture, "Agriculture",
		occupationOfHouseholdHead.Industry, "Industry",
		occupationOfHouseholdHead.Service, "Service",
		occupationOfHouseholdHead.Business, "Business",
		occupationOfHouseholdHead.Others, "Others",
	)

	donutTotalNumberOfHouseholdMembers := fmt.Sprintf(`
	<div id="donutTotalNumberOfHouseholdMembers" style="width: 600px;height:400px; align:center" class="x_content"></div>
	<script type="text/javascript">
	var chartDom = document.getElementById('donutTotalNumberOfHouseholdMembers');
	var OccupationOfHouseholdHeadChart = echarts.init(chartDom);
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
				name: 'Total number of household members',
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

	option && OccupationOfHouseholdHeadChart.setOption(option);

</script>
`,
		totalNumberOfHouseholdMembers.Male, "Male",
		totalNumberOfHouseholdMembers.Female, "Female",
		totalNumberOfHouseholdMembers.Hijra, "Hijra",
	)

	donutTotalNumberOfHouseholdMWorker := fmt.Sprintf(`
	<div id="donutTotalNumberOfHouseholdMWorker" style="width: 600px;height:400px; align:center" class="x_content"></div>
	<script type="text/javascript">
	var chartDom = document.getElementById('donutTotalNumberOfHouseholdMWorker');
	var OccupationOfHouseholdHeadChart = echarts.init(chartDom);
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
				name: 'Total number of household agricultural worker',
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

	option && OccupationOfHouseholdHeadChart.setOption(option);

</script>
`,
		totalNumberOfHouseholdWorkers.Male, "Male",
		totalNumberOfHouseholdWorkers.Female, "Female",
		totalNumberOfHouseholdWorkers.Hijra, "Hijra",
	)

	donutTotalNumberOfHouseholdMWorker1014 := fmt.Sprintf(`
	<div id="donutTotalNumberOfHouseholdMWorker1014" style="width: 600px;height:400px; align:center" class="x_content"></div>
	<script type="text/javascript">
	var chartDom = document.getElementById('donutTotalNumberOfHouseholdMWorker1014');
	var OccupationOfHouseholdHeadChart = echarts.init(chartDom);
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
				name: 'Total number of household agricultural worker (Age: 10 – 14)',
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

	option && OccupationOfHouseholdHeadChart.setOption(option);

</script>
`,
		totalNumberOfHouseholdWorkers1014.Male, "Male",
		totalNumberOfHouseholdWorkers1014.Female, "Female",
		totalNumberOfHouseholdWorkers1014.Hijra, "Hijra",
	)

	donutTotalNumberOfHouseholdMWorker15plus := fmt.Sprintf(`
	<div id="donutTotalNumberOfHouseholdMWorker15plus" style="width: 600px;height:400px; align:center" class="x_content"></div>
	<script type="text/javascript">
	var chartDom = document.getElementById('donutTotalNumberOfHouseholdMWorker15plus');
	var OccupationOfHouseholdHeadChart = echarts.init(chartDom);
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
				name: 'Total number of household agricultural worker (Age: 15 plus)',
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

	option && OccupationOfHouseholdHeadChart.setOption(option);

</script>
`,
		totalNumberOfHouseholdWorkers15plus.Male, "Male",
		totalNumberOfHouseholdWorkers15plus.Female, "Female",
		totalNumberOfHouseholdWorkers15plus.Hijra, "Hijra",
	)

	tableData := fmt.Sprintf(`
	<tr>
		<td rowspan="8" scope="rowgroup"><b>Education of the household head</b></td>
		<td>No education</td>
		<td>%s</td>
		<td>%.2f%%</td>
		<td rowspan="8" scope="rowgroup">	
			%s				
		</td>
	</tr>
	<tr>
		<td>Class-I-V</td>
		<td>%s</td>
		<td>%.2f</td>
		
	</tr>
	<tr>
		<td>Class VI-IX</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>
	<tr>
		<td>SSC Passed</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>		
	<tr>
		<td>HSC Passed</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>	
	<tr>
		<td>Degree Passed</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>	
	<tr>
		<td>Master Passed</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>	
	<tr>
		<td>Total</td>
		<td>%s</td>
		<td>100%%</td>
	</tr>	
	<tr>
		<td rowspan="6" scope="rowgroup"><b>Occupation of household head</b></td>
		<td>Agriculture</td>
		<td>%s</td>
		<td>%.2f%%</td>
		<td rowspan="6" scope="rowgroup">		
			%s
		</td>
	</tr>
	<tr>
		<td>Industry</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td>Service</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td>Business</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td>Others</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>
	<tr>
		<td>Total</td>
		<td>%s</td>
		<td>%.2f%%</td>
	</tr>		
	<tr>
		<td rowspan="4" scope="rowgroup"><b>Total number of household members</b></td>
		<td>Male</td>
		<td>%s</td>
		<td>%.2f%%</td>
		<td rowspan="4" scope="rowgroup">		
			%s
		</td>
	</tr>	
	<tr>
		<td>Female</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td>Hijra</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td>Total</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>		
	<tr>
		<td rowspan="4" scope="rowgroup"><b>Total number of household agricultural worker</b></td>
		<td>Male</td>
		<td>%s</td>
		<td>%.2f%%</td>
		<td rowspan="4" scope="rowgroup">		
			%s
		</td>
	</tr>	
	<tr>
		<td>Female</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td>Hijra</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td>Total</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td rowspan="4" scope="rowgroup"><b>Total number of household agricultural worker (Age: 10 – 14)</b></td>
		<td>Male</td>
		<td>%s</td>
		<td>%.2f%%</td>
		<td rowspan="4" scope="rowgroup">		
			%s
		</td>
	</tr>	
	<tr>
		<td>Female</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td>Hijra</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td>Total</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td rowspan="4" scope="rowgroup"><b>Total number of household agricultural worker (Age: 15 plus)</b></td>
		<td>Male</td>
		<td>%s</td>
		<td>%.2f%%</td>
		<td rowspan="4" scope="rowgroup">		
			%s
		</td>
	</tr>	
	<tr>
		<td>Female</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td>Hijra</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>	
	<tr>
		<td>Total</td>
		<td>%s</td>
		<td>%.2f</td>
	</tr>							
	`,
		p.Sprintf("%d", educationOftheHouseholdHead.NoEducation),
		(float64(educationOftheHouseholdHead.NoEducation)/float64(educationOftheHouseholdHead.TotalEducation))*100,
		donutEducationOftheHouseholdHead,
		p.Sprintf("%d", educationOftheHouseholdHead.Class_I_V),
		(float64(educationOftheHouseholdHead.Class_I_V)/float64(educationOftheHouseholdHead.TotalEducation))*100,
		p.Sprintf("%d", educationOftheHouseholdHead.Class_VI_IX),
		(float64(educationOftheHouseholdHead.Class_VI_IX)/float64(educationOftheHouseholdHead.TotalEducation))*100,
		p.Sprintf("%d", educationOftheHouseholdHead.SccPassed),
		(float64(educationOftheHouseholdHead.SccPassed)/float64(educationOftheHouseholdHead.TotalEducation))*100,
		p.Sprintf("%d", educationOftheHouseholdHead.HscPassed),
		(float64(educationOftheHouseholdHead.HscPassed)/float64(educationOftheHouseholdHead.TotalEducation))*100,
		p.Sprintf("%d", educationOftheHouseholdHead.DegreePassed),
		(float64(educationOftheHouseholdHead.DegreePassed)/float64(educationOftheHouseholdHead.TotalEducation))*100,
		p.Sprintf("%d", educationOftheHouseholdHead.MasterPassed),
		(float64(educationOftheHouseholdHead.MasterPassed)/float64(educationOftheHouseholdHead.TotalEducation))*100,
		p.Sprintf("%d", educationOftheHouseholdHead.TotalEducation),

		p.Sprintf("%d", occupationOfHouseholdHead.Agriculture),
		(float64(occupationOfHouseholdHead.Agriculture)/float64(occupationOfHouseholdHead.Total))*100,
		donutOccupationOfHouseholdHead,
		p.Sprintf("%d", occupationOfHouseholdHead.Industry),
		(float64(occupationOfHouseholdHead.Industry)/float64(occupationOfHouseholdHead.Total))*100,
		p.Sprintf("%d", occupationOfHouseholdHead.Service),
		(float64(occupationOfHouseholdHead.Service)/float64(occupationOfHouseholdHead.Total))*100,
		p.Sprintf("%d", occupationOfHouseholdHead.Business),
		(float64(occupationOfHouseholdHead.Business)/float64(occupationOfHouseholdHead.Total))*100,
		p.Sprintf("%d", occupationOfHouseholdHead.Others),
		(float64(occupationOfHouseholdHead.Others)/float64(occupationOfHouseholdHead.Total))*100,
		p.Sprintf("%d", occupationOfHouseholdHead.Total),
		(float64(occupationOfHouseholdHead.Total)/float64(occupationOfHouseholdHead.Total))*100,

		p.Sprintf("%d", totalNumberOfHouseholdMembers.Male),
		(float64(totalNumberOfHouseholdMembers.Male)/float64(totalNumberOfHouseholdMembers.Total))*100,
		donutTotalNumberOfHouseholdMembers,
		p.Sprintf("%d", totalNumberOfHouseholdMembers.Female),
		(float64(totalNumberOfHouseholdMembers.Female)/float64(totalNumberOfHouseholdMembers.Total))*100,
		p.Sprintf("%d", totalNumberOfHouseholdMembers.Hijra),
		(float64(totalNumberOfHouseholdMembers.Hijra)/float64(totalNumberOfHouseholdMembers.Total))*100,
		p.Sprintf("%d", totalNumberOfHouseholdMembers.Total),
		(float64(totalNumberOfHouseholdMembers.Total)/float64(totalNumberOfHouseholdMembers.Total))*100,

		p.Sprintf("%d", totalNumberOfHouseholdWorkers.Male),
		(float64(totalNumberOfHouseholdWorkers.Male)/float64(totalNumberOfHouseholdWorkers.Total))*100,
		donutTotalNumberOfHouseholdMWorker,
		p.Sprintf("%d", totalNumberOfHouseholdWorkers.Female),
		(float64(totalNumberOfHouseholdWorkers.Female)/float64(totalNumberOfHouseholdWorkers.Total))*100,
		p.Sprintf("%d", totalNumberOfHouseholdWorkers.Hijra),
		(float64(totalNumberOfHouseholdWorkers.Hijra)/float64(totalNumberOfHouseholdWorkers.Total))*100,
		p.Sprintf("%d", totalNumberOfHouseholdWorkers.Total),
		(float64(totalNumberOfHouseholdWorkers.Total)/float64(totalNumberOfHouseholdWorkers.Total))*100,

		p.Sprintf("%d", totalNumberOfHouseholdWorkers1014.Male),
		(float64(totalNumberOfHouseholdWorkers1014.Male)/float64(totalNumberOfHouseholdWorkers1014.Total))*100,
		donutTotalNumberOfHouseholdMWorker1014,
		p.Sprintf("%d", totalNumberOfHouseholdWorkers1014.Female),
		(float64(totalNumberOfHouseholdWorkers1014.Female)/float64(totalNumberOfHouseholdWorkers1014.Total))*100,
		p.Sprintf("%d", totalNumberOfHouseholdWorkers1014.Hijra),
		(float64(totalNumberOfHouseholdWorkers1014.Hijra)/float64(totalNumberOfHouseholdWorkers1014.Total))*100,
		p.Sprintf("%d", totalNumberOfHouseholdWorkers1014.Total),
		(float64(totalNumberOfHouseholdWorkers1014.Total)/float64(totalNumberOfHouseholdWorkers1014.Total))*100,

		p.Sprintf("%d", totalNumberOfHouseholdWorkers15plus.Male),
		(float64(totalNumberOfHouseholdWorkers15plus.Male)/float64(totalNumberOfHouseholdWorkers15plus.Total))*100,
		donutTotalNumberOfHouseholdMWorker15plus,
		p.Sprintf("%d", totalNumberOfHouseholdWorkers15plus.Female),
		(float64(totalNumberOfHouseholdWorkers15plus.Female)/float64(totalNumberOfHouseholdWorkers15plus.Total))*100,
		p.Sprintf("%d", totalNumberOfHouseholdWorkers15plus.Hijra),
		(float64(totalNumberOfHouseholdWorkers15plus.Hijra)/float64(totalNumberOfHouseholdWorkers15plus.Total))*100,
		p.Sprintf("%d", totalNumberOfHouseholdWorkers15plus.Total),
		(float64(totalNumberOfHouseholdWorkers15plus.Total)/float64(totalNumberOfHouseholdWorkers15plus.Total))*100,
	)

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	<h4>Result<small> ফলাফল</small></h4>
	<h5>Data for table number : %s</h5>
	<table class="table">
		<thead>
			<tr>
				<th>Report</th>
				<th>Indicator</th>
				<th>Total</th>
				<th>Percentage</th>
				<th>Graph</th>
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
