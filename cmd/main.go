package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/fs"
	"net/http"
	"os"
)

type logsRequest struct {
	Logs []log `json:"logs"`
}

type singleLogRequest struct {
	Log log `json:"log"`
}

type log struct {
	Action   string   `json:"action"`
	NodeData nodeData `json:"node"`
	Time     int      `json:"time"`
}

type nodeData struct {
	Node      int        `json:"node"`
	Position  []float64  `json:"position"`
	Edges     [][]string `json:"edges"`
	Toplabel  []string   `json:"toplabel"`
	Botlabel  []string   `json:"botlabel"`
	Valuation string
}

type handler interface {
	persistLogs() echo.HandlerFunc
	persistSingleLog() echo.HandlerFunc
	getFile() *os.File
}

type Handler struct {
	File *os.File
}

func initHandler() handler {
	file, err := os.OpenFile("logs.jsonl", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	return Handler{
		File: file,
	}
}

func main() {
	h := initHandler()
	defer func(file *os.File) {
		fmt.Print("Closing file")
		if err := file.Close(); err != nil {
			panic(err)
		}
	}(h.getFile())

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("logs", h.persistLogs())
	e.POST("log", h.persistSingleLog())

	e.Logger.Fatal(e.Start(":8000"))
}

func (h Handler) persistSingleLog() echo.HandlerFunc {
	return func(c echo.Context) error {
		body := new(singleLogRequest)

		if err := c.Bind(body); err != nil {
			c.Logger().Error(err)
		}

		if err := h.AppendLog(body.Log); err != nil {
			c.Logger().Error(err)
		}

		return c.String(http.StatusCreated, "Persisted log successfully!")
	}
}

func (h Handler) persistLogs() echo.HandlerFunc {
	return func(c echo.Context) error {
		body := new(logsRequest)

		err := c.Bind(body)
		if err != nil {
			c.Logger().Error(err)
		}

		bytes, err := json.Marshal(body)
		if err != nil {
			c.Logger().Error(err)
		}

		if err := os.WriteFile("logs.jsonl", bytes, fs.ModePerm); err != nil {
			c.Logger().Error(err)
		}

		return c.String(http.StatusCreated, "Persisted logs successfully!")
	}
}

func (h Handler) getFile() *os.File {
	return h.File
}

func (h Handler) AppendLog(l log) error {
	bytes, err := json.Marshal(l)
	if err != nil {
		return err
	}

	if _, err := h.File.Write(bytes); err != nil {
		return err
	}

	if _, err := h.File.WriteString("\n"); err != nil {
		return err
	}

	return nil
}
