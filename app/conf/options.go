package conf

import (
	"time"

	"github.com/spf13/viper"
)

// App ...
type App struct {
	Mode         string        `json:"mode" yaml:"mode"`   // "dev" | "prod" | "test"
	Grace        bool          `json:"grace" yaml:"grace"` // 默认 false, 可选 pprof 为空不使用
	Host         string        `json:"host" yaml:"host"`
	Port         int           `json:"port" yaml:"port"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
	MaxPingCount int           `json:"max_ping_count" yaml:"max_ping_count"` // pingServer函数try的次数
	JWTSecret    string        `json:"jwt_secret" yaml:"jwt_secret"`
	TLS          *TLS          `json:"tls" yaml:"tls"`
	AutoTLS      *AutoTLS      `json:"auto_tls" yaml:"auto_tls"`
}

// TLS ...
type TLS struct {
	Enabled         bool   `json:"enabled" yaml:"enabled"`                     //  是否开启https
	Port            int    `json:"port" yaml:"port"`                           // https 端口
	CertPath        string `json:"cert_path" yaml:"cert_path"`                 // the cert file (leave empty when using letsencrypt)
	KeyPath         string `json:"key_path" yaml:"key_path"`                   // the cert key (leave empty when using letsencrypt)
}

// AutoTLS ...
type AutoTLS struct {
	Enabled   bool     `json:"enabled" yaml:"enabled"`       // if the certificate should be requested from letsencrypt
	AcceptTos bool     `json:"accept_tos" yaml:"accept_tos"` // if you accept the tos from letsencrypt
	Folder    string   `json:"folder" yaml:"folder"`         // the directory of the cache from letsencrypt
	Hosts     []string `json:"hosts" yaml:"hosts"`           // the hosts for which letsencrypt should request certificates
}

// DB ...
type DB struct {
	Dialect      string `json:"dialect" yaml:"dialect"`
	Host         string `json:"host" yaml:"host"`
	Port         int    `json:"port" yaml:"port"`
	Username     string `json:"username" yaml:"username"`
	Password     string `json:"password" yaml:"password"`
	Database     string `json:"database" yaml:"database"`
	Ssl          bool   `json:"ssl" yaml:"ssl"`
	MaxIdleConns int    `json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns" yaml:"max_open_conns"`
	Logging      bool   `json:"logging" yaml:"logging"`
}

// Log ...
type Log struct {
	Path   string `json:"path" yaml:"path"`
	Level  string `json:"level" yaml:"level"`   // "trace" | "debug" | "info" | "warn" | "error" | "fatal" | "panic" | ""
	Format string `json:"format" yaml:"format"` // "pretty" | "json"
}

// Mail ...
type Mail struct {
	Enabled    bool   `json:"enabled" yaml:"enabled"`
	Username   string `json:"username" yaml:"username"`
	Password   string `json:"password" yaml:"password"`
	MailServer string `json:"mail_server" yaml:"mail_server"`
	Port       int    `json:"port" yaml:"port"`
	From       string `json:"from" yaml:"from"`
}

// Cache ...
type Cache struct {
	Type    string        `json:"type" yaml:"type"`
	Timeout time.Duration `json:"timeout" yaml:"timeout"`
	Redis   *CacheRedis   `json:"redis" yaml:"redis"`
}

// CacheRedis ...
type CacheRedis struct {
	Host      string `json:"host" yaml:"host"`
	Port      int    `json:"port" yaml:"port"`
	Password  string `json:"password" yaml:"password"`
	DB        int    `json:"db" yaml:"db"`
	KeyPrefix string `json:"key_prefix" yaml:"key_prefix"`
}

// ConfigTpl ...
type ConfigTpl struct {
	App      *App   `json:"app" yaml:"app"`
	Database *DB    `json:"database" yaml:"database"`
	Log      *Log   `json:"log" yaml:"log"`
	Mail     *Mail  `json:"mail" yaml:"mail"`
	Cache    *Cache `json:"cache" yaml:"cache"`
}

// NewAppConfig ...
func NewAppConfig(cfg *viper.Viper) *App {
	return &App{
		Grace:        cfg.GetBool("grace"),
		Mode:         cfg.GetString("mode"),
		Host:         cfg.GetString("host"),
		Port:         cfg.GetInt("port"),
		ReadTimeout:  cfg.GetDuration("read_timeout"),
		WriteTimeout: cfg.GetDuration("write_timeout"),
		IdleTimeout:  cfg.GetDuration("idle_timeout"),
		TLS: &TLS{
			Enabled: cfg.GetBool("tls.enabled"),
			Port: cfg.GetInt("tls.port"),
			CertPath: cfg.GetString("tls.cert_path"),
			KeyPath:  cfg.GetString("tls.key_path"),
		},
		AutoTLS: &AutoTLS{
			Enabled: cfg.GetBool("auto_tls.enabled"),
			AcceptTos: cfg.GetBool("auto_tls.auto_tls"),
			Folder:  cfg.GetString("auto_tls.folder"),
			Hosts:    cfg.GetStringSlice("auto_tls.hosts"),
		},
	}
}

// NewAppConfig ...
func NewDBConfig(cfg *viper.Viper) *DB {
	return &DB{
		Dialect:      cfg.GetString("dialect"),
		Host:         cfg.GetString("host"),
		Port:         cfg.GetInt("port"),
		Username:     cfg.GetString("username"),
		Password:     cfg.GetString("password"),
		Database:     cfg.GetString("database"),
		Ssl:          cfg.GetBool("ssl"),
		MaxIdleConns: cfg.GetInt("max_idle_conns"),
		MaxOpenConns: cfg.GetInt("max_open_conns"),
		Logging:      cfg.GetBool("logging"),
	}
}

// NewLogConfig ...
func NewLogConfig(cfg *viper.Viper) *Log {
	return &Log{
		Path:   cfg.GetString("path"),
		Level:  cfg.GetString("level"),
		Format: cfg.GetString("format"),
	}
}

// NewMailConfig ...
func NewMailConfig(cfg *viper.Viper) *Mail {
	return &Mail{
		Enabled:    cfg.GetBool("enabled"),
		Username:   cfg.GetString("username"),
		Password:   cfg.GetString("password"),
		MailServer: cfg.GetString("mail_server"),
		Port:       cfg.GetInt("port"),
		From:       cfg.GetString("from"),
	}
}

// NewCacheConfig ...
func NewCacheConfig(cfg *viper.Viper) *Cache {
	return &Cache{
		Type:    cfg.GetString("type"),
		Timeout: cfg.GetDuration("timeout"),
		Redis: &CacheRedis{
			Host:      cfg.GetString("host"),
			Port:      cfg.GetInt("port"),
			Password:  cfg.GetString("password"),
			DB:        cfg.GetInt("db"),
			KeyPrefix: cfg.GetString("key_prefix"),
		},
	}
}
