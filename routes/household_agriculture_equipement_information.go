package routes

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (srv *Server) FormatHouseholdAgricultureEquipementInformation(division, district, upazilla, union, mouza string, q *searchQuery) (tableAndDonut string, err error) {
	householdLandInformation, err := srv.Db.GetHouseholdAgricultureEquipement(division, district, upazilla, union, mouza)
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
			formatNumber(hli.NumberOfReportingHoldings),
			formatNumber(hli.TotalNumber),
			formatNumber(hli.NumberOfNonMechanicalDevice),
			formatNumber(hli.NumberOfDieselDevice),
			formatNumber(hli.NumberOfElectricalDevice),
		)
	}

	tableAndDonut = fmt.Sprintf(`
	<div class="x_content">
	<h4>Result<small> </small></h4>
	<h5>Data for table name : %s</h5>
	<h7>Source: Bangladesh Bureau of Statistics. Report produced by Agriculture (Crops, Fisheries and Livestock) Census 2018 Project.</h7>
	<table id="datatable-buttons" class="table table-striped">
		<thead>
			<tr>
				<th>Report</th>
				<th>Number of reporting holdings</th>
				<th>Total number</th>
				<th>Number of non-mechanical device</th>
				<th>Number of diesel device</th>
				<th>Number of electrical device</th>
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

func formatNumber(number uint) string {
	p := message.NewPrinter(language.English)
	if number > 0 {
		return p.Sprintf("%d", number)
	} else {
		return ""
	}
}
