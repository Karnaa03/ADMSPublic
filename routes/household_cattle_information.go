package routes

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (srv *Server) FormatHouseholdCattleInformation(division, district, upazilla, union, mouza string, q *searchQuery, geoLocation string) (tableAndDonut string, err error) {
	p := message.NewPrinter(language.English)
	householdLandInformation, err := srv.Db.GetHouseholdCattlenformation(division, district, upazilla, union, mouza)
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
		</tr>
								
		`,
			hli.Name,
			p.Sprintf("%d", hli.NumberOfReportingHoldings),
			p.Sprintf("%d", hli.TotalNumberOfCattle),
			p.Sprintf("%.2f", hli.AverageTypeOfCattlePerHolding),
		)
	}

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	Data for table name : %s
	<table id="datatable-buttons" class="table table-striped">
	<thead>
	<tr>
	<th>Cattle</th>
	<th>Number of reporting holdings</th>
	<th>Total number of cattle</th>
	<th>Average type of cattle per holding</th>
	</tr>
	</thead>
	<tbody>
	
	%s
	</tbody>
 
	</table>
	</div>
	<h7>Source: Agriculture Census 2019, Bangladesh Bureau of Statistics.</h7>
	`,
		fmt.Sprintf("%s<br>%s", getTableGenerationName(q.TableNumber), geoLocation),
		tableData)

	return
}
