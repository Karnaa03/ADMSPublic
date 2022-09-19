package routes

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (srv *Server) FormatHouseholdPoultryInformation(division, district, upazilla, union, mouza string, q *searchQuery, geoLocation string) (tableAndDonut string, err error) {
	p := message.NewPrinter(language.English)
	householdLandInformation, err := srv.Db.GetHouseholdPoultryInformation(division, district, upazilla, union, mouza)
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
			p.Sprintf("%d", hli.TotalNumberOfPoultry),
			p.Sprintf("%.4f", hli.AverageTypeOfPoultryPerHolding),
		)
	}

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	<h4>Result<small> </small></h4>
	
	<table id="datatable-buttons" class="table table-striped">
	<thead>
	<tr>
	<th class="text-wrap" style="width: 500px;">Data for table name : %s</th>
	<th></th>
	<th></th>
	<th></th>
	</tr>
	
	</thead>
	<tbody>
	<tr>
	<th>Poultry</th>
	<th>Number of reporting holdings</th>
	<th>Total number of poultry</th>
	<th>Average type of poultry per holding</th>
	</tr>
	%s
	</tbody>
	</table>
	</div>
	<h7>Source: Agriculture Census 2019, Bangladesh Bureau of Statistics.</h7>
	`,
		fmt.Sprintf("%s  %s", getTableGenerationName(q.TableNumber), geoLocation),
		tableData)

	return
}
