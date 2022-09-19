package routes

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func (srv *Server) FormatHouseholdAgricultureEquipementInformation(division, district, upazilla, union, mouza string, q *searchQuery, geoLocation string) (tableAndDonut string, err error) {
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
	
	<table id="datatable-buttons" class="table table-striped">
	<thead>
	<tr>
	<th class="text-wrap" style="width: 500px;">Data for Table Name : %s</th>
	<th></th>
	<th></th>
	<th></th>
	<th></th>
	<th></th>
	</tr>
	
	</thead>
	<tbody>
	<tr>
	<th>Report</th>
	<th>Number of reporting holdings</th>
	<th>Total number</th>
	<th>Number of non-mechanical device</th>
	<th>Number of diesel device</th>
	<th>Number of electrical device</th>
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

func formatNumber(number uint) string {
	p := message.NewPrinter(language.English)
	if number > 0 {
		return p.Sprintf("%d", number)
	} else {
		return ""
	}
}
