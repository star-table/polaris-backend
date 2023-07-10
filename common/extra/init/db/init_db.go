package db

import (
	"fmt"
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/polaris-backend/common/core/consts"
	"runtime"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var (
	_sourcePath         = "file://./init/db/migrations/"
	_sourceLocalPath    = "file://../../../init/db/migrations/"
	_sourceAppLocalPath = "file://../init/db/migrations/"
	log                 = logger.GetDefaultLogger()
)

type MigrationsLog struct {
	Log *logger.SysLogger
}

func (s MigrationsLog) Printf(format string, v ...interface{}) {
	s.Log.Infof(format, v...)
}

func (s MigrationsLog) Verbose() bool {
	return true
}

func DbMigrations(env string, dbHost string, port int, user string, password string, dbName string) errs.SystemErrorInfo {

	sqlFile := _sourcePath
	if runtime.GOOS != consts.LinuxGOOS {
		sqlFile = _sourceLocalPath
	}
	if env == "test" {
		sqlFile = _sourceLocalPath
	}
	applicationName := config.GetApplication().Name

	if applicationName == consts.AppApplicationName {
		sqlFile = _sourceAppLocalPath
	}

	sqlFile = sqlFile + applicationName

	jstr, _ := json.ToJson(config.GetApplication())
	log.Info(jstr)
	log.Info("init db env: " + env + "; filePath: " + sqlFile)

	mysqlstr := fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s", user, password, dbHost, port, dbName)
	log.Info("init db mysqlstr: " + mysqlstr)
	m, err := migrate.New(
		sqlFile,
		mysqlstr)

	if err != nil {
		log.Error(" init db migrations fail. " + strs.ObjectToString(err))
		return errs.BuildSystemErrorInfo(errs.InitDbFail)
	}

	migLog := MigrationsLog{Log: log}
	m.Log = migLog

	fmt.Println(m.Version())
	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		log.Error(" init db migrations up fail. " + strs.ObjectToString(err))
		return errs.BuildSystemErrorInfo(errs.InitDbFail)
	}
	return nil
}
