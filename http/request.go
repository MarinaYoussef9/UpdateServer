package http

import (
	"net/http"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"github.com/gin-gonic/gin"
)

// GetRequests is http handler to represent all the requests available in the UpdateServer
func GetRequests(c *gin.Context) {
	requests := core.GetRequests(c.Param("channel"), c.Param("product"))
	c.HTML(http.StatusOK, "requests.html", gin.H{
		"requests": requests,
	})
}

// GetSignature
func GetSignature(c *gin.Context) {
	releaseSignature := core.GetSignature(c.Query("id"))
	c.String(http.StatusOK, releaseSignature)
}