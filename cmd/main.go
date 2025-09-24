package main

import (
	"log"

	_ "github.com/Komilov31/event-booker/docs"

	"github.com/Komilov31/event-booker/cmd/app"
)

// @title Event Booker API
// @version 1.0
// @description Event Booker service to create and manage events and bookings.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	if err := app.Run(); err != nil {
		log.Fatal("could not start server: ", err)
	}
}
