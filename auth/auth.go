package auth

import (
	"net/http"
	"net/url"

	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"

	gin_oidc "git.solutions.im/Solutions.IM/ginOidc"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/conf"
)

func InitAuth(config conf.Config, router *gin.Engine) (authParam gin_oidc.InitParams, err error) {
	issuer, err := url.Parse(config.OpenIdURL)
	if err != nil {
		return
	}
	clientURL, err := url.Parse(config.BaseUrl)
	if err != nil {
		return
	}
	postLogout, err := url.Parse(config.BaseUrl)
	if err != nil {
		return
	}

	authParam = gin_oidc.InitParams{
		Router:       router,
		ClientId:     config.OpenIdClientID,
		ClientSecret: config.OpenIdClientSecret,
		Issuer:       *issuer,
		ClientUrl:    *clientURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		ErrorHandler: func(c *gin.Context) {
			// gin_oidc pushes a new error before any "ErrorHandler" invocation
			// message := c.Errors.Last().Error()
			// redirect to ErrorEndpoint with error message
			c.Redirect(http.StatusInternalServerError, config.BaseUrl+"error.html")
			// redirectToErrorPage(c, "http://example2.domain/error", message)
			// when "ErrorHandler" ends "c.Abort()" is invoked - no further handlers will be invoked
		},
		PostLogoutUrl: *postLogout,
		LogoutPath:    config.OpenIdLogoutPath,
	}
	return
}
