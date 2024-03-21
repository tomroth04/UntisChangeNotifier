package utils

import (
	"github.com/syndtr/goleveldb/leveldb"
	untisApi "github.com/tomroth04/untisAPI"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// CheckForUpdate gets data from Untis and check for changes
func CheckForUpdate(changeDetector *UntisChangeDetector, lessonDB *leveldb.DB, homeworkDB *leveldb.DB) {
	count := 0
	for {
		if count != 0 {
			time.Sleep(SleepTime(time.Now()))
		}
		count++
		client := untisApi.NewClient(
			os.Getenv("UNTIS_URL"),
			os.Getenv("UNTIS_SCHOOL"),
			"ANDROID",
			os.Getenv("UNTIS_USERNAME"),
			os.Getenv("UNTIS_PASSWORD"),
		)
		// "antiope.webuntis.com", "lmrl", "ANDROID", "RotTo652", "3S2QRV5IJQR4E5ZR"
		if err := client.Login(); err != nil {
			slog.Error("logging in failed", "error", err)
			continue
		}

		endDate := ValidateEndDate(client)

		// check timetable for the next 30 days
		lessons, err := client.GetOwnTimetableForRange(time.Now(), endDate, false)
		if err != nil {
			slog.Error("getting timetable failed", "error", err)
			continue
		}
		changeDetector.CheckLessons(lessons, lessonDB)

		// check homework for the next 30 days
		homeworks, err := client.GetHomeworksFor(time.Now(), time.Now().Add(time.Hour*24*30), false)
		if err != nil {
			slog.Error("getting homework failed", "error", err)
			continue
		}
		changeDetector.CheckHomework(homeworks, homeworkDB)

		if err := client.Logout(); err != nil {
			slog.Error("logging out failed", "error", err)
			continue
		}
		//slog.Info("refreshed data")
		runtime.GC()
	}
}
