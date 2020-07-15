package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Reply Log
type Reply struct {
	Response  string    `json:"response"`
	Timestamp time.Time `json:"timestamp"`
	Random    int       `json:"random"`
}

func main() {
	e := echo.New()
	godotenv.Load()
	env := os.Getenv("GOLANG_ENV")
	if env == "production" {
		godotenv.Load(".env.production")
	} else {
		godotenv.Load(".env.develop")
	}
	Port := os.Getenv("APP_PORT")

	e.Use(
		middleware.Recover(),
		middleware.Logger(),
		middleware.RequestID(),
	)

	e.GET("/", func(c echo.Context) error {
		r := &Reply{
			Response:  "Server is running",
			Timestamp: time.Now().UTC(),
			Random:    rand.Intn(1000),
		}
		sr, _ := json.Marshal(r)
		return c.String(http.StatusOK, string(sr))
	})

	e.Pre(APIVersion)

	e.Logger.Fatal(e.Start(":" + Port))
}

// APIVersion Header Based Versioning
func APIVersion(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		headers := req.Header

		apiVersion := headers.Get("version")
		if apiVersion != "" {
			req.URL.Path = fmt.Sprintf("/%s%s", apiVersion, req.URL.Path)
			return next(c)
		}
		return next(c)
	}
}
