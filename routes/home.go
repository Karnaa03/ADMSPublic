package routes

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
)

func (srv *Server) homePage(footer string) gin.IRoutes {
	return srv.router.GET("/", func(c *gin.Context) {
		header, err := templates.RenderHeader(c)
		if err != nil {
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"Header": template.HTML(header),
			"Footer": template.HTML(footer),
		})
	})
}
