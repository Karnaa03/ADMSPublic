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
	<h4>Result<small> </small></h4>
	<h5>Data for table name : %s</h5>
	<table id="datatable-buttons" class="table table-striped">
	<thead>
	<tr>
	<th>Report</th>
	<th>Number of reporting holdings</th>
	<th>Number of farm holdings</th>
	<th>Total Area (acres) of own land</th>
	<th>Total farm holding area</th>
	<th>Average area (acres) per farm holding</th>
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
