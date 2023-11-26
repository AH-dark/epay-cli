package conn

import (
	"fmt"
	"sync"

	gormlogrus "github.com/onrik/gorm-logrus"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/AH-dark/epay-cli/internal/log"
)

var (
	globalDB *gorm.DB = nil
	once              = sync.Once{}
)

func initDB(c *cli.Context) (*gorm.DB, error) {
	ctx := c.Context

	var dialect gorm.Dialector
	switch c.String("database.driver") {
	case "mysql", "mariadb":
		dialect = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			c.String("database.user"),
			c.String("database.password"),
			c.String("database.host"),
			c.Int("database.port"),
			c.String("database.name"),
			c.String("database.charset"),
		))
	case "postgres", "postgresql":
		dialect = postgres.Open(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			c.String("database.host"),
			c.Int("database.port"),
			c.String("database.user"),
			c.String("database.password"),
			c.String("database.name"),
			c.String("database.sslmode"),
		))
	default:
		err := fmt.Errorf("unsupported database driver: %s", c.String("database.driver"))
		log.Log(ctx).WithError(err).Error("select database driver failed")
		return nil, err
	}

	db, err := gorm.Open(dialect, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.String("database.prefix"),
			SingularTable: true,
			NameReplacer:  NewReplacer(),
		},
		Logger: gormlogrus.New(),
	})
	if err != nil {
		log.Log(ctx).WithError(err).Error("open database failed")
		return nil, err
	}

	if c.String("database.driver") == "mysql" || c.String("database.driver") == "mariadb" {
		db.Set("gorm:table_options", fmt.Sprintf("DEFAULT CHARSET = %s ENGINE = InnoDB", c.String("database.charset")))
	}

	return db, nil
}

func GetDB(c *cli.Context) *gorm.DB {
	once.Do(func() {
		db, err := initDB(c)
		if err != nil {
			log.Log(c.Context).WithError(err).Fatal("init database failed")
			return
		}

		globalDB = db
	})

	return globalDB.WithContext(c.Context)
}
