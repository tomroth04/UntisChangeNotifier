package utils

import (
	"github.com/containrrr/shoutrrr"
	untisApi "github.com/tomroth04/untisAPI/types"
	"log/slog"
	"os"
	"strings"
)

type NotificationData struct {
	Topic    string `json:"topic"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

// getHomeworkNotificationFormat formats the homework notification
func getHomeworkNotificationFormat(homework untisApi.Homework) string {
	var ans strings.Builder
	ans.WriteString("Subject: ")
	ans.WriteString(homework.LessonName)
	ans.WriteString("\n")
	ans.WriteString("Date: ")
	ans.WriteString(homework.GetDueDate().Format("Monday, 02 January 2006"))
	ans.WriteString("\n")
	ans.WriteString("Description: ")
	ans.WriteString(homework.Text)
	if homework.Remark != "" {
		ans.WriteString("\nRemark: ")
		ans.WriteString(homework.Remark)
	}
	return ans.String()
}

// Format the lesson notification
func getLessonNotificationFormat(lesson untisApi.GenericLesson) string {
	var ans strings.Builder
	ans.WriteString("Subject: ")
	ans.WriteString(lesson.GetSubject())
	ans.WriteString("\n")
	ans.WriteString("Date: ")
	ans.WriteString(lesson.GetDateFormatted())
	ans.WriteString("\n")
	ans.WriteString("Time: ")
	ans.WriteString(lesson.GetStartTimeFormatted())

	// Include lesson info and substitute text if available
	if lesson.GetLessonInfo() != "" {
		ans.WriteString("\nLesson-Text: ")
		ans.WriteString(lesson.GetLessonInfo())
	}
	if lesson.GetSubstituteText() != "" {
		ans.WriteString("\nSubstitute-Text: ")
		ans.WriteString(lesson.GetSubstituteText())
	}

	return ans.String()
}

func notify(title string, message string) {
	err := shoutrrr.Send(os.Getenv("NOTIFICATION_URL")+"?title="+title, message)
	if err != nil {
		slog.Error("error sending notification", "error", err)
		return
	}
}
