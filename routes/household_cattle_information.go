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
			p.Sprintf("%.4f", hli.AverageTypeOfCattlePerHolding),
		)
	}

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	<h4>Result<small> </small></h4>
	<h5>Data for table name : %s</h5>
	<table id="datatable-buttons" class="table table-striped">
	<thead>
	<tr>
	<th>Report</th>
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
	<h7>Source: Bangladesh Bureau of Statistics. Report produced by Agriculture (Crops, Fisheries and Livestock) Census 2018 Project.</h7>
	`,
		fmt.Sprintf("%s GEO CODE : %s", getTableGenerationName(q.TableNumber), geoLocation),
		tableData)

	return
}
