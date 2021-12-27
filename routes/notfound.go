package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (srv *Server) notFound() {
	srv.router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusOK, "page_404.html", gin.H{})
	})
}
