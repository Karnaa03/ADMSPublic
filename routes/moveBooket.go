package routes

import (
	"fmt"
	"html/template"
	"net/http"

	ginoidc "git.solutions.im/Solutions.IM/ginOidc"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

type MoveBookletRequest struct {
	BookletNumber  string `form:"BookletNumber"`
	SourceBox      string `form:"SourceBox"`
	DestinationBox string `form:"DestinationBox"`
}

func (srv *Server) moveBooklet(footer string) {
	srv.router.GET("/warehouse/moveBooklet.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		srv.moveBookletOkWithData(c, header, footer, MoveBookletRequest{})
	})

	srv.router.POST("/warehouse/moveBooklet.html", func(c *gin.Context) {
		header, _ := templates.RenderHeader(c)
		id, err := ginoidc.GetIdentity(c)
		if err != nil {
			log.Error(err)
			srv.moveBookletWithError(c, header, footer, "who are you ?", MoveBookletRequest{})
			return
		}

		var mr MoveBookletRequest
		err = c.ShouldBind(&mr)
		if err != nil {
			log.Error(err)
			srv.moveBookletWithError(c, header, footer, "unable to find this booklet for moving", mr)
		}
		srv.move(c, header, footer, id, mr)
	})
}

func (srv *Server) move(c *gin.Context, header, footer string, id ginoidc.Identity, r MoveBookletRequest) {
	b, err := srv.Db.GetBooklet(r.BookletNumber)
	if err != nil {
		log.Error(err)
		srv.moveBookletWithError(c, header, footer, "unable to find booklet with this number", r)
	}
	err = b.MoveInBox(id, r.SourceBox, r.DestinationBox, &srv.Db)
	if err != nil {
		log.Error(err)
		srv.moveBookletWithError(c, header, footer, fmt.Sprintf("unable to move booklet : %s", err), r)
	}
	srv.moveBookletOkWithData(c, header, footer, MoveBookletRequest{})
}

func (srv *Server) moveBookletOkWithData(c *gin.Context, header, footer string, r MoveBookletRequest) {
	name := ginoidc.GetValue(c, "name")

	var content = gin.H{
		"Name":          name,
		"Header":        template.HTML(header),
		"Footer":        template.HTML(footer),
		"BookletNumber": r.BookletNumber,
	}

	c.HTML(http.StatusOK, "moveBooklet.html", content)
}

func (srv *Server) moveBookletWithError(c *gin.Context, header, footer, alertMsg string, r MoveBookletRequest) {
	alert, err := templates.RenderAlert(alertMsg)
	if err != nil {
		log.Error(err)
	}
	log.Error(alertMsg, err)
	name := ginoidc.GetValue(c, "name")

	// tableData := srv.searchInArchivesWithQuery(sr)

	c.HTML(http.StatusOK, "moveBooklet.html", gin.H{
		"Name":           name,
		"Header":         template.HTML(header),
		"Footer":         template.HTML(footer),
		"Alert":          template.HTML(alert),
		"BookletNumber":  r.BookletNumber,
		"SourceBox":      r.SourceBox,
		"DestinationBox": r.DestinationBox,
	})
}
