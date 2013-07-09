package db_access

import (
	"database/sql"
	"strconv"
	"github.com/robfig/revel"
	_ "github.com/go-sql-driver/mysql"
	model "avatar/app/models"
)

const (

)


//获取所有时间触发类型的作业
func QueryJobsByCondition(conditions string) (mapJobs map[int]*model.Job) {
	var db *sql.DB
	var err error
	defer func (){
		if e := recover(); e != nil {
			revel.ERROR.Println("query all jobs error: ",e)
		}
		closeDBConn(db)
	}()
	mapJobs = make(map[int]*model.Job)
	db, err = getDBConn()
	if err != nil {
		revel.ERROR.Println("get database conn error: ",err.Error())
		return
	}

	jobSql := `
	SELECT 
	id,job_name,user_id,is_timing_job,job_type,job_cycle,job_status,day_offset,session_id,server_id,
	job_commond,enabled,need_callback,create_date,last_run_date,last_change_time
	FROM job_base 
	WHERE 1=1 
	` + conditions

	//revel.INFO.Println(jobSql)

	rows, eq := db.Query(jobSql)

	if eq != nil {
		revel.ERROR.Println("query jobs error: ",eq.Error())
		return
	}

	for rows.Next() {
		job := new(model.Job)
		row_err := rows.Scan(&job.JobId,&job.JobName,&job.UserId,&job.IsTimingJob,&job.JobType,&job.JobCycle,&job.JobStatus,&job.DayOffset,&job.SessionId, &job.ServerId, &job.JobCommond, &job.Enabled, &job.NeedCallback, &job.CreateDate, &job.LastRunDate, &job.LastChangeTime)
		if row_err != nil {
			revel.ERROR.Println("row scan error: ",row_err.Error())
			return
		}
		mapJobs[job.JobId] = job
	}

	return
}

//查询作业定时信息
func QureyTimerByJobId(jobId int) (timer *model.JobTimer) {
	var db *sql.DB
	var err error
	defer func (){
		if e := recover(); e != nil {
			revel.ERROR.Println("query job'timer info error: ",e)
		}
		closeDBConn(db)
	}()
	db, err = getDBConn()
	if err != nil {
		revel.ERROR.Println("get database conn error: ",err.Error())
		return
	}

	sql := `
	select id, second_interval,day_timer,week_timer,is_end_month,month_timer,quarterly_timer,once_job_timer,job_id 
	from job_timer 
	where job_id = 
	`
	sql = sql + strconv.Itoa(jobId)

	//revel.INFO.Println(sql)

	row := db.QueryRow(sql)

	timer = new(model.JobTimer)
	row_err := row.Scan(&timer.Id,&timer.DayTimer,&timer.WeekTimer,&timer.IsEndMonth,&timer.MonthTimer,&timer.QuarterlyTimer,&timer.OnceJobTimer,&timer.JobId)

	if row_err != nil {
		revel.ERROR.Println("query timer, row scan error: ",row_err.Error())
		return
	}
	return

}
//插入job至数据库
func InsertJob(job *model.Job) (flag bool){
	var db *sql.DB
	var err error
	defer func (){
		if e := recover(); e != nil {
			revel.ERROR.Println("query all jobs error: ",e)  
		}
		closeDBConn(db)
	}()

	db, err = getDBConn()
	if err != nil {
		revel.ERROR.Println("get database conn error: ",err.Error())
		return
	}
	sql := "insert into job_base(job_name,user_id,server_id,job_commond,create_date,last_run_date) values (?,?,?,?,?,?)"
	stmt, err := db.Prepare(sql)
	if err != nil {
		revel.ERROR.Println("prepare sql error: ",err.Error()) 
		return
	}

	res, err := stmt.Exec(job.JobName,job.UserId,job.ServerId,job.JobCommond,job.CreateDate,job.LastRunDate)
	if err != nil {
		revel.ERROR.Println("exec sql error: ",err.Error()) 
		return
	}

	ra,_ := res.RowsAffected()
	if ra > 0{
		return true
	}
	return
}

//更新作业信息
func UpdateJob(updateSql string) (flag bool){
	var db *sql.DB
	var err error
	defer func (){
		if e := recover(); e != nil {
			revel.ERROR.Println("query all jobs error: ",e)  
		}
		closeDBConn(db)
	}()

	db, err = getDBConn()
	if err != nil {
		revel.ERROR.Println("get database conn error: ",err.Error())
		return
	}
	stmt, err := db.Prepare(updateSql)
	if err != nil {
		revel.ERROR.Println("prepare sql error: ",err.Error()) 
		return
	}

	res, err := stmt.Exec()
	if err != nil {
		revel.ERROR.Println("exec sql error: ",err.Error()) 
		return
	}

	ra,_ := res.RowsAffected()
	if ra > 0{
		return true
	}
	return
}

//插入作业状态信息
func InsertJobStatus(jobId int, sessionId int, jobStatus string, txDate string) {
	var db *sql.DB
	var err error
	defer func (){
		if e := recover(); e != nil {
			revel.ERROR.Println("insert job status error: ",e)  
		}
		closeDBConn(db) 
	}()

	db, err = getDBConn()
	if err != nil {
		revel.ERROR.Println("get database conn error: ",err.Error())
	}
	sql := "insert into job_status(job_id,session_id,job_status,last_run_date) values (?,?,?,?)"
	stmt, err := db.Prepare(sql)
	if err != nil {
		revel.ERROR.Println("prepare sql error: ",err.Error()) 
	}

	_, err = stmt.Exec(jobId, sessionId, jobStatus, txDate)
	if err != nil {
		revel.ERROR.Println("exec sql error: ",err.Error()) 
	}
}
