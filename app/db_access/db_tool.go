package db_access

import (
	"database/sql"
	"github.com/robfig/revel"
	_ "github.com/go-sql-driver/mysql"
)

const (
	DB_DRIVER = "mysql"
	CONN_STR = "avatar:avatar@tcp(localhost:3306)/avatar?charset=utf8"
)

//断开数据库连接
func closeDBConn(db_conn *sql.DB){
	//revel.INFO.Printf("close db conn ...........%v\n",db_conn)
	err := db_conn.Close()
	if err != nil{
		revel.ERROR.Println("close db conn error: ",err.Error())
	}
}

//创建数据库连接
func getDBConn() (db_conn *sql.DB,err error){
	defer func(){
		if e := recover(); e != nil {
			revel.INFO.Println("get database conn faild....................")
			closeDBConn(db_conn)
		}
	}()
	db_conn, e1 := sql.Open(DB_DRIVER,CONN_STR)
	if e1 != nil { 
		return db_conn, e1
	}

	_, e2 := db_conn.Query("select 1")
	if e2 != nil {
		return db_conn, e2
	}
	return db_conn,nil

}
