package main

import (
	"github.com/cektrendstudio/cektrend-engine-go/config"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetReportCaller(true)

	app := config.Init()
	config.Catch(app.Start())
}
