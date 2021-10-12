package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/usrbinsam/d2-weapon-scrounger/db"
	"gorm.io/gorm"
)

var Session *gorm.DB

func generateToken() string {
	b := make([]byte, 32)
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
		io.WriteString(rw, "mssing 'code' or 'state' parameter. try again")
		return
	}

	var user db.User
	Session.Where("state = ?", state).Take(&user)
	user.BungieAuthCode = &code

	err := user.RequestBungieAccessToken()
	if err != nil {
		io.WriteString(rw, "failed to request access token, try again")
		return
	}

	io.WriteString(rw, "successfully authorized. you may close this page.")
}

func main() {
	if os.Getenv("APPMODE") == "dev" {
		godotenv.Load(".development.env")
	} else {
		godotenv.Load(".env")
	}

}
