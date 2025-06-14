package lib

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func SetLogLevelFromEnv() {
	logLevel := os.Getenv("LOG")
	if logLevel == "" {
		logLevel = "info"
	}
	switch logLevel {
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	default:
		LogFatal("Invalid LOG level: %s", logLevel)
	}
}

func LogFatal(msg string, args ...any) {
	slog.Error(msg, args...)
	panic(msg)
}

func BootstrapApp() Context {
	godotenv.Load()

	note := flag.String("note", "", "Note ID to edit")
	teamPath := flag.String("team", "", "Path to team")
	flag.Parse()

	if *note == "" {
		LogFatal("Note ID is required")
	}

	SetLogLevelFromEnv()

	ctx := Context{
		Note:       *note,
		TeamPath:   *teamPath,
		BaseUrl:    os.Getenv("HACKMD_BASE_URL"),
		ApiToken:   os.Getenv("HACKMD_API_TOKEN"),
		Editor:     os.Getenv("EDITOR"),
		EditorArgs: os.Getenv("EDITOR_ARGS"),
		Client:     &http.Client{},
	}

	if ctx.BaseUrl == "" {
		ctx.BaseUrl = "https://api.hackmd.io/v1/"
	}
	if !strings.HasSuffix(ctx.BaseUrl, "/") {
		ctx.BaseUrl += "/"
	}
	if ctx.ApiToken == "" {
		LogFatal("Env var HACKMD_API_TOKEN is required")
	}
	if ctx.Editor == "" {
		LogFatal("Env var EDITOR is required")
	}
	slog.Debug("Context", "context", ctx)
	return ctx
}
