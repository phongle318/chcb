package db

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fpt-corp/fptshop/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var connection *sqlx.DB

func init() {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		config.Env.DbUsername,
		config.Env.DbPassword,
		config.Env.DbHost,
		config.Env.DbPort,
		config.Env.DbName)
	connection = sqlx.MustConnect("mysql", dataSource)
}

func validateIds(param string) error {
	ids := strings.Split(param, ",")
	for _, id := range ids {
		_, err := strconv.Atoi(id)
		if err != nil {
			return err
		}
	}
	return nil
}
