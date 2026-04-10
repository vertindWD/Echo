package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局配置变量
var Conf = new(AppConfig)

// AppConfig 整个配置的聚合根
type AppConfig struct {
	*App       `mapstructure:"app"`
	*Log       `mapstructure:"log"`
	*MySQL     `mapstructure:"mysql"`
	*Redis     `mapstructure:"redis"`
	*Snowflake `mapstructure:"snowflake"`
	*Auth      `mapstructure:"auth"`
}

type App struct {
	Name string `mapstructure:"name"`
	Mode string `mapstructure:"mode"`
	Port int    `mapstructure:"port"`
}

type Log struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQL struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxLifetime  int    `mapstructure:"max_lifetime"`
}

type Redis struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type Snowflake struct {
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
}

type Auth struct {
	JwtSecret string `mapstructure:"jwt_secret"`
	JwtExpire int    `mapstructure:"jwt_expire"`
}

func Init() (err error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 将读取的配置反序列化到全局变量 Conf 中
	if err := viper.Unmarshal(Conf); err != nil {
		return fmt.Errorf("配置文件反序列化失败: %w", err)
	}

	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("检测到配置文件修改，正在热重载...")
		// 每次修改后，必须重新反序列化到 Conf 中
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("热重载配置失败: %v\n", err)
		}
	})

	return nil
}
