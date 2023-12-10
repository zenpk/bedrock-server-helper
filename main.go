package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zenpk/bedrock-server-helper/dal"
)

const (
	path = "mc.db"
)

func main() {
	db := &dal.Db{}
	if err := db.Connect(path); err != nil {
		panic(err)
	}
	defer db.Db.Close()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", _)

	e.Logger.Fatal(e.Start(":1323"))
}
