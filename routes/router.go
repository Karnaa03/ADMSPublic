package routes

import (
	"os"

	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/conf"
	"git.solutions.im/XeroxAgriCensus/AgriTracking/model"
	"git.solutions.im/XeroxAgriCensus/AgriTracking/s3"
	"git.solutions.im/XeroxAgriCensus/AgriTracking/templates"
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

	// srv.AuthParam, err = auth.InitAuth(srv.Config, srv.router)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// srv.CookieStore = cookie.NewStore([]byte("secret"))

	// srv.router.Use(sessions.Sessions("agritracking", srv.CookieStore))
	// // manage authentication
	// srv.router.Use(ginoidc.Init(srv.AuthParam))
	// // manage authorization
	// enforcer, err := casbin.NewEnforcer("authz_model.conf", "authz_policy.csv")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// srv.router.Use(ginoidc.NewAuthorizer(enforcer))

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
	srv.stage1(footer)
	srv.stage2(footer)
	srv.stage3(footer)
	srv.stage4(footer)
	srv.stage5(footer)
	srv.stage6(footer)
	srv.currentStatus(footer)
	srv.registeredStats()
	srv.inBatchStats()
	srv.cutStats()
	srv.preparedStats()
	srv.scannStats()
	srv.archivedStats()
	srv.freeze(footer)
	srv.archiveFrozen(footer)
	srv.registerArchiveBox(footer)
	srv.searchInArchives(footer)
	srv.moveBooklet(footer)
	srv.search(footer)
	srv.tallySheet()
	srv.questionnaire()
	srv.updateGeoCode(footer)
	srv.updateTallySheet(footer)
}

func (srv *Server) registerStatic() {
	srv.router.Static("../vendors", "node_modules/gentelella/vendors")
	srv.router.Static("./images", "node_modules/gentelella/production/images")
	srv.router.Static("./css", "node_modules/gentelella/production/css")
	srv.router.Static("./build", "node_modules/gentelella/build")
	srv.router.LoadHTMLGlob("templates/html/*.html")
}

func (srv *Server) Run() {
	err := srv.router.Run(srv.Config.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("server available on : %s", srv.Config.BaseUrl)
}
