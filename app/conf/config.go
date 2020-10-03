package conf

import (
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/go-impatient/gaia/pkg/conf"
)

// Config ...
var AppConfig = new(App)
var DBConfig = new(DB)
var LogConfig = new(Log)
var MailConfig = new(Mail)
var CacheConfig = new(Cache)
var Config = new(ConfigTpl)

func InitConfig(name string) (*ConfigTpl, error) {
	// 获取文件后缀
	fileSuffix := path.Ext(name)
	if len(fileSuffix) > 1 {
		// 去掉后缀中.字符
		fileSuffix = fileSuffix[1:]
	}

	if err := conf.LoadConfig(name, fileSuffix); err != nil {
		return nil, err
	}

	// 应用配置初始化
	app := viper.Sub("app")
	if app == nil {
		return nil, errors.New("No found `app` in the configuration")
	}
	AppConfig = NewAppConfig(app)

	// 数据库配置
	db := viper.Sub("database")
	if db == nil {
		return nil, errors.New("No found `database` in the configuration")
	}
	DBConfig = NewDBConfig(db)

	// 日志配置
	log := viper.Sub("log")
	if log == nil {
		return nil, errors.New("No found `log` in the configuration")
	}
	LogConfig = NewLogConfig(log)

	// 邮箱配置
	mail := viper.Sub("mail")
	if mail == nil {
		return nil, errors.New("No found `mail` in the configuration")
	}
	MailConfig = NewMailConfig(mail)

	// 缓存配置
	cache := viper.Sub("cache")
	if cache == nil {
		return nil, errors.New("No found `cache` in the configuration")
	}
	CacheConfig = NewCacheConfig(cache)

	// 全部配置
	Config = &ConfigTpl{
		App:      AppConfig,
		Database: DBConfig,
		Log:      LogConfig,
		Mail:     MailConfig,
		Cache:    CacheConfig,
	}
	return Config, nil
}
