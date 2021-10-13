package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/usrbinsam/d2-weapon-scrounger/auth"
)

func main() {

	dev := os.Getenv("APPMODE") == "dev"

	authServer := flag.Bool("auth", false, "launch auth webserver server")
	flag.BoolVar(&dev, "dev", dev, "enable dev mode (overrides APPMODE envvar)")
	flag.Parse()

	if dev {
		log.Println("loading develompent environment")
		godotenv.Load(".development.env")
	} else {
		log.Println("loading production environment")
		godotenv.Load(".env")
	}

	if *authServer {
		log.Print("launching auth-server mode")
		auth.StartAuthServer(os.Getenv("AUTH_LISTEN_ADDR"))
	}
}
