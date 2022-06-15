package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Main website",
	})
}
