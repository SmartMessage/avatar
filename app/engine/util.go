package engine

import (
	"time"
	"fmt"
	"github.com/robfig/revel"
	model "avatar/app/models"
	dba "avatar/app/db_access"
)

const (
	READY = "Ready"
	WAITING = "Waiting"
	ASSIGNED = "Assigned"
	RUNNING = "Running"
	FAILED = "Failed"
	DONE = "Done"
	AVATAR = "Avatar"
	DATE_FORMAT = "2006-01-02"
	WITHNANOS = "2006-01-02 15:04:05"
	YMD = " 00:00:00"
)

var (

)

// 查询所有时间触发的作业
func QueryTimingJobs() (mapJobs map[int]*model.Job){
	conditions := "and is_timing_job='N' and enabled='Y' and last_run_date<'" + GetCurrDate() + "' and job_status in ('" + READY + "','" + DONE + "')"
	return dba.QueryJobsByCondition(conditions)
}

// 查询所有job
func QueryAllJobs() (mapJobs map[int]*model.Job){
	conditions := ""
	return dba.QueryJobsByCondition(conditions)
}

//查询时间触发的时间设置
func QureyTimerByJobId(jobId int) *model.JobTimer {
	return dba.QureyTimerByJobId(jobId)
}

//更新作业状态为Waiting
func updateJob2Waiting(jobId int, sessionId int, lastRunDate string) bool {
	txDate := GetAfterDate(lastRunDate)
	updateSql := fmt.Sprintf("update job_base set job_status='Waiting', session_id=session_id+1 ,last_run_date= '%s' job_id = %d",txDate,jobId)
	flag := dba.UpdateJob(updateSql)
	if flag {
		insertJobStatus(jobId, sessionId, WAITING, txDate)
	}
	return flag
}
//插入作业状态
func insertJobStatus(jobId int, sessionId int, jobStatus string, txDate string) {
	dba.InsertJobStatus(jobId,sessionId, jobStatus, txDate)
}

//插入avatar goroutine 监控日志
func insertAvatarMonitorLog(goroutineName,status string){
	insertServiceMonitorLog(AVATAR,goroutineName,status,1)
}

//插入监控日志
func insertServiceMonitorLog(serviceType string,serviceName string,serviceStatus string, jobId int) {
	dba.InsertServiceMonitorLog(serviceType, serviceName, serviceStatus, jobId)
}

//时间判断
func judgmentJobTime(job *model.Job) bool {
	defer func (){
		if e := recover(); e != nil {
			revel.ERROR.Printf("judgment job time error , job is: \n %v \n error is: %v",job,e)
		}
	}()

	/*
	jobId := job.JobId
	jobStatus := job.JobStatus
	jobCycle := job.JobCycle
	jobTxDate := job.LastRunDate
	currDate := GetCurrDate()

	timer := QureyTimerByJobId(jobId)
	*/

	//judg
	return false
}

//获取当前时间
func GetCurrTime() time.Time {
	return time.Now()
}

//获取当前日期
func GetCurrDate() string {
	return time.Now().Format(DATE_FORMAT)
}

//获取前一天日期
func GetBeforeDate(date string) string {
	t, _ := time.Parse(WITHNANOS, date + YMD)
	return t.AddDate(0,0,-1).Format(DATE_FORMAT)
}

//获取后一天日期
func GetAfterDate(date string) string {
	t, _ := time.Parse(WITHNANOS, date + YMD)
	return t.AddDate(0,0,1).Format(DATE_FORMAT)

}
