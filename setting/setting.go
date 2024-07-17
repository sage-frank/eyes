package setting

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type MConfig struct {
	Dev     *IPath `mapstructure:"dev"`
	Release *IPath `mapstructure:"release"`
	Prod    *IPath `mapstructure:"prod"`
}

type IPath struct {
	Path string `mapstructure:"path"`
}

type AutoConf struct {
	Host      string `mapstructure:"host"`
	Path      string `mapstructure:"path"`
	SecretKey string `mapstructure:"key"`
}

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
	*Jwt         `mapstructure:"jwt"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	ShowSQL      bool   `mapstructure:"show_sql"`
}

type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Password     string `mapstructure:"password"`
	Port         int    `mapstructure:"port"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type LogConfig struct {
	Level       string `mapstructure:"level"`
	Filename    string `mapstructure:"filename"`
	ErrFilename string `mapstructure:"err_filename"`
	MaxSize     int    `mapstructure:"max_size"`
	MaxAge      int    `mapstructure:"max_age"`
	MaxBackups  int    `mapstructure:"max_backups"`
}

type Jwt struct {
	PublicKey string `mapstructure:"public_key"`
	AppID     string `mapstructure:"app_id"`
	AppKey    string `mapstructure:"app_key"`
	AppSecret string `mapstructure:"app_secret"`
	Portal    string `mapstructure:"portal"`
	BaseURL   string `mapstructure:"base_url"`
}

func LoadRemoteConf(autoconf *AutoConf) (*AppConfig, error) {
	var (
		conf AppConfig
		err  error
	)
	fmt.Println(autoconf)
	if err = viper.AddSecureRemoteProvider("etcd3", autoconf.Host, autoconf.Path, autoconf.SecretKey); err != nil {
		return nil, fmt.Errorf("viper.AddSecureRemoteProvide: %w", err)
	}
	viper.SetConfigType("yaml")
	if err = viper.ReadRemoteConfig(); err != nil {
		return nil, fmt.Errorf("viper.ReadRemoteConfig: %w", err)
	}

	if err = viper.Unmarshal(&conf); err != nil {
		return nil, fmt.Errorf(" viper.Unmarshal(&conf): %w", err)
	}

	return &conf, nil
}

func AutoLoad(env string) (*AutoConf, error) {
	autoConf := new(AutoConf)
	executable, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("获取可执行文件路径失败: %v", err)
	}
	configPath := path.Dir(executable)
	if env == "dev" {
		configPath = "."
	}
	c := viper.New()
	c.AddConfigPath(path.Join(configPath, "config"))
	c.SetConfigName("autoload")
	c.SetConfigType("yaml")

	err = c.ReadInConfig()
	if err != nil {
		return nil, err
	}
	autoConf.Host = c.GetString(fmt.Sprintf("%s.host", env))
	autoConf.Path = c.GetString(fmt.Sprintf("%s.path", env))
	autoConf.SecretKey = c.GetString(fmt.Sprintf("%s.key", env))
	return autoConf, nil
}

func LoadConf(env string) (conf *AppConfig, err error) {
	autoConf, err := AutoLoad(env)
	if err != nil {
		return nil, fmt.Errorf("读取autoload配置文件失败: %w", err)
	}

	conf, err = LoadRemoteConf(autoConf)
	if err != nil {
		fmt.Println("从远程配置中心读取配置文件失败:", err.Error())
		return nil, fmt.Errorf("从远程配置中心读取配置文件失败:%w", err)
	}

	return conf, nil
}
