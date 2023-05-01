package dto

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type ProfileWorkoutCreateModel struct {
	WorkoutPlan string          `json:"workoutPlan"`
	Schedule    []DaysOfWeekDto `json:"schedule"`
}

type ProfileWorkoutShortModel struct {
	Id string `json:"id"`
}

type DaysOfWeekDto int

const (
	Sunday DaysOfWeekDto = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func (s DaysOfWeekDto) String() string {
	return toString[s]
}

var toString = map[DaysOfWeekDto]string{
	Sunday:    "Sunday",
	Monday:    "Monday",
	Tuesday:   "Tuesday",
	Wednesday: "Wednesday",
	Thursday:  "Thursday",
	Friday:    "Friday",
	Saturday:  "Saturday",
}

var toID = map[string]DaysOfWeekDto{
	"Sunday":    Sunday,
	"Monday":    Monday,
	"Tuesday":   Tuesday,
	"Wednesday": Wednesday,
	"Thursday":  Thursday,
	"Friday":    Friday,
	"Saturday":  Saturday,
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
	if value, ok := toID[j]; ok {
		*s = value
		return nil
	}
	return fmt.Errorf("can't convert value %s to valid day of week", j)
}
