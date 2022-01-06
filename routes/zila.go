package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/ADMSPublic/model"
	"git.solutions.im/XeroxAgriCensus/ADMSPublic/templates"
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
			srv.zilaOkWithData(c, header, footer, "", "")
		default:
			log.Infof("delete booklet number %s", q.BookletToDelete)
			err := srv.DeleteBooklet(q.BookletToDelete)
			if err != nil {
				log.Error(err)
			}
			srv.zilaOkWithData(c, header, footer, "", "")
		}
	})

	srv.router.POST("/production/zila.html", func(c *gin.Context) {
		division := c.PostForm("Division")
		district := c.PostForm("District")
		tableName := c.PostForm("TableName")
		header, _ := templates.RenderHeader(c)

		fmt.Printf("division : %s , district : %s, tableName : %s\n", division, district, tableName)

		srv.zilaOkWithData(c, header, footer, division, district)
	})
}

func (srv *Server) zilaOkWithData(c *gin.Context, header, footer, division, district string) {
	data := gin.H{
		"Header":    template.HTML(header),
		"Footer":    template.HTML(footer),
		"TableData": template.HTML(srv.GetChecked()),
		"Division":  division,
		"District":  district,
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
