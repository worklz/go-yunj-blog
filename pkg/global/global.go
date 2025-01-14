package global

import (
	"github.com/worklz/yunj-blog-go/config"

	"github.com/go-playground/validator/v10"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/olivere/elastic/v7"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

var (
	Config   *config.Config
	Logger   *logrus.Logger
	MySQL    *gorm.DB
	Redis    *redis.Pool
	EsClient *elastic.Client
	Validate *validator.Validate
	Corn     *cron.Cron
)
