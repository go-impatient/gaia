package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type Config struct {
	Name   string // 配置文件名
	Format string // 配置文件类型
}

// ReadConfig 导入配置文件并解析
func (c *Config) ReadConfig() error {
	if c.Name != "" { // 如果指定了配置文件, 解析指定的配置文件
		viper.SetConfigFile(c.Name)
	} else { // 否则, 没有指定配置文件, 则解析默认的配置文件
		// 找到项目路径
		path, _ := filepath.Abs("./../../")
		viper.AddConfigPath(filepath.Join(path, "config"))

		// 或者从部署项目的主目录中查找
		viper.AddConfigPath("/etc/gaia/")
		viper.AddConfigPath("$HOME/.gaia/")

		viper.SetConfigName("config")
	}

	// Config's format: "json" | "toml" | "yaml" | "yml"
	viper.SetConfigType(c.Format)

	// 读取匹配的环境变量
	viper.AutomaticEnv()

	// 读取环境变量的前缀为Gaia, 会自动大写(GAIA)
	// 例如: 当前设置的前缀为 Gaia. 定义一个环境变量名为 GAIA_LOG_PATH, 会自动转换为 log.path, 就可以使用 viper.GetString("log.path")
	viper.SetEnvPrefix("Gaia")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// viper 解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Using config file: %s [%s]\n", viper.ConfigFileUsed(), err)
		return err
	}

	return nil
}

// LoadConfig 加载配置文件
func LoadConfig(name string, format string) error {
	c := &Config{
		Name:   name,
		Format: format,
	}

	// 初始化配置文件
	if err := c.ReadConfig(); err != nil {
		return err
	}

	// 调用监控配置文件变化并热加载程序
	c.WatchConfig()

	return nil
}

// WatchConfig 监控配置文件变化并热加载程序
func (c *Config) WatchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s \n", e.Name)
	})
}
