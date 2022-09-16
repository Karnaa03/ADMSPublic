package routes

import (
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.solutions.im/XeroxAgriCensus/ADMSPublic/templates"
)

func (srv *Server) contact_us(footer string) gin.IRoutes {
	return srv.router.GET("/production/contact_us.html", func(c *gin.Context) {
		header, err := templates.RenderHeader(c)
		if err != nil {
			return
		}
		c.HTML(http.StatusOK, "contact_us.html", gin.H{
			"Header": template.HTML(header),
			"Footer": template.HTML(footer),
		})
	})
}