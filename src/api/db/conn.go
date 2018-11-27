// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com
package db

import (
	"api/logger"
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	. "github.com/polaris1119/config"
	"os"
	"sync"
	"time"
)

var (
	sqlLogFile    = "/data/log/test/sql.log"
	engineGroup   *xorm.EngineGroup
	dns, dnsSlave string
	dbLock        sync.RWMutex
)

func init() {
	EngineGroup()
}

func fillDns(mysqlConfig map[string]string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		mysqlConfig["user"],
		mysqlConfig["password"],
		mysqlConfig["host"],
		mysqlConfig["port"],
		mysqlConfig["dbname"],
		mysqlConfig["charset"])
}

func initEngine() error {
	if engineGroup != nil {
		return nil
	}

	var (
		err error
	)

	//主库
	masterEngine, err := xorm.NewEngine("mysql", dns)
	if err != nil {
		return err
	}
	//从库
	slaveEngine, err := xorm.NewEngine("mysql", dnsSlave)
	if err != nil {
		return err
	}

	engineGroup, err = xorm.NewEngineGroup(masterEngine, []*xorm.Engine{slaveEngine})
	if err != nil {
		return err
	}

	maxIdle := ConfigFile.MustInt("mysql", "max_idle", 2)
	maxConn := ConfigFile.MustInt("mysql", "max_conn", 10)

	engineGroup.SetMaxIdleConns(maxIdle)
	engineGroup.SetMaxOpenConns(maxConn)

	showSQL := ConfigFile.MustBool("xorm", "show_sql", false)
	logLevel := ConfigFile.MustInt("xorm", "log_level", 1)
	env := ConfigFile.MustValue("global", "env", "release")

	if env != "release" {
		file := sqlLogger(sqlLogFile)
		sqlLogFile := xorm.NewSimpleLogger(file)
		engineGroup.SetLogger(sqlLogFile)
		engineGroup.ShowSQL(showSQL)
	}
	engineGroup.Logger().SetLevel(core.LogLevel(logLevel))
	logger.Infof("database init finished, master dns: %s, slave dns: %s", dns, dnsSlave)
	fmt.Printf("database init finished, master dns: %s, slave dns: %s, engineGroup: %+v\n", dns, dnsSlave, engineGroup)

	//检查连接数据库是否正常
	err = engineGroup.Ping()
	if err != nil {
		logger.Errorf("ping database fail, %v", err)
	}

	//5s检查数据库连接是否正常
	go func() {
		for {
			timer := time.NewTimer(time.Second * 60)
			select {
			case <-timer.C:
				err := engineGroup.Ping()
				if err != nil {
					logger.Infof("ping database fail, %s, err: %v", dns, err)

					dbLock.Lock()
					//释放连接池, 重启动连接池
					engineGroup.Close()
					engineGroup = nil
					defer dbLock.Unlock()

					return
				} else {
					logger.Infof("ping database successful")
				}
			}
		}
	}()

	return nil
}

//主库从库engine group
func EngineGroup() *xorm.EngineGroup {
	dbLock.Lock()
	defer dbLock.Unlock()

	if engineGroup == nil {
		mysqlConfig, err := ConfigFile.GetSection("mysql")
		if err != nil {
			fmt.Println("get mysql config error:", err)
			panic("mysql init fail")
		}

		mysqlSlaveConfig, err := ConfigFile.GetSection("mysql_read1")
		if err != nil {
			fmt.Println("get mysql slave config err:", err)
			panic("mysql slave init fail")
		}

		dns = fillDns(mysqlConfig)
		dnsSlave = fillDns(mysqlSlaveConfig)

		// 启动时就打开数据库连接
		if err = initEngine(); err != nil {
			panic(err)
		}
	}

	return engineGroup
}

//主库engine
func StdMasterDB() *xorm.Engine {
	if engineGroup == nil {
		EngineGroup()
	}

	return engineGroup.Master()
}

//从库engine
func StdSlaveDB() *xorm.Engine {
	if engineGroup == nil {
		EngineGroup()
	}

	return engineGroup.Slave()
}

func sqlLogger(sqlLogFile string) *os.File {
	file, err := os.OpenFile(sqlLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("init sql.log file fail, err = %v", err)
		return nil
	}

	return file
}
