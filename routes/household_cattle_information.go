package routes

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (srv *Server) FormatHouseholdCattleInformation(division, district, upazilla, union, mouza string, q *searchQuery) (tableAndDonut string, err error) {
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
			p.Sprintf("%.2f", (float64(hli.TotalNumberOfCattle)/float64(hli.NumberOfReportingHoldings))),
		)
	}

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	<h4>Result></h4>
	<h5>Data for table number : %s</h5>
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
	`,
		getTableGenerationName(q.TableNumber),
		tableData)

	return
}
