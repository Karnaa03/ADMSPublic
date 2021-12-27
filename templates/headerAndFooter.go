package templates

import (
	"bytes"
	"html/template"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	gin_oidc "git.solutions.im/Solutions.IM/ginOidc"
)

type HeaderUser struct {
	Name string
	User gin_oidc.Identity
}

func RenderHeader(c *gin.Context) (content string, err error) {
	// user, err := gin_oidc.GetIdentity(c)
	// if err != nil {
	// 	return
	// }
	header, err := ioutil.ReadFile("templates/html/header.html")
	if err != nil {
		return
	}
	t, err := template.New("header").Parse(string(header))
	if err != nil {
		return
	}
	var h bytes.Buffer

	err = t.Execute(&h, HeaderUser{
		// Name: strings.Title(user.FullName),
		// User: user,
	})
	if err != nil {
		return
	}
	return h.String(), nil
}

func RenderFooter(version, hostname string) (content string, err error) {
	footer, err := ioutil.ReadFile("templates/html/footer.html")
	if err != nil {
		return
	}
	t, err := template.New("header").Parse(string(footer))
	if err != nil {
		return
	}
	var h bytes.Buffer

	err = t.Execute(&h, struct {
		Version  string
		Hostname string
	}{
		Version:  version,
		Hostname: hostname,
	})
	if err != nil {
		return
	}
	return h.String(), nil
}

func RenderAlert(alertMessage string) (content string, err error) {
	alert, err := ioutil.ReadFile("templates/html/alert.html")
	if err != nil {
		return
	}
	t, err := template.New("alert").Parse(string(alert))
	if err != nil {
		return
	}
	var h bytes.Buffer

	err = t.Execute(&h, struct{ AlertMessage string }{AlertMessage: alertMessage})
	if err != nil {
		return
	}
	return h.String(), nil
}

func RenderWarning(alertMessage string) (content string, err error) {
	alert, err := ioutil.ReadFile("templates/html/warning.html")
	if err != nil {
		return
	}
	t, err := template.New("alert").Parse(string(alert))
	if err != nil {
		return
	}
	var h bytes.Buffer

	err = t.Execute(&h, struct{ AlertMessage string }{AlertMessage: alertMessage})
	if err != nil {
		return
	}
	return h.String(), nil
}
