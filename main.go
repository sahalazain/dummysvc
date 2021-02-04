package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

//KeyReq key req
type KeyReq struct {
	Key string `json:"key,omitempty" mapstructure:"key"`
}

func main() {
	log.SetLevel(log.DebugLevel)
	e := echo.New()

	e.GET("/ok/*", func(c echo.Context) error {
		log.Debug("[get]", c.Request().URL.Path)
		return c.String(http.StatusOK, "OK")
	})

	e.GET("/header/*", func(c echo.Context) error {
		log.Debug("[get]", c.Request().URL.Path)
		return c.JSON(http.StatusOK, c.Request().Header)
	})

	e.GET("/query/*", func(c echo.Context) error {
		log.Debug("[get]", c.Request().URL.Path)
		return c.JSON(http.StatusOK, c.QueryParams())
	})

	e.POST("/echo/*", func(c echo.Context) error {
		log.Debug("[post]", c.Request().URL.Path)
		obj := make(map[string]interface{})
		if err := c.Bind(&obj); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, obj)
	})

	e.GET("/menu/bundle", serveBundle)
	e.GET("/menu", func(c echo.Context) error {
		hash := hash64(menu)
		etag := c.Request().Header.Get("If-None-Match")
		if etag == hash {
			return c.NoContent(http.StatusNotModified)
		}

		c.Response().Header().Set("Etag", hash)
		return c.JSON(http.StatusOK, menu)
	})

	e.GET("/", func(c echo.Context) error {
		return c.File("index.html")
	})

	e.Logger.Fatal(e.Start(":8000"))
}

var menu = map[string][]string{
	"default": {
		"home",
		"user",
		"account",
		"dashboard",
	},
	"admin": {
		"config",
		"admin",
		"policy",
	},
}

//BundleManifest bundle manifest object
type BundleManifest struct {
	Roots []string `json:"roots,omitempty" mapstructure:"roots"`
}

func serveBundle(c echo.Context) error {
	dir := os.TempDir()

	path := dir + "/bundle.tar.gz"

	if !FileExists(path) {
		if err := menuToBundle(menu); err != nil {
			return err
		}
	}

	hash, err := HashFile(path)
	if err != nil {
		log.Error("Error hashing file ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": "Error calculating bundle hash " + err.Error()})
	}

	etag := c.Request().Header.Get("If-None-Match")

	if etag == hash {
		return c.NoContent(http.StatusNotModified)
	}

	c.Response().Header().Set(echo.HeaderContentType, "application/gzip")
	c.Response().Header().Set("Etag", hash)

	return c.File(path)
}

func menuToBundle(menu map[string][]string) error {
	dir := os.TempDir()

	bundleRoot := BundleManifest{
		Roots: []string{"menu"},
	}

	if err := ObjectToFile(dir+"/.manifest", bundleRoot); err != nil {
		return err
	}

	if err := ObjectToFile(dir+"/data.json", menu); err != nil {
		return err
	}

	if err := CompressFile(dir+"/bundle.tar.gz", []string{dir + "/data.json", dir + "/.manifest"}, []string{"menu/data.json", ".manifest"}); err != nil {
		return err
	}
	os.Remove(dir + "/.manifest")
	return os.Remove(dir + "/data.json")
}
