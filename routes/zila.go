package routes

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/ADMSPublic/templates"
)

type Test struct {
	Number    string `form:"BookletNumber"`
	GeoCodeID string `form:"GeoCode"`
	Size      uint   `form:"BookletSize"`
}

/*
Only the booklet number is unique
*/
func (srv *Server) zila(footer string) {
	srv.router.GET("/production/zila.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		type zilaQuery struct {
			Division  string `form:"Division"`
			District  string `form:"District"`
			TableName string `form:"TableName"`
		}
		var q zilaQuery
		err := c.ShouldBind(&q)
		if err != nil {
			log.Error(err)
			srv.zilaWithError(c, header, footer, "unprocessable query", "", "")
		}
		srv.zilaOkWithData(c, header, footer, "", "", "")
	})

	srv.router.GET("/production/getReport", func(c *gin.Context) {
		tlno := c.Query("no")
		content, err := srv.S3.Get(srv.Config.S3Config.Bucket, tlno)
		if err != nil && err.Error() != io.EOF.Error() {
			log.Error(err)
		}
		c.Data(http.StatusOK, "application/pdf", content)
	})

	srv.router.POST("/production/zila.html", func(c *gin.Context) {
		division := strings.Trim(strings.Split(c.PostForm("DivisionNumber"), "-")[0], " ")
		division_num, err := strconv.Atoi(division)
		if err == nil {
			division = fmt.Sprintf("%.2d", division_num)
		}
		district := strings.Trim(strings.Split(c.PostForm("DistrictNumber"), "-")[0], " ")
		district_num, err := strconv.Atoi(district)
		if err == nil {
			district = fmt.Sprintf("%.2d", district_num)
		}
		tableName := c.PostForm("TableName")
		header, _ := templates.RenderHeader(c)

		fmt.Printf("division : %s , district : %s, tableName : %s\n", division, district, tableName)

		srv.zilaOkWithData(c, header, footer, division, district, tableName)
	})
}

func (srv *Server) zilaOkWithData(c *gin.Context, header, footer, division, district, tableNumber string) {
	report := ""
	if !(tableNumber != "") && !(district != "") && !(division != "") {
		report = ""
	} else {
		reportName := fmt.Sprintf("Zila-Series-Report-%s-%s-%s.pdf", division, district, tableNumber)
		report = fmt.Sprintf(`
		<div class="col-md-12 col-sm-12 col-xs-12">
		<embed src="/production/getReport?no=%s" width=100%% height=1000
			type='application/pdf' />
		</div>`, reportName)
	}
	data := gin.H{
		"Header":      template.HTML(header),
		"Footer":      template.HTML(footer),
		"Division":    division,
		"District":    district,
		"TableNumber": tableNumber,
		"PDF":         template.HTML(report),
	}

	c.HTML(http.StatusOK, "zila.html", data)
}

func (srv *Server) zilaWithError(c *gin.Context, header, footer, alertMsg, division, district string) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)

	data := gin.H{
		"Header":      template.HTML(header),
		"Footer":      template.HTML(footer),
		"AlertQrCode": template.HTML(alert),
		"Division":    division,
		"District":    district,
	}
	c.HTML(http.StatusOK, "zila.html", data)
}
