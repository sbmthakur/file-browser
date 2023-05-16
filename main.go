package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type FileDetails struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Size int    `json:"size"`
}

type Resp struct {
	FileDetails
	Contents []FileDetails `json:"contents"`
}

func get_files(file_path string) string {

	dir := "/home/shubham/workspace/file-browser/storage/" + file_path
	fmt.Println(dir)

	// Open the directory.
	f, err := os.Open(dir)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer f.Close()

	// Read the directory entries.
	entries, err := f.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	files := make([]FileDetails, 0)

	for _, entry := range entries {

		t, s := "", 0

		if entry.IsDir() {
			t = "dir"
		} else {
			t = "file"
			s = int(entry.Size())
		}

		new_file := FileDetails{Name: entry.Name(), Type: t, Size: s}
		files = append(files, new_file)
	}

	r := Resp{FileDetails: FileDetails{Name: "/", Type: "dir", Size: 0}, Contents: files}

	data, _ := json.Marshal(r)
	return string(data)
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
	e.Static("/static", "assets")
	e.GET("/files", file_handler)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func file_handler(c echo.Context) error {
	file_path := c.QueryParam("path")
	fmt.Println(file_path)
	return c.String(http.StatusOK, get_files(file_path))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, get_files(""))
}
