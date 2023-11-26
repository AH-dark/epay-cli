package log

import (
	"context"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var once = sync.Once{}

func InitLogger(_ context.Context, debug bool) (*logrus.Logger, error) {
	logger := logrus.StandardLogger()
	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.TextFormatter{})

	if debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger, nil
}

func GetLogger(ctx context.Context) *logrus.Logger {
	once.Do(func() {
		_, err := InitLogger(ctx, false)
		if err != nil {
			logrus.WithContext(ctx).WithError(err).Fatal("failed to init logger")
			return
		}

		logrus.WithContext(ctx).Info("logger initialized")
	})

	return logrus.StandardLogger()
}

func Log(ctx context.Context) *logrus.Entry {
	return GetLogger(ctx).WithContext(ctx)
}
