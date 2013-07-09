package models

import (
	//"database/sql"
	//"fmt"
)

type Job struct {
	JobId int				`json:"id"`
	JobName string			`json:"job_name"`
	UserId int				`json:"user_id"`
	IsTimingJob string 		`json:"is_timing_job"`
	JobType string			`json:"job_type"`
	JobCycle string			`json:"job_cycle"`
	JobStatus string		`json:"job_status"`
	DayOffset int			`json:"day_offset"`
	SessionId int			`json:"session_id"`
	ServerId int			`json:"server_id"`
	JobCommond string		`json:"job_commond"`
	Enabled string			`json:"enabled"`
	NeedCallback string		`json:"need_callback"`
	CreateDate string		`json:"create_date"`
	LastRunDate  string		`json:"last_run_date"`
	LastChangeTime string   `json:"last_change_time"`

}

type JobTimer struct {
	Id int					`json:"id"`
	DayTimer string			`json:"day_timer"`
	WeekTimer string		`json:"week_timer"`
	IsEndMonth string		`json:"is_end_month"`
	MonthTimer string		`json:"month_timer"`
	QuarterlyTimer string   `json:"quarterly_timer"`
	OnceJobTimer string		`json:"once_job_timer"`
	JobId int				`json:"job_id"`
}

type JobTimeWindow struct {

}

type JobDependency struct {

}

type JobStream struct {

}

type JobStatus struct {

}

type JobCallbackEvent struct {

}

type JobCallbackEventLog struct {

}

type JobLog struct {

}


