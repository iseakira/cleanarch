package database

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
  InstanceSQLite int = iota
  InstanceMySQL
)

var (
  errInvalidSQLDatabaseInstance = errors.New("invalid sql db instance")
)

func NewDatabaseSQLFactory(instance int) (db *gorm.DB, err error) {
  switch instance {
  case InstanceMySQL:
    configs := NewConfigMySQL()
    dsn := fmt.Sprintf(
      "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
      configs.User,
      configs.Password,
      configs.Host,
      configs.Port,
      configs.Database)
    db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
  case InstanceSQLite:
    configs := NewConfigSQLite()
    db, err = gorm.Open(sqlite.Open(configs.Database), &gorm.Config{})
  default:
    return nil, errInvalidSQLDatabaseInstance
  }
  return db, err
}