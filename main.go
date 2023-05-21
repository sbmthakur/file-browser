package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type FileDetails struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size string `json:"size"`
}

type Resp struct {
	FileDetails
	Contents []FileDetails `json:"contents"`
}

func get_files(file_path string) (string, error) {

	dir := "/home/shubham/workspace/file-browser/storage/" + file_path
	fmt.Println(dir)

	// Open the directory.
	f, err := os.Open(dir)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Read the directory entries.
	entries, err := f.Readdir(-1)
	if err != nil {
		return "", err
	}

	files := make([]FileDetails, 0)

	for _, entry := range entries {

		t, s := "", ""

		if entry.IsDir() {
			t = "dir"
		} else {
			t = "file"
			s = convert_unit(entry.Size())
		}

		new_file := FileDetails{Name: entry.Name(), Type: t, Size: s}
		files = append(files, new_file)
	}

	r := Resp{FileDetails: FileDetails{Name: "/", Type: "dir", Size: ""}, Contents: files}

	data, _ := json.Marshal(r)
	return string(data), nil
}

func convert_unit(size int64) string {
	s := [3]string{"Bytes", "KB", "MB"}
	i := 0
	ss := float64(size)

	for ss > 1024 {
		ss = ss / 1024
		i++
	}

	ss = math.Ceil(ss)
	return strconv.FormatFloat(ss, 'f', -1, 64) + " " + s[i]
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Routes
	e.Static("/static", "web-app/build/static")
	//e.Use(middleware.Static("web-app/build"))

	e.GET("/:dirPath", func(c echo.Context) error {
		return c.File("web-app/build/index.html")
	})

	e.GET("/files", file_handler)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func file_handler(c echo.Context) error {
	file_path := c.QueryParam("path")
	response, err := get_files(file_path)

	if err == nil {
		return c.String(http.StatusOK, response)
	}

	if strings.Contains(err.Error(), "no such file") {
		return c.String(http.StatusNotFound, "Not Found")
	}

	return c.String(http.StatusInternalServerError, "Something failed")
}

func redirect_root(c echo.Context) error {
	return c.Redirect(http.StatusMovedPermanently, "/folders")
}
