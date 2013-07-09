package db_access

import (
	"database/sql"
	"fmt"
	"github.com/robfig/revel"
	_ "github.com/go-sql-driver/mysql"
	model "avatar/app/models"
)

const (

)


//获取所有时间触发类型的作业
func QueryUserByCondition(conditions string) (mapUser map[int]*model.AvatarUser) {
	var db *sql.DB
	var err error
	defer func (){
		if e := recover(); e != nil {
			revel.ERROR.Println("query all jobs error: ",e)
		}
		closeDBConn(db)
	}()
	mapUser = make(map[int]*model.AvatarUser)
	db, err = getDBConn()
	if err != nil {
		revel.ERROR.Println("get database conn error: ",err.Error())
		return
	}

	sql := "SELECT user_id,user_name,user_passwd,user_mobile_phone,role_type FROM  avatar_user WHERE 1=1 " + conditions
	rows, eq := db.Query(sql)

	if eq != nil {
		revel.ERROR.Println("query jobs error: ",eq.Error())
		return
	}

	for rows.Next() {
		user := new(model.AvatarUser)
		row_err := rows.Scan(&user.UserId,&user.UserName, &user.UserPasswd, &user.UserMobilePhone, &user.RoleType)
		if row_err != nil {
			revel.ERROR.Println("row scan error: ",row_err.Error())
			return
		}
		mapUser[user.UserId] = user
	}

	return
}

//查询用户信息
func QureyByUserName(userName string) (user *model.AvatarUser) {
	var db *sql.DB
	var err error
	defer func (){
		if e := recover(); e != nil {
			revel.ERROR.Println("query user info error: ",e)
		}
		closeDBConn(db)
	}()
	db, err = getDBConn()
	if err != nil {
		revel.ERROR.Println("get database conn error: ",err.Error())
		return
	}
	user = new(model.AvatarUser)
	sql := fmt.Sprintf("SELECT user_id,user_name,user_passwd,user_mobile_phone,role_type FROM  avatar_user WHERE user_name = '%s' ", userName)
	row := db.QueryRow(sql)
	row_err := row.Scan(&user.UserId,&user.UserName, &user.UserPasswd, &user.UserMobilePhone, &user.RoleType)
	if row_err != nil {
		revel.ERROR.Println("query user, row scan error: ",row_err.Error())
		return
	}
	return
}

