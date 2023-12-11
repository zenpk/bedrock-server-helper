package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/zenpk/bedrock-server-helper/dal"
	"log"
)

const (
	path   = "mc.db"
	jwkEnd = "https://example.com/public-key"
)

type jwtCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	db := &dal.Db{}
	if err := db.Connect(path); err != nil {
		panic(err)
	}
	defer db.Db.Close()

	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	jwtConfig := echojwt.Config{
		KeyFunc: getKey,
	}
	e.Use(echojwt.WithConfig(jwtConfig))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			user := c.Get("user").(*jwt.Token) // User is injected by the JWT middleware
			claims := user.Claims.(jwtCustomClaims)
			log.Printf("username: %v\n", claims.Username)
			return nil
		},
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogURI:       true,
		LogURIPath:   true,
		LogRoutePath: true,
		LogStatus:    true,
	}))

	handlers := Handlers{Db: db}
	e.GET("/backups/list", handlers.backupsList)
	e.GET("/versions/list", handlers.backupsList)

	e.Logger.Fatal(e.Start(":1323"))
}

func getKey(token *jwt.Token) (interface{}, error) {
	keySet, err := jwk.Fetch(context.Background(), jwkEnd)
	if err != nil {
		return nil, err
	}
	keyID, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("expecting JWT header to have a key ID in the kid field")
	}
	key, found := keySet.LookupKeyID(keyID)
	if !found {
		return nil, fmt.Errorf("unable to find key %q", keyID)
	}
	var pubKey interface{}
	if err := key.Raw(&pubKey); err != nil {
		return nil, fmt.Errorf("unable to get the public key. Error: %s", err.Error())
	}
	return pubKey, nil
}
