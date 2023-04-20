package model

import "time"

type EntryTime struct {
	StartTime time.Time `json:"start_time`
	PhotoId   int       `json:"photo_id"`
}
type ExitTime struct {
	EndTime time.Time `json:"end_time"`
	PhotoId int       `json:"photo_id"`
}
type TimeRecord struct {
	Id        int       `json:"id"`
	Employee  *int      `json:"employee"`
	EntryTime EntryTime `json:"entryTime"`
	ExitTime  *ExitTime `json:"exitTime"`
}

type AddTimeRecord struct {
	EmployeeId *int      `json:"employee_id"`
	EntryTime  EntryTime `json:"entryTime"`
}

func (t AddTimeRecord) ToTimeRecord(id int) TimeRecord {
	timeRec := TimeRecord{Id: id, Employee: t.EmployeeId, EntryTime: t.EntryTime}

	return timeRec

}

type UpdateTimeRecord struct {
	Id        int       `json:"id"`
	Employee  *int      `json:"employee"`
	EntryTime EntryTime `json:"entryTime"`
	ExitTime  *ExitTime `json:"exitTime"`
}

func (t UpdateTimeRecord) ToUpdateTimeRecord(id int) TimeRecord {
	timeRec := TimeRecord{Id: id, Employee: t.Employee, EntryTime: t.EntryTime, ExitTime: t.ExitTime}

	return timeRec

}

type DateTime time.Time
