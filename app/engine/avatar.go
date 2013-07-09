package engine

import (
	"github.com/robfig/revel"
	"time"
	model "avatar/app/models"
)

const (
	TRUE = "Y"
	FALSE = "N"
	TIMER = "job_timer"
	SCHEDULE = "job_schedule"
	START = "start"
	RESTART = "restart"
	STOP = "stop"
	ALIVE = "Alive"
	DEAD = "Dead"
)

var (
	chanTimer chan(*model.Job)
	chanReceiver chan(*model.Job)
	chanNotification chan(string)
	heartbeatTimeout time.Duration  //秒级
	heartbeatInterval time.Duration  //秒级
	engineRunFlag bool
	mapGoroutineTime map[string]time.Time
)

//心跳计时器
func heartbeatMonitor(){
	go func(){
		for engineRunFlag {
			for goroutine, lastHeartbeatTime := range mapGoroutineTime {
				switch (goroutine) {
				case TIMER:
					if time.Now().After(lastHeartbeatTime.Add(heartbeatTimeout)) {
						revel.WARN.Println("JobTimer    heartbeat Timeout, So start It now..........")
						startJobTimer()
						insertAvatarMonitorLog(goroutine,RESTART)
					}else{
						insertAvatarMonitorLog(goroutine,ALIVE)
						revel.INFO.Printf("[goroutine: %s   ]--->[status: %s]....................\n",TIMER, ALIVE)  
					}
				case SCHEDULE:
					if time.Now().After(lastHeartbeatTime.Add(heartbeatTimeout)) {
						revel.WARN.Println("JobScheduler heartbeat Timeout, So start It now............")
						startJobScheduler()
						insertAvatarMonitorLog(goroutine,RESTART)
					}else{
						insertAvatarMonitorLog(goroutine,ALIVE)
						revel.INFO.Printf("[goroutine: %s   ]--->[status: %s]....................\n",SCHEDULE, ALIVE)  
					}
				default:
					if time.Now().After(lastHeartbeatTime.Add(heartbeatTimeout)) {
						revel.WARN.Printf("%s     heartbeat Timeout..............................\n",goroutine)
					}else{
						revel.INFO.Printf("%s send heartbeat %s ..............................\n",goroutine, ALIVE)
					}
				}
				time.Sleep(heartbeatInterval)
			}
		}
	}()
}

//接收goroutine的心跳报告(含远程Na'vi机器的心跳) 
func ReceiveGoroutineHeartbeat(goroutine string, serverReport *model.ServerReport){
	mapGoroutineTime[goroutine] = time.Now()
	if serverReport != nil {
		//记录远端 Na'vi 机器的机器健康报告 默认：每10次心跳发送一次服务器健康报告
	}
}

//启动作业定时器，以制定间隔获取数据库中满足执行条件的时间触发类型的作业
func startJobTimer(){
	go func() {
		defer func(){
			if e := recover(); e != nil {
				revel.ERROR.Println("goroutine jobTimer error:",e)
			}
		}()
		for {
			if !engineRunFlag {
				chanNotification <- TRUE
				return
			}
			m := QueryTimingJobs()
			for _, job := range m {
				//判断作业是否满足时间触发的条件
				flag :=  judgmentJobTime(job)
				if !flag {
					//revel.INFO.Printf("Judge not by time, job[%d] .....................",job.jobId)
					continue
				}

				flag = updateJob2Waiting(job.JobId, job.SessionId, job.LastRunDate)
				if flag {
					revel.INFO.Printf("update job[%d] Waiting, push it to chanTimer......",job.JobId)
					chanTimer <- job
				}else{
					revel.WARN.Printf("update job[%d] Waiting failed.....................",job.JobId)
				}
			}
			time.Sleep(heartbeatInterval)
			ReceiveGoroutineHeartbeat(TIMER, nil)
		}
	}()
}

//接收远端Na'vi机器发送回来的作业执行结果信号
func NaviSignalReceiver(jobId, status, runlog string) {
	defer func(){
		if e := recover(); e != nil {
			revel.ERROR.Println("NviSignalReceiver error,................err:",e)
		}
	}()
	// 更新作业状态

	//插入作业执行日志

	if engineRunFlag {
	//查找作业的下游作业，循环放入channel 并更循环到的下游作业状态及插入status表数据

		chanReceiver <- &model.Job{JobId:1}
	}

}

//启动作业调度器，从channel中获取waiting状态的作业进行调度
func startJobScheduler(){
	go func (){
		defer func(){
			if e := recover(); e != nil {
				revel.ERROR.Println("job scheduler error,................err:",e)
			}
			//退出前将chan中的Waiting作业处理完成
			for {
				rJob, ok := <-chanReceiver
				if ok {
					executeJob(rJob)
				}
			}
			for {
				tJob, ok := <-chanTimer
				if ok {
					executeJob(tJob)
				}
			}
		}()

		for {
			select {
			case tJob, ok := <-chanTimer:
				if ok {
					executeJob(tJob)
				}
			case rJob, ok := <-chanReceiver:
				if ok {
					executeJob(rJob)
				}
			case _, ok := <-chanNotification:
				if ok {
					return
				}
			default:
				//revel.INFO.Println("job scheduler sleeping...............................")
				time.Sleep(heartbeatInterval)
				ReceiveGoroutineHeartbeat(SCHEDULE, nil)
			}

		}
	}()
}

//处理运行完成的作业的回调事件
func execCallbackEvent(jobId string){

}

//执行作业：将作业发送至远端机器的Na'vi
func executeJob(job *model.Job){
	revel.INFO.Println("execute job begin.............",job)
}

//引擎初始化
func initialization() bool {
	revel.INFO.Println("begin initialze avatar engine ................................")
	defer func() bool {
		if e := recover(); e != nil {
			revel.ERROR.Printf("initialize ......................................err:\n%v",e)
			return false
		}
		return true
	}()

	chanTimer = make (chan *model.Job,5)
	chanReceiver = make (chan *model.Job,5)
	chanNotification = make (chan string)
	mapGoroutineTime = make (map[string]time.Time)

	heartbeatTimeout = time.Second * time.Duration(6)
	heartbeatInterval = time.Second * time.Duration(3)

	//初始化待监控的服务器(goroutine) 列表
	mapGoroutineTime[TIMER] = time.Now()
	mapGoroutineTime[SCHEDULE] = time.Now()

	return true
}

func ControlAvatar(cmd string){
	switch cmd {
	case START :
		engineRunFlag = true
	case STOP :
		engineRunFlag = false
		default :
		revel.WARN.Printf("commond [%s] error .................................", cmd)
	}
}

//启动引擎
func StartEngine(){
	flag := initialization()
	//NaviSignalReceiver("","","")
	startJobTimer()
	startJobScheduler()
	heartbeatMonitor()
	revel.INFO.Println("initialze avatar engine, the result is....................",flag) 
}

/*
待优化列表：
3、 mysql数据库报错 连接数过多
*/
