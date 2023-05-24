package initializers

import (
	"fmt"

	"github.com/MikhailR1337/task-sync-x/domain/migrate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgDb struct {
	*gorm.DB
}

var DB PgDb

func InitDb() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		Cfg.PgHost,
		Cfg.PgUser,
		Cfg.PgPassword,
		Cfg.PgDb,
		Cfg.PgPort,
	)
	// dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	err = migrate.Migrate(db)
	if err != nil {
		return err
	}
	DB = PgDb{db}
	return nil
}
