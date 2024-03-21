package utils

import (
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tidwall/gjson"
	"github.com/tomroth04/untisAPI/types"
	"log/slog"
	"os"
	"strconv"
)

// GetDbPath gets the path to db
func GetDbPath() string {
	// Check if running under windows or macos -> use local dir for db
	if os.Getenv("DB_LOCATION") != "" {
		if os.Getenv("DB_LOCATION") == "WD" {
			return ""
		}
		return os.Getenv("DB_LOCATION")
	}
	return "/untis/"
}

// OpenDatabase opens a database with the given suffix
func OpenDatabase(suffix string) *leveldb.DB {
	db, err := leveldb.OpenFile(GetDbPath()+suffix, nil)
	if err != nil {
		slog.Error("error opening "+suffix, "error", err)
		os.Exit(1)
	}
	return db
}

// PopulateLessons populates the lessons map in the change detector with the lessons from the db
func PopulateLessons(changeDetector *UntisChangeDetector, db *leveldb.DB) error {
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		lesson := types.GenericLesson{}
		lesson.R = gjson.ParseBytes(value)

		i, err := strconv.Atoi(string(key))
		if err != nil {
			slog.Error("error converting string to int during lesson population", "error", err)
			continue
		}

		changeDetector.Lessons[i] = lesson
	}
	iter.Release()
	if iter.Error() != nil {
		slog.Error("error populating lessons", "error", iter.Error())
		return iter.Error()
	}

	return nil
}

// PopulateHomework populates the homework map in the change detector with the homework from the db
func PopulateHomework(changeDetector *UntisChangeDetector, homeDB *leveldb.DB) error {
	iter := homeDB.NewIterator(nil, nil)
	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		lesson := types.Homework{}
		if err := json.Unmarshal(value, &lesson); err != nil {
			slog.Error("error unmarshalling homework during population", "error", err)
			continue
		}
		i, err := strconv.Atoi(string(key))
		if err != nil {
			slog.Error("error converting str to int during homework population", "error", err)
			continue
		}
		changeDetector.Homeworks[i] = lesson
	}
	iter.Release()
	if iter.Error() != nil {
		slog.Error("error populating homework", "error", iter.Error())
		return iter.Error()
	}

	return nil
}
