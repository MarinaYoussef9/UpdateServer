package core

import (
	"log"
	"time"

	"code.videolan.org/GSoC2017/Marco/UpdateServer/database"
)

// GetRequests return all requests recorded at the database associated with a channel and product
func GetRequests(channel, product string) []database.UpdateRequest {
	var requests []database.UpdateRequest
	db.DB.Where("product = ? AND channel = ?", product, channel).Order("created_at desc").Find(&requests)
	ProcessCreatedSince(&requests)
	return requests
}

func NewRequest(request database.UpdateRequest) {
	db.DB.Create(&request)
}

// ProcessCreatedSince initiate the CreateSince Section at update_requests
func ProcessCreatedSince(requests *[]database.UpdateRequest) {
	TimeNow := time.Now().UTC()
	for i := 0; i < len(*requests); i++ {
		Duration := TimeNow.Sub((*requests)[i].CreatedAt.UTC())
		(*requests)[i].CreatedSince.Month = int(Duration.Hours() / (24 * 30))
		(*requests)[i].CreatedSince.Day = TimeNow.Day() - (*requests)[i].CreatedAt.UTC().Day()
		(*requests)[i].CreatedSince.Hour = TimeNow.Hour() - (*requests)[i].CreatedAt.UTC().Hour()
		(*requests)[i].CreatedSince.Minute = TimeNow.Minute() - (*requests)[i].CreatedAt.UTC().Minute()
		(*requests)[i].CreatedSince.Second = TimeNow.Second() - (*requests)[i].CreatedAt.UTC().Second()
	}
}

// GetSignature function return the signature of a given release.
func GetSignature(release_id string) string {
	var release database.Release
	db.DB.First(&release, "id = ?", release_id)
	return release.Signature
}

// ReleaseMap function map the incoming update_request with the most suitable update release
func ReleaseMap(request database.UpdateRequest) (database.Release, bool) {
	var release database.Release
	// First , find and count the available releases match the request specs
	var releasescount int
	var releases []database.Release

	db.DB.Where("product = ? AND channel = ? AND os = ? AND os_arch = ? AND os_ver >= ?",
		request.Product, request.Channel, request.OS, request.OsArch, request.OsVer).Find(&releases).Count(&releasescount)
	log.Println(releasescount, releases)

	if releasescount == 0 {
		return release, false
	} else {
		return releases[0], true
	}

}
