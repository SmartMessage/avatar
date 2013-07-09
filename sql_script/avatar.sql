----------------------create user and database and grant--------------------------
mysql -u root -p   --wolf

insert into mysql.user(Host,User,Password) values("localhost","avatar",password("avatar"));

flush privileges;

exit;
mysql -u root -p 

create database avatar;

grant all privileges on avatar.* to avatar@% identified by 'avatar';

flush privileges;
exit;

mysql -hlocalhost -uavatar -p --avatar

-----------------------create tables --------------------------------------------

drop table job_timer;
drop table job_time_window;
drop table job_dependency;
drop table job_stream;
drop table job_status;
drop table job_callback_event_log;
drop table job_callback_event;
drop table job_log;
drop table service_monitor_cmd;
drop table service_monitor_log;
drop table server_monitor_log;
drop table avatar_log;
drop table avatar_user;
drop table job_base;
drop table server_info;


create table server_info(
  `id` int(11) NOT NULL AUTO_INCREMENT,  
  `server_name` varchar(20) NOT NULL,  
  `user_id` int(11) NOT NULL,   
  `enabled`  enum('Y','N') DEFAULT 'Y',  
  `host_name` varchar(40) NOT NULL,  
  `server_ip` varchar(15) NOT NULL,
  `create_date`  date, 
  `last_change_time` datetime,
  PRIMARY KEY (`id`)  
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;

create table job_base(
  `id` int(11) NOT NULL AUTO_INCREMENT,  
  `job_name` varchar(20) NOT NULL,  
  `user_id` int(11) NOT NULL,  
  `is_timing_job` enum('Y','N') DEFAULT 'N',  
  `job_type` enum('Common','MapReduce','Hive','Sqoop','Monitor') DEFAULT 'Common',  
  `job_cycle` enum('M','D','W','M','Q','O') DEFAULT 'D',  
  `job_status` enum('Ready','Waiting','Assigned','Running','Failed','Done')  DEFAULT 'Ready',
  `day_offset` int(2) DEFAULT 0,
  `server_id` int(11) NOT NULL ,
  `job_commond` varchar(2000),  
  `enabled`  enum('Y','N') DEFAULT 'Y',  
  `need_callback` enum('Y','N') DEFAULT 'N',
  `create_date` date NOT NULL,
  `last_run_date`  date, 
  `last_change_time` datetime,
  PRIMARY KEY (`id`),
  foreign key (`server_id`) references  server_info (`id`)  on delete cascade
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;


create table job_timer(
  `id` int(11) NOT NULL AUTO_INCREMENT,  
  `minute_interval` int(11),
  `day_timer` varchar(6),
  `week_timer` varchar(7),
  `is_end_month`  enum('Y','N') DEFAULT 'N',
  `month_timer` varchar(31),
  `quarterly_timer` enum('1','2','3','4'),
  `once_job_timer` varchar(19),
  `job_id` int(11) NOT NULL,
  `last_change_time` datetime,
  PRIMARY KEY (`id`),
  foreign key (`job_id`) references  job_base (`id`)  on delete cascade
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;


create table job_time_window(
  `id` int(11) NOT NULL AUTO_INCREMENT,  
  `begin_hour` varchar(2),
  `end_hour` varchar(2),
  `job_id` int(11) NOT NULL,
  `last_change_time` datetime,
  PRIMARY KEY (`id`),
  foreign key (`job_id`) references  job_base (`id`)  on delete cascade
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;



create table job_dependency(
  `id` int(11) NOT NULL AUTO_INCREMENT,  
  `job_id` int(11) NOT NULL,
  `dependent_job_id` int(11) NOT NULL,
  `enabled`  enum('Y','N') DEFAULT 'Y' ,
  `last_change_time` datetime,
  PRIMARY KEY (`id`),
  foreign key (`job_id`) references  job_base (`id`)  on delete cascade
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;


create table job_stream(
  `id` int(11) NOT NULL AUTO_INCREMENT,  
  `job_id` int(11) NOT NULL,
  `downstream_job_id` int(11) NOT NULL,
  `enabled`  enum('Y','N') DEFAULT 'Y' ,
  `last_change_time` datetime,
  PRIMARY KEY (`id`),
  foreign key (`job_id`) references  job_base (`id`)  on delete cascade
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;




create table job_status(
  `id` int(11) NOT NULL AUTO_INCREMENT,  
  `job_id` int(11) NOT NULL,
  `job_status` enum('Ready','Waiting','Assigned','Running','Failed','Done'),
  `last_run_date`  date,  
  `last_change_time` datetime,
  PRIMARY KEY (`id`),
  foreign key (`job_id`) references  job_base (`id`)  on delete cascade
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;


create table job_callback_event(
  `id` int(11) NOT NULL AUTO_INCREMENT,  
  `job_id` int(11) NOT NULL,
  `enabled`   enum('Y','N') DEFAULT 'Y' ,
  `event_type` enum('HTTP','Script') DEFAULT 'HTTP',  
  `event` varchar(2000),
  `last_change_time` datetime,
  PRIMARY KEY (`id`),
  foreign key (`job_id`) references  job_base (`id`)  on delete cascade
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;


create table job_callback_event_log(
  `id` int(11) NOT NULL AUTO_INCREMENT,  
  `event_id` int(11) NOT NULL,
  `event_status` int(2),
  `event_log` varchar(2000),
  `last_change_time` datetime,
  PRIMARY KEY (`id`),
  foreign key (`event_id`) references  job_callback_event (`id`)  on delete cascade
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;

create table job_log(
  `id` int(11) NOT NULL AUTO_INCREMENT,  
  `job_id` int(11) NOT NULL,
  `run_job_log` varchar(2000),
  `last_change_time` datetime,
  PRIMARY KEY (`id`),
  foreign key (`job_id`) references  job_base (`id`)  on delete cascade
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;




create table service_monitor_cmd(
  `id` int(11) NOT NULL AUTO_INCREMENT,  
  `job_id` int(11) NOT NULL,
  `monitor_cmd` varchar(2000),
  `recovery_cmd` varchar(2000),
  `send_email` varchar(1000),
  `last_change_time` datetime,
  PRIMARY KEY (`id`),
  foreign key (`job_id`) references  job_base (`id`)  on delete cascade
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;



create table service_monitor_log(
  `id` int(11) NOT NULL AUTO_INCREMENT,  
  `job_id` int(11) NOT NULL,
  `status_log` enum('Normal','Error','Restart') DEFAULT 'Normal',
  `last_change_time` datetime,
  PRIMARY KEY (`id`),
  foreign key (`job_id`) references  job_base (`id`)  on delete cascade
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;


create table server_monitor_log(
  `id` int(11) NOT NULL AUTO_INCREMENT, 
  `job_id` int(11) NOT NULL,
  `host_name` varchar(100),
  `server_ip` varchar(15) NOT NULL,
  `cpu_utilization` decimal(3, 2), 
  `cpu_load` decimal(3, 2), 
  `mem_utilization` decimal(3, 2), 
  `disk_utilization` decimal(3, 2), 
  `disk_io` decimal(3, 2), 
  `network_io` decimal(3, 2), 
  `last_change_time` datetime,
  PRIMARY KEY (`id`),
  foreign key (`job_id`) references  job_base (`id`)  on delete cascade
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;


create table avatar_log(
  `id` int(11) NOT NULL AUTO_INCREMENT,  
  `goroutine_name` varchar(100),
  `status_log` enum('Normal','Error','Restart') DEFAULT 'Normal',
  `last_change_time` datetime,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;


  create table avatar_user(
    `user_id` int(11) NOT NULL AUTO_INCREMENT, 
    `user_name` varchar(20) NOT NULL, 
    `user_mobile_phone` int(11),
    `role_type` enum('common','admin') NOT NULL DEFAULT 'common',
    `last_change_time` datetime,
    PRIMARY KEY (`user_id`)
  ) ENGINE=InnoDB  DEFAULT CHARSET=utf8;