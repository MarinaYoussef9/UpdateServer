package core

import (
	"encoding/json"
	"strconv"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
)

// GetReleases return all releases recorded at the database orded by id
func GetReleases() []database.Release {
	var releases []database.Release
	db.DB.Order("id").Find(&releases)
	return releases
}

// GetRelease
func GetRelease(id string) database.Release {
	var release database.Release
	db.DB.Find(&release, "id = ?", id)
	return release
}

// NewRelease
func NewRelease(release *database.Release) {
	db.DB.Create(&release)
}

// EditRelease
func EditRelease(release *database.Release, id string, bindingSignature string, bindingRelease string) {
	id_buf, _ := strconv.Atoi(id)
	release.ID = uint(id_buf)
	db.DB.First(&release)
	// This looks super ugly, I must improve it.
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
	json.Unmarshal([]byte(bindingRelease), &buf)
	release.Channel = buf.Channel
	release.OS = buf.OS
	release.OsVer = buf.OsVer
	release.OsArch = buf.OsArch
	release.ProductVersion = buf.ProductVersion
	release.URL = buf.URL
	release.Title = buf.Title
	release.Description = buf.Description
	release.Product = buf.Product
	release.Signature = bindingSignature

	db.DB.Save(&release)
}

func DeleteRelease(release_id string) {
	var release database.Release
	db.DB.Where("id = ?", release_id).Delete(&release)
}
