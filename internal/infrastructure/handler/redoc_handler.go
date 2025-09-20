package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RedocHandler struct{}

func NewRedocHandler() *RedocHandler {
	return &RedocHandler{}
}

func (h *RedocHandler) Register(router *gin.RouterGroup) {
	router.GET("", h.ServeRedoc)
}

func (h *RedocHandler) ServeRedoc(c *gin.Context) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>API Documentation - Redoc</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css?family=Montserrat:300,400,700|Roboto:300,400,700" rel="stylesheet">
    <style>
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
    <redoc spec-url='/docs/doc.json'></redoc>
    <script src="https://cdn.jsdelivr.net/npm/redoc/bundles/redoc.standalone.js"> </script>
</body>
</html>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}
