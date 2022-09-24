package routes

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (srv *Server) FormatHouseholdLandInformation(division, district, upazilla, union, mouza string, q *searchQuery, geoLocation string) (tableAndDonut string, err error) {
	p := message.NewPrinter(language.English)
	householdLandInformation, err := srv.Db.GetHouseholdLandInformation(division, district, upazilla, union, mouza)
	if err != nil {
		return "", err
	}

	tableData := ""
	for _, hli := range householdLandInformation {
		tableData += fmt.Sprintf(`
		<tr>
			<td><b>%s</b></td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		</tr>
								
		`,
			hli.Name,
			p.Sprintf("%d", hli.NumberOfReportingHoldings),
			p.Sprintf("%d", hli.NumberOfFarmHoldings),
			p.Sprintf("%.2f", hli.TotalAreaOfOwnLand),
			p.Sprintf("%.2f", hli.TotalFarmHoldingArea),
			p.Sprintf("%.2f", (float64(hli.TotalFarmHoldingArea)/float64(hli.NumberOfFarmHoldings))),
		)
	}

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	Data for table name : %s
	
	<table id="datatable-buttons" class="table table-striped">
	<thead>
	<tr>
	<th>Report</th>
	<th>Number of reporting holdings</th>
	<th>Number of farm holdings</th>
	<th>Total Area (acres)</th>
	<th>Total farm holding area (acres)</th>
	<th>Average area (acres) per farm holding</th>
	</tr>	
	</thead>
	<tbody>
	%s
	</tbody>
	<tfoot>
	<tr>
	  <th>Source: Agriculture Census 2019, Bangladesh Bureau of Statistics</th>
	  <th></th>
	  <th></th>
	</tr>
  </tfoot>
	</table>
	</div>
	<h7>Source: Agriculture Census 2019, Bangladesh Bureau of Statistics.</h7>
	`,
		fmt.Sprintf("%s<br>%s", getTableGenerationName(q.TableNumber), geoLocation),
		tableData)

	return
}
