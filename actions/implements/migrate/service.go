package migrate

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"gorm.io/gorm"

	"github.com/AH-dark/epay-cli/actions/factory"
	"github.com/AH-dark/epay-cli/database/conn"
	"github.com/AH-dark/epay-cli/database/model"
	"github.com/AH-dark/epay-cli/internal/log"
	"github.com/AH-dark/epay-cli/pkg/utils"
)

type service struct {
}

func NewService() factory.ActionService {
	return &service{}
}

func (*service) getDefaultDatabaseConfig(c *cli.Context) (map[string]string, error) {
	var data map[string]string

	logger := log.Log(c.Context)

	path := utils.GetAbsolutePath(c.String("migrate.default_config"))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.WithError(err).Error("default config file not exists")
		return data, nil
	}

	logger = logger.WithField("path", path)

	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		logger.WithError(err).Error("open default config file failed")
		return data, err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			logger.WithError(err).Warn("close default config file failed")
		}
	}(f)

	if err := json.NewDecoder(f).Decode(&data); err != nil {
		logger.WithError(err).Error("decode default config file failed")
		return data, err
	}

	return data, nil
}

func (*service) getDefaultPaymentTypes(c *cli.Context) ([]model.Type, error) {
	var data []model.Type

	logger := log.Log(c.Context)

	path := utils.GetAbsolutePath(c.String("migrate.default_payment_types"))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		logger.WithError(err).Error("default payment types file not exists")
		return data, nil
	}

	logger = logger.WithField("path", path)

	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		logger.WithError(err).Error("open default payment types file failed")
		return data, err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			logger.WithError(err).Warn("close default payment types file failed")
		}
	}(f)

	if err := json.NewDecoder(f).Decode(&data); err != nil {
		logger.WithError(err).Error("decode default payment types file failed")
		return data, err
	}

	return data, nil
}

func (*service) getDefaultGroup(_ *cli.Context) *model.Group {
	return &model.Group{
		GID:  0,
		Name: "默认分组",
		Info: `{"1":{"type":"","channel":"-1","rate":""},"2":{"type":"","channel":"-1","rate":""},"3":{"type":"","channel":"-1","rate":""}}`,
	}
}

func (*service) createDatabaseConfig(c *cli.Context, configs map[string]string) error {
	db := conn.GetDB(c)
	ctx := c.Context

	for k, v := range configs {
		logger := log.Log(ctx).WithField("key", k).WithField("val", v)
		logger.Debug("create default config")

		// check if config exists
		if err := db.Model(&model.Config{}).
			Where("k = ?", k).
			First(&model.Config{}).
			Error; err == nil {
			logger.Debug("config already exists")
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WithError(err).Error("check default config failed")
			return err
		}

		// create config
		if err := db.
			Model(&model.Config{}).
			Create(&model.Config{
				Key: k,
				Val: v,
			}).Error; err != nil {
			logger.WithError(err).Error("create default config failed")
			return err
		}
	}

	return nil
}

func (svc *service) Do(c *cli.Context) error {
	ctx := c.Context
	db := conn.GetDB(c)

	log.Log(ctx).Info("auto migrate database tables")
	if err := db.AutoMigrate(
		&model.AlipayRisk{},
		&model.Anounce{},
		&model.Batch{},
		&model.Channel{},
		&model.Config{},
		&model.Domain{},
		&model.Group{},
		&model.Log{},
		&model.Order{},
		&model.Plugin{},
		&model.Record{},
		&model.Regcode{},
		&model.Risk{},
		&model.Roll{},
		&model.Settle{},
		&model.Type{},
		&model.User{},
		&model.Weixin{},
	); err != nil {
		log.Log(ctx).WithError(err).Error("auto migrate database tables failed")
		return err
	}

	// alter auto increment if user table is empty
	if err := db.Model(&model.User{}).First(&model.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		log.Log(ctx).Info("alter auto increment")
		if err := db.Model(&model.User{}).
			Exec(fmt.Sprintf(
				"alter table `%s` AUTO_INCREMENT = 1000",
				db.NamingStrategy.TableName("user"),
			)).Error; err != nil {
			log.Log(ctx).WithError(err).Error("alter auto increment failed")
			return err
		}
	}

	log.Log(ctx).Info("get default config")
	defaultDatabaseConfig, err := svc.getDefaultDatabaseConfig(c)
	if err != nil {
		log.Log(ctx).WithError(err).Error("get default config failed")
		return err
	}

	log.Log(ctx).Info("create default config")
	if err := svc.createDatabaseConfig(c, defaultDatabaseConfig); err != nil {
		log.Log(ctx).WithError(err).Panic("create default config failed")
		return err
	}

	log.Log(ctx).Info("init app config")
	initData := map[string]string{
		"syskey":  c.String("app.syskey"),
		"build":   time.Now().Format("2006-01-02"),
		"cronkey": c.String("app.cronkey"),
	}
	if err := svc.createDatabaseConfig(c, initData); err != nil {
		log.Log(ctx).WithError(err).Panic("create app init config failed")
		return err
	}

	log.Log(ctx).Info("create default payment types")
	defaultPaymentTypes, err := svc.getDefaultPaymentTypes(c)
	if err != nil {
		log.Log(ctx).WithError(err).Error("get default payment types failed")
		return err
	}

	for _, t := range defaultPaymentTypes {
		if err := db.Model(&model.Type{}).
			Where("name = ?", t.Name).
			First(&model.Type{}).Error; err == nil {
			continue
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Log(ctx).WithError(err).Error("check default payment type failed")
			return err
		}

		if err := db.Model(&model.Type{}).Create(&t).Error; err != nil {
			log.Log(ctx).WithError(err).Error("create default payment type failed")
			return err
		}
	}

	log.Log(ctx).Info("create default group")
	defaultGroup := svc.getDefaultGroup(c)
	if err := db.Model(&model.Group{}).
		Where("gid = ?", defaultGroup.GID).
		First(&model.Group{}).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Log(ctx).WithError(err).Error("check default group failed")
		return err
	} else if err := db.Model(&model.Group{}).Create(&defaultGroup).Error; err != nil {
		log.Log(ctx).WithError(err).Error("create default group failed")
		return err
	}

	fmt.Printf(`
System Key: %s
Build Time: %s
Cron Key: %s

Admin Username: admin
Admin Password: 123456
`, initData["syskey"], initData["build"], initData["cronkey"])

	return nil
}
