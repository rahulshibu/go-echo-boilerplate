package main

import (
	"github.com/rahulshibu/go-echo-boilerplate/internal/configs"
	"github.com/rahulshibu/go-echo-boilerplate/internal/controllers"
	"github.com/rahulshibu/go-echo-boilerplate/internal/database"
	midw "github.com/rahulshibu/go-echo-boilerplate/internal/middleware"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

func setRouters(e *echo.Echo) {
	//set routers
	controllers.RegisterRoutes(e)
	controllers.RegisterErrorHanders(e)
}

func main() {

	//initialize application with toml file
	if err := configs.Init(); err != nil {
		log.Error(err)
	}

	//initialize echo
	e := echo.New()

	//tweaks for server
	e.HideBanner = true
	e.Debug = true

	// Middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	//Add Auth Layer
	e.Use(midw.Authorize)

	//start the databse
	dbconn := database.ConnectSQL()
	defer dbconn.Close()

	//set migration
	//dbconn.AutoMigrate(&models.User{}, &models.Status{})

	//set routers
	setRouters(e)

	// Start server
	e.Logger.Fatal(e.Start(configs.AppConfig.Server.Host + ":" + configs.AppConfig.Server.Port))

}
