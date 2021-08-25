package model

import "time"

type CheckInRecord struct {
	id              int       `json:"id"`
	Name            string    `json:"name"` //注意:struct名稱開頭必須要大寫...否則無法寫入mongoDB!!!不知道為什麼...
	Check_in_time   string    `json:"checkintime"`
	Pic             string    `json:"pic"`
	Leave_type      string    `json:"leavetype"`
	Date            string    `json:"date"`
	Department      string    `json:"department"`
	Position        string    `json:"position"`
	DateTimeToday   time.Time `json:"datetimetoday"`
	DateTimeCheckIn time.Time `json:"datetimecheckin"`
}
