package http

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// RouterInit function initiate the gin engine main handler
func RouterInit() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/", Auth)
	admin := auth.Group("/admin/dashboard")
	{
		admin.GET("/channels", GetChannels)
		admin.GET("/releases", GetReleases)
		admin.GET("/release/:id", GetRelease)
		admin.GET("/channel/:name", GetChannel)
		admin.GET("/channels/add", AddChannel)
		admin.GET("/newrelease", AddRelease)
		admin.GET("/addsignature/:reference", AddSignature)
		admin.GET("/del_release/:id", DelRelease)
		admin.GET("/add_rule/:id", AddRule)
		admin.GET("/delete_rule/:rule/:id", DeleteRule)

		admin.POST("/new_rule/:rule/:id", NewRule)
		admin.POST("/edit_release/:id", EditRelease)
		admin.POST("/verifysignature/:reference", VerifySignature)
		admin.POST("/new_release", NewRelease)
		admin.POST("/new_channel", NewChannel)
	}

	pub := router.Group("/u/:product/:channel")
	{
		pub.GET("/requests", GetRequests)
		pub.GET("/signature", GetSignature)
		pub.GET("/update", Update)
	}
	router.GET("act", GetAct)
	router.POST("act", GetAct)

	router.LoadHTMLGlob("html/*.html")
	return router
}

// Run function initaite the UpdateServer
func Run(addr string) {
	RouterInit().Run(":" + addr)
}

func Auth(c *gin.Context) {
	if checkAuth(c) {
		c.Next()
	} else {
		c.Writer.Header().Set("WWW-Authenticate", "Basic realm=UpdateServer")
		c.AbortWithStatus(http.StatusUnauthorized)

	}

}
func checkAuth(c *gin.Context) bool {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)
	if len(auth) != 2 {
		return false
	}

	base, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		return false
	}

	UserData := strings.SplitN(string(base), ":", 2)
	if len(UserData) != 2 {
		return false
	}
	username := UserData[0]
	password := UserData[1]
	return username == "admin" && password == "admin"
}
