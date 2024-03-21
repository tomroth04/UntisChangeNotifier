package utils

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	. "github.com/tomroth04/untisAPI/types"
	"log/slog"
	"strconv"
	"time"
)

type UntisChangeDetector struct {
	Lessons   map[int]GenericLesson
	Homeworks map[int]Homework
}

// NewChangeDetector creates a new detector to detect changes in the timetable
func NewChangeDetector() UntisChangeDetector {
	return UntisChangeDetector{
		Lessons:   map[int]GenericLesson{},
		Homeworks: map[int]Homework{},
	}
}

// CheckLessons checks for changes in the lessons and notifies the user if there are any
func (u *UntisChangeDetector) CheckLessons(lessons []GenericLesson, db *leveldb.DB) {
	batch := new(leveldb.Batch)
	// Match lessons to already existing lessons
	for _, lesson := range lessons {
		e, existing := u.Lessons[lesson.GetLessonId()]
		if existing {
			if !lesson.IsEqual(e) {
				CheckLessonChanges(u, lesson, false)
				u.Lessons[lesson.GetLessonId()] = lesson
			}
		} else {
			u.Lessons[lesson.GetLessonId()] = lesson
			CheckLessonChanges(u, lesson, true)
		}
		// store lesson in db
		batch.Put([]byte(strconv.Itoa(lesson.GetLessonId())), []byte(lesson.R.String()))
	}
	// persist changes in batch to db
	if err := db.Write(batch, nil); err != nil {
		slog.Error("error writing lessons batch update", "error", err)
	}

	// delete outdated lessons
	for k, value := range lessons {
		if value.GetDate().Before(time.Now().Add(time.Hour * -48)) { // Check if lesson is older than 48 hours
			err := db.Delete([]byte(strconv.Itoa(value.GetLessonId())), nil)
			if err != nil {
				slog.Error("error deleting lesson from db", "error", err)
			}
			delete(u.Lessons, k)
		}
	}
}

func ToJson(i interface{}) []byte {
	b, err := json.Marshal(i)
	if err != nil {
		slog.Error("error converting to json", "error", err)
		return []byte{}
	}
	return b
}

// CheckHomework checks for changes in the homework and notifies the user if there are any
func (u *UntisChangeDetector) CheckHomework(homeworks []Homework, db *leveldb.DB) {
	batch := new(leveldb.Batch)

	// Match homeworks to already existing homeworks
	for _, homework := range homeworks {
		_, existing := u.Homeworks[homework.Id]
		if !existing {
			slog.Info("New homework", "homework", homework)
			go notify(
				"New Homework: "+homework.LessonName,
				getHomeworkNotificationFormat(homework),
			)
		}
		// store homework in db
		batch.Put([]byte(strconv.Itoa(homework.Id)), ToJson(homework))
		u.Homeworks[homework.Id] = homework
	}

	// write batch to db
	if err := db.Write(batch, nil); err != nil {
		slog.Error("error writing homework batch to db", "error", err)
	}

	// delete outdated homeworks
	for k, value := range homeworks {
		value.GetDueDate()
		if value.GetDueDate().Before(time.Now().Add(time.Hour * -48)) { // Check if homework is older than 48 hours
			delete(u.Lessons, k)
			if err := db.Delete([]byte(strconv.Itoa(value.Id)), nil); err != nil {
				slog.Error("error deleting homework from db", "error", err)
			}
		}
	}
}

// TODO: remove the firstTime parameter and use the existing lesson map to check didn't exist before
// CheckLessonChanges checks if the lesson has been cancelled or replaced and notifies the user if so
func CheckLessonChanges(detector *UntisChangeDetector, lesson GenericLesson, firstTime bool) {
	if !firstTime {
		slog.Info("Detected change in lesson", "lesson", lesson.ToString())
	}
	if lesson.IsCancelled() {
		go notify(
			"Free: "+lesson.GetSubject(),
			getLessonNotificationFormat(lesson),
		)
	} else if lesson.IsIrregular() {
		go notify(
			"Irregular: "+lesson.GetSubject(),
			getLessonNotificationFormat(lesson),
		)
	} else if lesson.IsReplaced() && detector.Lessons[lesson.GetLessonId()].IsReplaced() {
		go notify(
			"Change supervision: "+lesson.GetSubject(),
			getLessonNotificationFormat(lesson),
		)
	} else if lesson.IsReplaced() {
		go notify(
			"Supervision: "+lesson.GetSubject(),
			getLessonNotificationFormat(lesson),
		)
	} else if !firstTime {
		go notify(
			"Change: "+lesson.GetSubject(),
			getLessonNotificationFormat(lesson),
		)
	}
}
