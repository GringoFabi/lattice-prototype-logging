package main

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	getLogFile() *os.File
	getLocalStorageLogFile() *os.File
}

type Handler struct {
	LogFile             *os.File
	LocalStorageLogFile *os.File
}

func initHandler() handler {
	logFile, err := os.OpenFile("logs.jsonl", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	localStorageLogFile, err := os.OpenFile("localStorageLogs.jsonl", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	return Handler{
		LogFile:             logFile,
		LocalStorageLogFile: localStorageLogFile,
	}
}

func main() {
	h := initHandler()
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}(h.getLogFile())

	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}(h.getLocalStorageLogFile())

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

		if err := h.AppendLogsToLocalStorageLogFile(body.Logs); err != nil {
			c.Logger().Error(err)
		}

		return c.String(http.StatusCreated, "Persisted logs successfully!")
	}
}

func (h Handler) getLogFile() *os.File {
	return h.LogFile
}

func (h Handler) getLocalStorageLogFile() *os.File {
	return h.LocalStorageLogFile
}

func (h Handler) AppendLog(l log) error {
	bytes, err := json.Marshal(l)
	if err != nil {
		return err
	}

	if _, err := h.LogFile.Write(bytes); err != nil {
		return err
	}

	if _, err := h.LogFile.WriteString("\n"); err != nil {
		return err
	}

	return nil
}

func (h Handler) AppendLogsToLocalStorageLogFile(logs []log) error {
	bytes, err := json.Marshal(logs)
	if err != nil {
		return err
	}

	if _, err := h.LocalStorageLogFile.Write(bytes); err != nil {
		return err
	}

	if _, err := h.LocalStorageLogFile.WriteString("\n"); err != nil {
		return err
	}

	return nil
}
