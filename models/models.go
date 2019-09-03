package models

import (
	"os"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego/orm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// RegisterDB ...
func RegisterDB() {
	// register model
	orm.RegisterModel(
		new(Event),
		new(EventRewards),
		new(EmailMessage),
	)

	orm.RegisterDriver("postgres", orm.DRPostgres)

	err := godotenv.Load()
	if err != nil {
		logs.Error("error", "Error loading .env file")
	}
	DBHOST := os.Getenv("DBHOST")

	orm.RegisterDataBase("default", "postgres", DBHOST)

}
