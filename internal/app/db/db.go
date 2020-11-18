package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/zliang90/kingRest/internal/app/conf"
	"github.com/zliang90/kingRest/internal/app/db/callbacks"
	"github.com/zliang90/kingRest/pkg/log"
)

var (
	// default db
	_db *gorm.DB

	// multiple mysql data source
	dbs = make(map[string]*gorm.DB)
)

// Init init default db
func Init(c *conf.Config) error {
	if _db != nil {
		return nil
	}
	if c == nil {
		return fmt.Errorf("config is nil")
	}
	return initDataSources(c)
}

// Close close all database connections
func Close() (err error) {
	defer func() {
		if e := recover(); e != nil {
			if er, ok := e.(error); ok {
				err = er
				return
			}
			err = fmt.Errorf("%v", e)
		}
	}()
	// close connection
	for _, _db := range dbs {
		if _db != nil {
			_db.Close()
		}
	}
	return
}

func initDataSources(c *conf.Config) error {
	var err error
	for k, db := range c.DataSources {
		dbs[k], err = gorm.Open("mysql", db.Addr)
		if err != nil {
			return err
		}
		// testing connection
		err = dbs[k].DB().Ping()
		if err != nil {
			return err
		}
		// setting db connection
		dbs[k].SingularTable(true)
		dbs[k].DB().SetMaxIdleConns(db.IdleConn)
		dbs[k].DB().SetMaxOpenConns(db.MaxConn)
		// db logger
		if db.Debug {
			dbs[k].LogMode(true)
			dbs[k].SetLogger(log.GetLogger())
		}
	}

	var ok bool
	if _db, ok = dbs["default"]; !ok {
		return fmt.Errorf("the 'default' datasource is not defined")
	}
	// Auto migrate
	log.Info("auto db migration")
	// defaultDB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8")
	_db.AutoMigrate(
		&User{},
	)
	// Replace delete callback
	_db.Callback().Delete().Replace("gorm:delete", callbacks.DeleteCallback)

	// data initialize
	if c.Env != "prod" {
		dataInitialize()
	}

	return nil
}

// GetDefaultDB get the default db object
func GetDefaultDB() *gorm.DB {
	return _db
}

// GetDbWithName get db object with name
func GetDbWithName(name string) (*gorm.DB, error) {
	if db, ok := dbs[name]; ok {
		return db, nil
	}
	return nil, fmt.Errorf("the '%s' database datasource not defined", name)
}
