package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"lattice-logging/pkg/analyze"
	"net/http"
	"os"
)

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

func getRunMode() string {
	if mode, ok := os.LookupEnv("MODE"); !ok {
		return "server"
	} else {
		return mode
	}
}

func main() {
	mode := getRunMode()
	switch mode {
	case "server":
		runServer()
	case "analyze":
		count1 := analyze.ReadJson("localStorage.json")
		count2 := analyze.ReadJsonLines("logs.jsonl")
		count3 := analyze.ReadJsonLinesArray("localStorageLogs.jsonl")
		zip := analyze.Zip(count1, count2, count3)

		fmt.Printf("\nAction,Count\n")
		for action, count := range zip {
			fmt.Printf("%s,%d\n", action, count)
		}
	}
}

func runServer() {
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
		body := new(analyze.SingleLogRequest)

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
		body := new(analyze.LogsRequest)

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

func (h Handler) AppendLog(l analyze.Log) error {
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

func (h Handler) AppendLogsToLocalStorageLogFile(logs []analyze.Log) error {
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
