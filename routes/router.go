package routes

import (
	"os"

	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/ADMSPublic/conf"
	"git.solutions.im/XeroxAgriCensus/ADMSPublic/model"
	"git.solutions.im/XeroxAgriCensus/ADMSPublic/s3"
	"git.solutions.im/XeroxAgriCensus/ADMSPublic/templates"
)

type Server struct {
	Config      conf.Config
	Db          model.Db
	S3          s3.S3
	CookieStore cookie.Store
	// AuthParam           ginoidc.InitParams
	router              *gin.Engine
	autoCompleteGeoCode string
}

func (srv *Server) InitRouter() (err error) {
	srv.router = gin.Default()
	srv.registerStatic()
	srv.registerPages()
	return
}

func (srv *Server) registerPages() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("unable to get hostname : %w", err)
	}
	footer, err := templates.RenderFooter(srv.Config.Version, hostname)
	if err != nil {
		log.Fatal(err)
	}
	srv.notFound()
	srv.homePage(footer)
	srv.national_report(footer)
	srv.zila(footer)
	srv.analytical(footer)
	srv.cropping(footer)
	srv.national_report_2008(footer)
	srv.national_report_1996(footer)
	srv.indicator(footer)
	srv.frequency(footer)
	srv.tableGeneration(footer)

}

func (srv *Server) registerStatic() {
	srv.router.Static("../vendors", "node_modules/gentelella/vendors")
	srv.router.Static("./images", "node_modules/gentelella/production/images")
	srv.router.Static("./css", "node_modules/gentelella/production/css")
	srv.router.Static("./build", "node_modules/gentelella/build")
	srv.router.Static("./img", "templates/img")
	srv.router.LoadHTMLGlob("templates/html/*.html")
}

func (srv *Server) Run() {
	err := srv.router.Run(srv.Config.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("server available on : %s", srv.Config.BaseUrl)
}
