package main

import (
	"do_aug25/db"
	"do_aug25/middleware"
	"do_aug25/models"
	"do_aug25/router"
	"flag"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	conn, _ := db.SQLiteConnect()

	var port int64
	printHelp := false
	flag.Int64Var(&port, "port", 8080, "Specify port to bind to.")
	flag.BoolVar(&printHelp, "h", false, "Print help message.")
	flag.Parse()

	if printHelp {
		flag.PrintDefaults()
		return
	}

	err := conn.AutoMigrate(&models.Cat{}, &models.Mission{},
		&models.Target{}, &models.TargetNote{})
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	middleware.InitLogging()
	r.Use(middleware.CustomLoggingMiddleware())
	router.SetupRoutes(r)
	models.InitValidators()
	r.Run(fmt.Sprintf(":%d", port))
}
