package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/zenpk/bedrock-server-helper/dal"
	"github.com/zenpk/bedrock-server-helper/runner"
	"log"
)

const (
	dbPath          = "./mc.db"
	mcPath          = "~/mc"
	baseWorldFolder = "base_world"
	serversFolder   = "servers"
	backupsFolder   = "backups"
)

var (
	jwkEnd = flag.String("jwk", "https://example.com", "JWK public key endpoint")
)

type jwtCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	flag.Parse()
	db := &dal.Db{}
	if err := db.Connect(dbPath); err != nil {
		panic(err)
	}
	defer db.Db.Close()
	runners := &runner.Runner{
		Db:              db,
		McPath:          mcPath,
		BaseWorldFolder: baseWorldFolder,
		ServersFolder:   serversFolder,
		BackupsFolder:   backupsFolder,
	}
	scheduler, err := runners.StartCron()
	if err != nil {
		panic(err)
	}
	defer scheduler.StopJobs()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"*"},
	}))
	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		KeyFunc: getKey,
	}
	e.Use(echojwt.WithConfig(jwtConfig))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			user := c.Get("user").(*jwt.Token) // User is injected by the JWT middleware
			claims := user.Claims.(*jwtCustomClaims)
			log.Printf("username: %v\n", claims.Username)
			log.Println(values.Error)
			return nil
		},
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogURI:       true,
		LogURIPath:   true,
		LogRoutePath: true,
		LogStatus:    true,
		LogError:     true,
	}))

	handlers := &Handlers{
		Db:     db,
		Runner: runners,
	}
	e.GET("/worlds/list", handlers.worldsList)
	e.POST("/worlds/create", handlers.createWorld)
	e.POST("/worlds/upload/:worldId", handlers.uploadWorld)
	e.GET("/backups/list/:worldId", handlers.backupsList)
	e.GET("/servers/list/:worldId", handlers.serversList)
	e.POST("/servers/get", handlers.getServer)
	e.POST("/servers/use", handlers.useServer)
	e.POST("/backups/backup", handlers.backup)
	e.POST("/backups/restore", handlers.restore)
	//e.POST("/worlds/start", handlers.start)
	//e.POST("/worlds/stop", handlers.stop)

	//e.Logger.Fatal(e.StartTLS(":1323", "cert.pem", "key.pem"))
	e.Logger.Fatal(e.Start(":1323"))
}

func getKey(token *jwt.Token) (interface{}, error) {
	keySet, err := jwk.Fetch(context.Background(), *jwkEnd)
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