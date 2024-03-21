package utils

import (
	"github.com/syndtr/goleveldb/leveldb"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func CloseDatabases(lessonDB *leveldb.DB, homeworkDB *leveldb.DB) {
	if err := lessonDB.Close(); err != nil {
		slog.Error("error closing lessonDB", "error", err)
	}
	if err := homeworkDB.Close(); err != nil {
		slog.Error("error closing lessonDB", "error", err)
	}
}

func HandleShutdownSignals(shutdownSignal chan os.Signal, lessonDB *leveldb.DB, homeworkDB *leveldb.DB) {
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)
	<-shutdownSignal
	CloseDatabases(lessonDB, homeworkDB)
	slog.Info("Received shutdown signal. Cleaning up...")
	slog.Info("Cleaned up. Exiting...")
}
