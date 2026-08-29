package docs

import (
	"net/http"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gin-gonic/gin"
)

func Docs(c *gin.Context) {
	html, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: "./openapi.yaml",
		CustomOptions: scalar.CustomOptions{
			PageTitle: "My API",
		},
	})

	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Data(
		http.StatusOK,
		"text/html; charset=utf-8",
		[]byte(html),
	)
}
