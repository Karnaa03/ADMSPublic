package model

var index = map[string]string{
	"tally_sheet_no": `CREATE INDEX IF NOT EXISTS geo_codes_agregated ON agregateds (geocode);`,
}
