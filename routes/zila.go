package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	log "github.com/sirupsen/logrus"

	gin_oidc "git.solutions.im/Solutions.IM/ginOidc"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/model"
	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

type Test struct {
	Number    string `form:"BookletNumber"`
	GeoCodeID string `form:"GeoCode"`
	Size      uint   `form:"BookletSize"`
}

const checked = "checked"

/*
Only the booklet number is unique
*/
func (srv *Server) zila(footer string) {
	srv.router.GET("/production/zila.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		type s1query struct {
			Size            string `form:"size"`
			BookletNumber   string `form:"number"`
			GeoCodeNumbers  string `form:"geo"`
			BookletToDelete string `form:"delete"`
		}
		var q s1query
		err := c.ShouldBind(&q)
		if err != nil {
			log.Error(err)
			srv.stage1WithError(c, header, footer, "unprocessable query", "", "", "")
		}
		switch q.BookletToDelete {
		case "":
			srv.OkWithData(c, header, footer, q.Size, q.BookletNumber, q.GeoCodeNumbers)
		default:
			log.Infof("delete booklet number %s", q.BookletToDelete)
			err := srv.DeleteBooklet(q.BookletToDelete)
			if err != nil {
				log.Error(err)
			}
			srv.OkWithData(c, header, footer, q.Size, "", q.GeoCodeNumbers)
		}
	})

	srv.router.POST("/production/zila.html", func(c *gin.Context) {
		size := c.PostForm("BookletSize")
		geoCode := c.PostForm("GeoCode")
		number := c.PostForm("BookletNumber")
		header, _ := templates.RenderHeader(c)

		id, err := gin_oidc.GetIdentity(c)
		if err != nil {
			srv.stage1WithError(c, header, footer, fmt.Sprintf("Who are you ??? : %s", err), number, geoCode, size)
		}

		var booklet model.Booklet
		booklet.GeoCodeID = geoCode
		booklet.Number = number

		err = model.RegisterNewBooklet(number, geoCode, size, id, &srv.Db)
		if err != nil {
			pgErr, ok := err.(pg.Error)
			if ok && pgErr.IntegrityViolation() {
				switch {
				case strings.Contains(pgErr.Error(), "duplicate key value violates unique constraint \"booklets_number_key\""):
					srv.stage1WithError(c, header, footer, "This booklet number has already been registered", number, geoCode, size)
				case strings.Contains(pgErr.Error(), "insert or update on table \"booklets\" violates foreign key constraint \"booklets_geo_code_id_fkey\""):
					srv.stage1WithError(c, header, footer, "This Geo code is unknown", number, geoCode, size)
				default:
					srv.stage1WithError(c, header, footer, pgErr.Error(), number, geoCode, size)
				}
			} else {
				srv.stage1WithError(c, header, footer, err.Error(), number, geoCode, size)
				// srv.RecordAction(c, "record", &booklet)
				log.Error(err)
			}
		}
		srv.OkWithData(c, header, footer, size, "", geoCode[0:14])
	})

}

func (srv *Server) OkWithData(c *gin.Context, header, footer, previousSize, bookletNumber, geoCodeNumber string) {
	data := gin.H{
		"Header":         template.HTML(header),
		"Footer":         template.HTML(footer),
		"TableData":      template.HTML(srv.GetChecked()),
		"BookletBarCode": bookletNumber,
		"BookletQRCode":  geoCodeNumber,
	}
	switch previousSize {
	case "100":
		data["R100"] = checked
	case "50":
		data["R50"] = checked
	case "25":
		data["R25"] = checked
	default:
		data["R100"] = checked
	}
	c.HTML(http.StatusOK, "zila.html", data)
}

func (srv *Server) stage1WithError(c *gin.Context, header, footer, alertMsg, bookletNumber, bookletGeoCode, previousSize string) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	// name := gin_oidc.GetValue(c, "name")

	data := gin.H{
		"Header":         template.HTML(header),
		"Footer":         template.HTML(footer),
		"AlertQrCode":    template.HTML(alert),
		"BookletBarCode": bookletNumber,
		"BookletQRCode":  bookletGeoCode,
		"TableData":      template.HTML(srv.GetChecked()),
	}
	switch previousSize {
	case "100":
		data["R100"] = checked
	case "50":
		data["R50"] = checked
	case "25":
		data["R25"] = checked
	default:
		data["R100"] = checked
	}
	c.HTML(http.StatusOK, "zila.html", data)
}

func (srv *Server) GetChecked() (data string) {
	var booklets []model.Booklet
	err := srv.Db.Conn.Model(&booklets).
		Where("status = ?", "registered").
		Relation("GeoCode").
		Limit(50).
		Order("registered_on DESC").
		Select()
	if err != nil {
		log.Error(err)
		return
	}
	for _, booklet := range booklets {
		data += fmt.Sprintf(`
		<tr>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%s</td>
			<td>%d</td>
			<td>
				<a href="#" data-href="/production/zila.html?delete=%s" data-toggle="modal" data-target="#confirm-delete">
					<center>
						<i class="fa fa-trash"></i>
					</center>
				</a>
			</td>
		</tr>
	`, booklet.RegisteredOn.Format(time.RFC1123),
			booklet.Number,
			booklet.GeoCodeID,
			booklet.GeoCode.NameDistrict,
			booklet.GeoCode.NameMouza,
			booklet.Size,
			booklet.Number)
	}
	return
}

func (srv *Server) DeleteBooklet(id string) (err error) {
	booklet := model.Booklet{
		Number: id,
	}
	event := model.Event{
		BookletNumber: id,
	}
	_, err = srv.Db.Conn.Model(&event).Where("booklet_number = ?", id).Delete()
	if err != nil {
		return
	}
	_, err = srv.Db.Conn.Model(&booklet).Where("number = ?", id).Delete()
	return err
}
