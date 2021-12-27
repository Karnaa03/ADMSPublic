package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	ginoidc "git.solutions.im/Solutions.IM/ginOidc"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/model"
	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

type archivedBox struct {
	ArchivedBoxNumber         string
	BookletNumber             string
	GeoCodeId                 string
	BookletSize               int
	District                  string
	Mouza                     string
	WarehouseRowNumber        int
	WarehouseShelfNumber      int
	WarehouseShelfLevelNumber int
}

func (srv *Server) registerArchiveBox(footer string) {
	srv.router.GET("/warehouse/registerArchiveBox.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		srv.registerArchiveBoxOkWithData(c, header, footer, "", 0, 0, 0)

	})

	srv.router.POST("/warehouse/registerArchiveBox.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		var registrationRequest model.ArchiveBox
		err := c.ShouldBind(&registrationRequest)
		if err != nil {
			log.Error(err)
			srv.registerArchiveBoxWithError(c, header, footer, err.Error(), "", 0, 0, 0)
			return
		}
		srv.registerArchiveBoxInWarehouse(c, header, footer, registrationRequest)
	})
}

func (srv *Server) registerArchiveBoxInWarehouse(c *gin.Context, header, footer string, registrationRequest model.ArchiveBox) {
	id, err := ginoidc.GetIdentity(c)
	if err != nil {
		log.Error(err)
		srv.registerArchiveBoxWithError(
			c,
			header,
			footer,
			err.Error(),
			registrationRequest.Number,
			registrationRequest.WarehouseRowNumber,
			registrationRequest.WarehouseShelfNumber,
			registrationRequest.WarehouseShelfLevelNumber,
		)
		return
	}

	archiveBox, err := srv.Db.GetArchiveBox(registrationRequest.Number)
	if err != nil {
		log.Error(err)
		srv.registerArchiveBoxWithError(
			c,
			header,
			footer,
			err.Error(),
			registrationRequest.Number,
			registrationRequest.WarehouseRowNumber,
			registrationRequest.WarehouseShelfNumber,
			registrationRequest.WarehouseShelfLevelNumber,
		)
	}
	err = archiveBox.RegisterInWarehouse(id, &srv.Db, registrationRequest)
	if err != nil {
		log.Error(err)
		srv.registerArchiveBoxWithError(
			c,
			header,
			footer,
			err.Error(),
			registrationRequest.Number,
			registrationRequest.WarehouseRowNumber,
			registrationRequest.WarehouseShelfNumber,
			registrationRequest.WarehouseShelfLevelNumber,
		)
	}
	abNumber := archiveBox.Number
	srv.registerArchiveBoxOkWithData(c, header, footer, abNumber, 0, 0, 0)
}

func (srv *Server) registerArchiveBoxOkWithData(c *gin.Context, header, footer, archiveBoxNumber string, warehouseRowNumber, warehouseShelfNumber, warehouseShelfLevelNumber int) {
	name := ginoidc.GetValue(c, "name")

	tableData, lastArchiveBoxNumber := srv.getArchiveBox(archiveBoxNumber)

	c.HTML(http.StatusOK, "registerArchiveBox.html", gin.H{
		"Name":             name,
		"Header":           template.HTML(header),
		"Footer":           template.HTML(footer),
		"ArchiveBoxNumber": lastArchiveBoxNumber,
		"TableData":        template.HTML(tableData),
	})
}

func (srv *Server) registerArchiveBoxWithError(c *gin.Context, header, footer, alertMsg string, archiveBoxNumber string, warehouseRowNumber, warehouseShelfNumber, warehouseShelfLevelNumber int) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := ginoidc.GetValue(c, "name")

	var zero = func(input int) string {
		if input == 0 {
			return ""
		} else {
			return strconv.Itoa(input)
		}
	}

	tableData, _ := srv.getArchiveBox("")

	c.HTML(http.StatusOK, "registerArchiveBox.html", gin.H{
		"Name":                      name,
		"Header":                    template.HTML(header),
		"Footer":                    template.HTML(footer),
		"Alert":                     template.HTML(alert),
		"ArchiveBoxNumber":          archiveBoxNumber,
		"WarehouseRowNumber":        fmt.Sprintf("%02s", zero(warehouseRowNumber)),
		"WarehouseShelfNumber":      fmt.Sprintf("%03s", zero(warehouseShelfNumber)),
		"WarehouseShelfLevelNumber": zero(warehouseShelfLevelNumber),
		"TableData":                 template.HTML(tableData),
	})
}

func (srv *Server) getArchiveBox(archivedBoxNumber string) (data, lastArchiveBoxNumber string) {
	var archives []archivedBox
	if archivedBoxNumber == "" {
		return
	}

	_, err := srv.Db.Conn.Query(&archives, `
	select
    a.number as archived_box_number,
    b.number as booklet_number,
    b.geo_code_id as geo_code_id,
    b.size as booklet_size,
    g.name_district as district,
    g.name_mouza as mouza,
    a.warehouse_row_number as warehouse_row_number,
    a.warehouse_shelf_number as warehouse_shelf_number,
    a.warehouse_shelf_level_number as warehouse_shelf_level_number
from archive_boxes a,
     booklets b,
     geo_codes g
where b.archive_box_number = a.number
  and b.geo_code_id = g.geocode_id
	and a.number = ?
order by a.number, b.number;`, archivedBoxNumber)
	if err != nil {
		log.Error(err)
		return
	}

	for _, archive := range archives {
		var storageLocaltion string
		if archive.WarehouseRowNumber == 0 && archive.WarehouseShelfNumber == 0 && archive.WarehouseShelfLevelNumber == 0 {
			storageLocaltion = "not yet registered"
		} else {
			storageLocaltion = fmt.Sprintf("Row : %d, Shelf : %d, Level : %d",
				archive.WarehouseRowNumber,
				archive.WarehouseShelfNumber,
				archive.WarehouseShelfLevelNumber)
		}
		data += fmt.Sprintf(`
		<tr>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%d</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
		</tr>
`,
			archive.ArchivedBoxNumber,
			archive.BookletNumber,
			archive.GeoCodeId,
			archive.BookletSize,
			archive.District,
			archive.Mouza,
			storageLocaltion)
	}

	if len(archives) == 1 {
		lastArchiveBoxNumber = archives[0].ArchivedBoxNumber
	}
	return
}
