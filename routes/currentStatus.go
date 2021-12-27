package routes

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	gin_oidc "git.solutions.im/Solutions.IM/ginOidc"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

type stats struct {
	tableName                  struct{} `pg:"select:stats"`
	RegisteredToday            int
	RegisteredYesterday        int
	InbatchToday               int
	InbatchYesterday           int
	IncuttingstationToday      int
	IncuttingstationYesterday  int
	InprescanningToday         int
	InprescanningYesterday     int
	InscanningstationToday     int
	InscanningstationYesterday int
	ArchivedToday              int
	ArchivedYesterday          int
}

type totalPerStatus struct {
	Status     string
	Count      int
	Percentage float64
}

type topPerformer struct {
	Station  string
	Operator string
	Total    int
}

type topPerformers []topPerformer

// topPerformerMap the key is topPerformer.Station
type topPerformerMap map[string][]topPerformer

func (tp topPerformers) ToMap() topPerformerMap {
	topPerfMap := topPerformerMap{}
	for _, performer := range tp {
		if v, ok := topPerfMap[performer.Station]; ok {
			v = append(v, performer)
			topPerfMap[performer.Station] = v
		} else {
			topPerfMap[performer.Station] = []topPerformer{performer}
		}
	}
	return topPerfMap
}

func (srv *Server) currentStatus(footer string) {
	srv.router.GET("/reporting/currentStatus.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		srv.currentStatusOkWithData(c, header, footer, "", "")
	})

}

func (srv *Server) currentStatusOkWithData(c *gin.Context, header, footer, shelfNumber, crateNumber string) {
	name := gin_oidc.GetValue(c, "name")
	s, err := srv.getStats()
	if err != nil {
		log.Error(err)
	}
	total, err := srv.total()
	if err != nil {
		log.Error(err)
	}
	filledTotal := computePercentage(sortTotal(total))

	overallTotal, err := srv.overallTotal()
	if err != nil {
		log.Error(err)
	}
	filledOverallTotal := computePercentage(sortTotal(overallTotal))

	topPerf, err := srv.topPerformer()
	if err != nil {
		log.Error(err)
	}

	topPerfMap := topPerf.ToMap()

	c.HTML(http.StatusOK, "currentStatus.html", gin.H{
		"Name":                      name,
		"Header":                    template.HTML(header),
		"Footer":                    template.HTML(footer),
		"Registered":                s.RegisteredToday,
		"RegisteredYesterday":       s.RegisteredYesterday,
		"InBatch":                   s.InbatchToday,
		"InBatchYesterday":          s.InbatchYesterday,
		"InCuttingStation":          s.IncuttingstationToday,
		"InCuttingStationYesterday": s.IncuttingstationYesterday,
		"InPreScanning":             s.InprescanningToday,
		"InPreScanningYesterday":    s.InprescanningYesterday,
		"InScanning":                s.InscanningstationToday,
		"InScanningYesterday":       s.InscanningstationYesterday,
		"Archived":                  s.ArchivedToday,
		"ArchivedYesterday":         s.ArchivedYesterday,
		"Total":                     filledTotal,
		"OverallTotal":              filledOverallTotal,
		"TopPerformer":              topPerfMap,
	})
}

func sortTotal(total []totalPerStatus) []totalPerStatus {
	sortedTotal := make([]totalPerStatus, len(total))
	counter := 0
	for _, order := range []string{"registered", "inBatch", "inCuttingStation", "inPreScanning", "inScanningStation", "scanned", "inIceBox", "archived"} {
		for _, t1 := range total {
			if t1.Status == order {
				sortedTotal[counter] = t1
				counter++
			}
		}
	}
	return sortedTotal
}

func (srv *Server) getStats() (s stats, err error) {
	s = stats{}
	_, err = srv.Db.Conn.Query(&s, "select * from stats")
	return
}

func variance(today, yesterday int) (percentage float64, color string, upDown string) {
	if today > 0 {
		percentage = float64((today-yesterday)/today) * 100
	} else {
		percentage = -100
	}
	if percentage > 0 {
		color = "green"
		upDown = "fa-sort-asc"
	} else {
		color = "red"
		upDown = "fa-sort-desc"
	}
	return
}

func (srv *Server) total() (status []totalPerStatus, err error) {
	status = []totalPerStatus{}
	_, err = srv.Db.Conn.Query(&status, `
select b.status, count(b.number)
from booklets b
group by b.status;`)
	return
}

func (srv *Server) overallTotal() (status []totalPerStatus, err error) {
	status = []totalPerStatus{}
	_, err = srv.Db.Conn.Query(&status, `
select 'registered' as status , count(b.registered_on)
from booklets b
where b.registered_on is not null
union all
select 'scanned' as status, count(b.scanned_on)
from booklets b
where b.scanned_on is not null
union all
select 'archived' as status, count(b.archived_on)
from booklets b
where b.archived_on is not null
;`)
	return
}

func (srv *Server) topPerformer() (top topPerformers, err error) {
	top = topPerformers{}
	_, err = srv.Db.Conn.Query(&top, `
select station, operator, total from (
select
       e.event_type as station,
       i.full_name as operator,
       count(e.event_type) as total,
       row_number() over (partition by e.event_type order by count(e.event_type) desc ) as rownum
from events e
         left join identities i on e.identity_id = i.id
where e.event_type in ('register','addInBatch','moveToCuttingStation','moveToPreScanning','moveToScanStation','archive')
group by i.full_name, e.event_type
order by total desc) tmp where rownum < 4;
	`)
	return
}

func computePercentage(t []totalPerStatus) (result []totalPerStatus) {
	max := 0
	for _, v := range t {
		if v.Count > max {
			max = v.Count
		}
	}
	result = []totalPerStatus{}
	for _, v := range t {
		result = append(result, totalPerStatus{
			Status:     v.Status,
			Count:      v.Count,
			Percentage: float64(v.Count) / float64(max) * 100,
		})
	}
	return
}
