package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/usrbinsam/d2-weapon-scrounger/db"
	"gorm.io/gorm"
)

var Session *gorm.DB

func generateToken() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("rand.Read() failed: %s", err.Error())
	}
	return hex.EncodeToString(b)
}

func formAuthUrl(state string) string {
	redirectUrl := os.Getenv("BASE_URL") + "/auth"

	return fmt.Sprintf(
		"https://www.bungie.net/en/oauth/authorize?"+
			"client_id=%s"+
			"&response_type=code"+
			"&state=%s"+
			"&redirect_url=%s",
		os.Getenv("CLIENT_ID"),
		state,
		redirectUrl,
	)
}

func authHandler(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	code := q.Get("code")
	state := q.Get("state")

	if code == "" || state == "" {
		log.Printf("auth: missing 'code' or 'state' parameter: %+v", q)
		io.WriteString(rw, "mssing 'code' or 'state' parameter. try again")
		return
	}

	var user db.User
	err := Session.Where("state = ?", state).Take(&user).Error

	if err != nil {
		log.Printf("invalid state %q: %s", state, err.Error())
		io.WriteString(rw, "state not found")
		return
	}

	user.BungieAuthCode = &code

	err = user.RequestBungieAccessToken()
	if err != nil {
		io.WriteString(rw, "failed to request access token, try again")
		return
	}

	Session.Save(&user)

	io.WriteString(rw, "successfully authorized. you may close this page.")
}

func authUrlHandler(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	discordId := q.Get("discord_id")

	// TODO: verify the discordId provided here with a cryptographic signature
	// prevents idiots from polluting the database with nonsense so only the
	// bot can issue valid /start requests.

	if discordId == "" {
		log.Printf("auth URL requested but no discord_id was provided")
		io.WriteString(rw, "'discord_id' not provided")
		return
	}

	var user db.User
	state := generateToken()
	Session.Where(db.User{DiscordId: &discordId}).Attrs(db.User{State: &state}).FirstOrCreate(&user)

	rw.Header().Add("Location", formAuthUrl(*user.State))
	rw.WriteHeader(http.StatusFound)
}

func StartAuthServer(addr string) {
	session, err := db.Open("app.db", nil)

	if err != nil {
		log.Fatalf("failed to open app.db: %s", err.Error())
	}

	// assign global Session for this package
	Session = session

	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/start", authUrlHandler)

	log.Printf("auth server starting on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
