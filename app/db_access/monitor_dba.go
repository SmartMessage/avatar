package db_access

import (
	//"database/sql"
	"github.com/robfig/revel"
	_ "github.com/go-sql-driver/mysql"
	//model "avatar/app/models"
)

const (

)

//插入service monitor log 至数据库
func InsertServiceMonitorLog(serviceType string,serviceName string,serviceStatus string, jobId int) {
	db, err := getDBConn()
	if err != nil {
		revel.ERROR.Println("get database conn error: ",err.Error())
		return
	}
	sql := "insert into service_monitor_log(service_type,service_name,status,job_id) values (?,?,?,?)"
	stmt, err := db.Prepare(sql)
	if err != nil {
		revel.ERROR.Println("prepare sql error: ",err.Error()) 
		return
	}

	_, e := stmt.Exec(serviceType,serviceName,serviceStatus,jobId)
	if e != nil {
		revel.ERROR.Println("exec sql error: ",e.Error()) 
	}
	return
}

