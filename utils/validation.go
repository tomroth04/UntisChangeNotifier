package utils

import (
	untisApi "github.com/tomroth04/untisAPI"
	"github.com/tomroth04/untisAPI/types"
	"time"
)

// ValidateEndDate validates the end date
// if the end date surpasses the end of the school year it will be set to the end of the school year
// this is to prevent the user from entering a date that is not in the school year and causing an error
func ValidateEndDate(client untisApi.Client) time.Time {
	endDate := time.Now().Add(time.Hour * 24 * 30)
	if year, err := client.GetLatestSchoolyear(false); err == nil {
		if year.EndDate.Before(types.Time(time.Now().Add(time.Hour * 24 * 30))) {
			endDate = time.Date(year.EndDate.ToTime().Year(), year.EndDate.ToTime().Month(),
				year.EndDate.ToTime().Day(), 0, 0, 0, 0, time.Local)
		}
	}
	return endDate
}

// TODO: also add validation for the start date
