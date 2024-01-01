package main

import (
	"github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/database"
	"github.com/Faqihrofiqi/task-5-pbi-btpns-faqih-rofiqi/router"
)

func main() {
	database.ConnectDatabase()
	r := router.InitRouter()

	// Mulai server Gin
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
