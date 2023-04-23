package dto

import (
	"bytes"
	"encoding/json"
)

type ProfileWorkoutCreateModel struct {
	WorkoutPlan string          `json:"workoutPlan"`
	Schedule    []DaysOfWeekDto `json:"schedule"`
}

type ProfileWorkoutShortModel struct {
	Id string `json:"id"`
}

type DaysOfWeekDto int64

const (
	Monday DaysOfWeekDto = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

func (s DaysOfWeekDto) String() string {
	return toString[s]
}

var toString = map[DaysOfWeekDto]string{
	Monday:    "Monday",
	Tuesday:   "Tuesday",
	Wednesday: "Wednesday",
	Thursday:  "Thursday",
	Friday:    "Friday",
	Saturday:  "Saturday",
	Sunday:    "Sunday",
}

var toID = map[string]DaysOfWeekDto{
	"Monday":    Monday,
	"Tuesday":   Tuesday,
	"Wednesday": Wednesday,
	"Thursday":  Thursday,
	"Friday":    Friday,
	"Saturday":  Saturday,
	"Sunday":    Sunday,
}

// MarshalJSON marshals the enum as a quoted json string
func (s DaysOfWeekDto) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *DaysOfWeekDto) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = toID[j]
	return nil
}
