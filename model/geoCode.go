package model

type GeoCodes struct {
	GeocodeID        string `pg:"type:ltree,pk"`
	Division         uint8
	District         uint8
	Upazilla         uint8
	Union            uint16
	Mouza            uint16
	Villages         *Villages
	CA               uint8
	Rmo              uint8
	NameDivision     string
	NameDistrict     string
	NameUpazilla     string
	NameUnion        string
	NameMouza        string
	NameCountingArea string
	NameRMO          string
}

type Villages struct {
	Villages map[int]Village
}

type Village struct {
	NameEn      string
	NameBengali string
}
