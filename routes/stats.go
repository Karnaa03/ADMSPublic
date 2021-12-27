package routes

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type MonthlyStats struct {
	Date  int64
	Stats DailyStats
}

type DailyStats struct {
	Registered   int
	AddedInBatch int
	Cut          int
	PreScanned   int
	Scanned      int
	Archived     int
}

type Stats struct {
	Date       int64 `json:"date"`
	NbBooklets int   `json:"qty"`
}

func (srv *Server) registeredStats() {
	srv.router.GET("/reporting/registered_stats", func(c *gin.Context) {
		var registered []Stats
		_, err := srv.Db.Conn.Query(&registered, `
select extract(epoch from date(date))*1000 as date,
       count(b.registered_on)     as nb_booklets
from generate_series(date(now() - interval '30 days'), date(now()), '1 day') as date
         left join booklets as b on date(date.date) = date(b.registered_on)
group by date(date.date)
order by date;`)
		if err != nil {
			log.Error(err)
		}
		c.JSON(200, registered)
	})
}

func (srv *Server) inBatchStats() {
	srv.router.GET("/reporting/inbatch_stats", func(c *gin.Context) {
		var addedInBatch []Stats
		_, err := srv.Db.Conn.Query(&addedInBatch, `
select extract(epoch from date(date))*1000 as date,
       count(b.added_in_batch_on)     as nb_booklets
from generate_series(date(now() - interval '30 days'), date(now()), '1 day') as date
         left join booklets as b on date(date.date) = date(b.added_in_batch_on)
group by date(date.date)
order by date;`)
		if err != nil {
			log.Error(err)
		}
		c.JSON(200, addedInBatch)
	})
}

func (srv *Server) cutStats() {
	srv.router.GET("/reporting/cut_stats", func(c *gin.Context) {
		var cut []Stats
		_, err := srv.Db.Conn.Query(&cut, `
select extract(epoch from date(date))*1000 as date,
       count(b.cut_on)     as nb_booklets
from generate_series(date(now() - interval '30 days'), date(now()), '1 day') as date
         left join booklets as b on date(date.date) = date(b.cut_on)
group by date(date.date)
order by date;`)
		if err != nil {
			log.Error(err)
		}
		c.JSON(200, cut)
	})
}

func (srv *Server) preparedStats() {
	srv.router.GET("/reporting/prepared_stats", func(c *gin.Context) {
		var prepared []Stats
		_, err := srv.Db.Conn.Query(&prepared, `
select extract(epoch from date(date))*1000 as date,
       count(b.prepared_on)     as nb_booklets
from generate_series(date(now() - interval '30 days'), date(now()), '1 day') as date
         left join booklets as b on date(date.date) = date(b.prepared_on)
group by date(date.date)
order by date;`)
		if err != nil {
			log.Error(err)
		}
		c.JSON(200, prepared)
	})
}

func (srv *Server) scannStats() {
	srv.router.GET("/reporting/scann_stats", func(c *gin.Context) {
		var scanned []Stats
		_, err := srv.Db.Conn.Query(&scanned, `
select extract(epoch from date(date))*1000 as date,
       count(b.scanned_on)     as nb_booklets
from generate_series(date(now() - interval '30 days'), date(now()), '1 day') as date
         left join booklets as b on date(date.date) = date(b.scanned_on)
group by date(date.date)
order by date;`)
		if err != nil {
			log.Error(err)
		}
		c.JSON(200, scanned)
	})
}

func (srv *Server) archivedStats() {
	srv.router.GET("/reporting/archived_stats", func(c *gin.Context) {
		var archivedOn []Stats
		_, err := srv.Db.Conn.Query(&archivedOn, `
select extract(epoch from date(date))*1000 as date,
       count(b.archived_on)     as nb_booklets
from generate_series(date(now() - interval '30 days'), date(now()), '1 day') as date
         left join booklets as b on date(date.date) = date(b.archived_on)
group by date(date.date)
order by date;`)
		if err != nil {
			log.Error(err)
		}
		c.JSON(200, archivedOn)
	})
}
