package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/core"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
	"code.videolan.org/GSoC2017/Marco/UpdateServer/utils"
	"github.com/gin-gonic/gin"
)

// GetReleases is http handler to represent all the releases available in the UpdateServer
func GetReleases(c *gin.Context) {
	releases := core.GetReleases()
	c.HTML(http.StatusOK, "releases.html", gin.H{
		"releases": releases,
	})
}

// GetRelease
func GetRelease(c *gin.Context) {
	var rules []database.Rule
	release := core.GetRelease(c.Param("id"))
	channels := core.GetChannels()
	c.HTML(http.StatusOK, "release.html", gin.H{
		"release":  release,
		"rules":    rules,
		"channels": channels,
	})
}

// AddRelease
func AddRelease(c *gin.Context) {
	channels := core.GetChannels()
	c.HTML(http.StatusOK, "newrelease.html", gin.H{
		"channels": channels,
	})
}

// NewRelease
func NewRelease(c *gin.Context) {
	var (
		release database.Release
		buf     struct {
			Channel        string `json:"channel"`
			OS             string `json:"os"`
			OsVer          string `json:"os_ver"`
			OsArch         string `json:"os_arch"`
			ProductVersion string `json:"product_ver"`
			URL            string `json:"url"`
			Title          string `json:"title"`
			Description    string `json:"desc"`
			Product        string `json:"product"`
		}
	)
	c.Bind(&release)
	// FIXME: Isn't there a way to handle that?
	buf.Channel = release.Channel
	buf.OS = release.OS
	buf.OsVer = release.OsVer
	buf.OsArch = release.OsArch
	buf.ProductVersion = release.ProductVersion
	buf.URL = release.URL
	buf.Title = release.Title
	buf.Description = release.Description
	buf.Product = release.Product

	ReleaseJSON, _ := json.Marshal(buf)
	c.SetCookie("release", string(ReleaseJSON), 0, "/", "", false, false)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/addsignature/new")
}

//DelRelease
func DelRelease(c *gin.Context) {
	core.DeleteRelease(c.Param("id"))
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/releases/")
}

// EditRelease
func EditRelease(c *gin.Context) {
	var release database.Release
	c.Bind(&release)

	ReleaseJSON, _ := json.Marshal(release)
	c.SetCookie("release", string(ReleaseJSON), 0, "/", "", false, false)
	c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/addsignature/edit?id="+c.Param("id"))
}

//AddSignature
func AddSignature(c *gin.Context) {
	release, _ := c.Cookie("release")
	var buf struct {
		Channel        string `json:"channel"`
		OS             string `json:"os"`
		OsVer          string `json:"os_ver"`
		OsArch         string `json:"os_arch"`
		ProductVersion string `json:"product_ver"`
		URL            string `json:"url"`
		Title          string `json:"title"`
		Description    string `json:"desc"`
		Product        string `json:"product"`
	}

	json.Unmarshal([]byte(string(release)), &buf)

	ReleaseDir := "static/releases/tmp"
	ReleaseJSON, _ := json.Marshal(buf)
	ioutil.WriteFile(ReleaseDir, ReleaseJSON, 0644)

	ref := c.Param("reference")
	query := "0"
	if ref == "edit" {
		query = c.Query("id")
	}
	c.HTML(http.StatusOK, "newsignature.html", gin.H{
		"status": string(ReleaseJSON),
		"ref":    ref,
		"query":  query,
	})
}

// VerifySignature
func VerifySignature(c *gin.Context) {
	var (
		binding struct {
			Content   string `form:"content" json:"content"`
			Signature string `form:"signature" json:"signature"`
		}
		release database.Release
	)

	c.Bind(&binding)
	json.Unmarshal([]byte(string(binding.Content)), &release)

	SignatureDir := "static/signatures/tmp.asc"
	ioutil.WriteFile(SignatureDir, []byte(binding.Signature), 0644)

	c.SetCookie("release", "", 0, "/", "", false, false)
	if utils.ProcessRelease(release) == true {
		if c.Param("reference") == "new" {
			release.Signature = binding.Signature
			core.NewRelease(&release)

		}
		if c.Param("reference") == "edit" {
			core.EditRelease(&release, c.Query("id"), binding.Signature, binding.Content)
		}

		os.Rename("static/releases/tmp", "static/releases/"+strconv.Itoa(int(release.ID)))
		os.Rename("static/signatures/tmp.asc", "static/signatures/"+strconv.Itoa(int(release.ID))+".asc")
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/release/"+strconv.Itoa(int(release.ID)))

	} else {
		c.Redirect(http.StatusMovedPermanently, "/admin/dashboard/newrelease")
	}
}