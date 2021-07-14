package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"regexp"
	"strings"
	"time"
)

/*
@Author : VictorTu
@Software: GoLand
*/

type DatabaseConfig struct {
	DriverName      string
	ConnectionURI   string
	MaxIdle         int
	MaxOpen         int
	ConnMaxLeftTime int
	BatchInsertSize int
	UseSingular     bool
	EnableLog       bool
}

type DatabaseInstance struct {
	gdb *gorm.DB
}

func InitInstance(config DatabaseConfig) *DatabaseInstance {
	instance := &DatabaseInstance{}
	instance.Init(config)
	return instance
}

func (instance *DatabaseInstance) Init(config DatabaseConfig) {
	db, err := gorm.Open(config.DriverName, config.ConnectionURI)
	if err != nil {
		panic(err)
	}

	db.SetLogger(&dbLogger{})

	db.DB().SetMaxIdleConns(config.MaxIdle)
	db.DB().SetMaxOpenConns(config.MaxOpen)

	maxLeftTime := time.Duration(config.ConnMaxLeftTime)
	if maxLeftTime < 30*time.Minute && maxLeftTime != 0 {
		maxLeftTime = 30 * time.Minute
	}
	db.DB().SetConnMaxLifetime(maxLeftTime)

	db.LogMode(config.EnableLog)

	instance.gdb = db
}

type dbLogger struct {
}

func formatSQL(sql string) string {
	nSQL := strings.ReplaceAll(sql, "\n", " ")
	nSQL = strings.TrimSpace(nSQL)

	replaceWhiteSpaces := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	nSQL = replaceWhiteSpaces.ReplaceAllString(nSQL, " ")

	return nSQL
}

func (*dbLogger) Print(v ...interface{}) {
	if len(v) >= 6 {
		sql := formatSQL(fmt.Sprint(v[3]))
		log.Println(fmt.Sprintf("%s ( %+v ) [%d rows affected or returned] [%+v], %+v", sql, v[4], v[5], v[2], v[1]))
	}
}
