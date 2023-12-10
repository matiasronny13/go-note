package main

import (
	"errors"
	"flag"
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/ghostrepo00/go-note/config"
	"github.com/ghostrepo00/go-note/internal/app"
	appconstant "github.com/ghostrepo00/go-note/internal/pkg/app_constant"
)

func getLogFileName(appConfig *config.AppConfig) (result string) {
	currentDate := time.Now()
	return strings.Join([]string{appConfig.Log.FolderPath, "/", currentDate.Format(appconstant.TimestampFormat), "_", appConfig.Log.FileName}, "")
}

func useSlog(appConfig *config.AppConfig) (logFile *os.File, err error) {
	logFile, err = os.OpenFile(getLogFileName(appConfig), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	slogHandler := slog.NewJSONHandler(
		io.MultiWriter(os.Stdout, logFile),
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	logger := slog.New(slogHandler)
	slog.SetDefault(logger)
	return
}

func main() {
	configPath := flag.String("config", "./config.json", "Config file path")
	flag.Parse()

	if config, err := config.NewAppConfig(*configPath); err != nil {
		panic(err)
	} else {
		if fileLog, err := useSlog(config); err == nil {

			defer fileLog.Close()

			slog.Info("App started")
			supabaseConnection := os.Getenv("SUPABASE")
			if supabaseConnection == "" {
				panic(errors.New("SUPABASE env var is not found"))
			}
			supabasePair := strings.Split(supabaseConnection, "|")
			config.SupabaseUrl = supabasePair[0]
			config.SupabaseKey = supabasePair[1]

			webServer := app.NewWebServer(config)
			webServer.Run()

		} else {
			panic(err)
		}
	}
}
