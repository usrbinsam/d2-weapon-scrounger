package db

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	DiscordId            *string `gorm:"index"`
	DiscordGuildId       *string
	BungieAuthCode       *string
	BungieAccessToken    *string
	BungieRefreshToken   *string
	BungieMembershipIdId *string
	State                *string `gorm:"uniqueIndex"`
}

type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	MembershipId string `json:"membership_id"`
}

// Request an access token & refresh token from the Bungie OAuth2 API
func (user *User) RequestBungieAccessToken() error {
	resp, err := http.PostForm(
		"https://www.bungie.net/platform/app/oauth/token/",
		url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {*user.BungieAuthCode},
			"client_id":     {os.Getenv("CLIENT_ID")},
			"client_secret": {os.Getenv("CLIENT_SECRET")},
		},
	)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("bungie api responded with status %d: %s", resp.StatusCode, body)
		return errors.New("non-200 status code from bungie api")
	}

	var creds Credentials
	err = json.NewDecoder(resp.Body).Decode(&creds)

	if err != nil {
		log.Printf("failed to decode response from bungie api: %s", err.Error())
		return err
	}

	user.BungieAccessToken = &creds.AccessToken
	user.BungieRefreshToken = &creds.RefreshToken
	user.BungieMembershipIdId = &creds.MembershipId

	return nil
}

// Open is a simplified wrapper around gorm.Open for sqlite.
func Open(filename string, config *gorm.Config) (*gorm.DB, error) {
	log.Printf("db: opening %q", filename)
	if config == nil {
		config = &gorm.Config{}
	}

	db, err := gorm.Open(sqlite.Open(filename), config)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
