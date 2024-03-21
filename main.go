package main

import (
	"UntisChangeNotifier/utils"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

// TODO: Modular config loading
// TODO: variables to add: authtype, username, password,school, notification url, secret, refresh interval

// Global shutdown signal channel
var shutdownSignal chan os.Signal

func main() {
	// Initialize shutdown signal channel
	shutdownSignal = make(chan os.Signal, 1)

	// Load config from env file
	slog.Info("Attempting to load config from .env file. This operation is not relevant for docker environments.")
	if err := godotenv.Load(); err != nil {
		slog.Warn("An error occurred while loading the .env file", "error", err)
	}
	slog.Info("UNTIS_URL:", os.Getenv("UNTIS_URL"))
	slog.Info("UNTIS_SCHOOL:", os.Getenv("UNTIS_SCHOOL"))
	slog.Info("UNTIS_USERNAME:", os.Getenv("UNTIS_USERNAME"))
	slog.Info("UNTIS_REFRESH_INTERVAL:", os.Getenv("UNTIS_REFRESH_INTERVAL"))
	slog.Info("NOTIFICATION_URL:", os.Getenv("NOTIFICATION_URL"))
	slog.Info("DB_LOCATION:", os.Getenv("DB_LOCATION"))

	changeDetector := utils.NewChangeDetector()

	slog.Info("Loading session from", "path", utils.GetDbPath())

	// open db
	lessonDB := utils.OpenDatabase("timetable")
	homeworkDB := utils.OpenDatabase("homework")

	slog.Info("connection opened to databases")

	// populate lessons and homework with outdated lessons and homework
	if utils.PopulateLessons(&changeDetector, lessonDB) != nil ||
		utils.PopulateHomework(&changeDetector, homeworkDB) != nil {
		return
	}
	slog.Info("Loaded session from databases")

	// fetch new lesson & homework data and check for changes
	go utils.CheckForUpdate(&changeDetector, lessonDB, homeworkDB)

	// handle shutdown signals gracefully
	utils.HandleShutdownSignals(shutdownSignal, lessonDB, homeworkDB)
}

// ENV variables:
// UNTIS_URL
// UNTIS_SCHOOL
// UNTIS_USERNAME
// UNTIS_PASSWORD
// UNTIS_REFRESH_INTERVAL
// NOTIFICATION_URL
// DB_LOCATION
