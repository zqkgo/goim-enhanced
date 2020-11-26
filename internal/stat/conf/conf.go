package conf

import (
	"flag"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/bilibili/discovery/naming"
	log "github.com/golang/glog"
	xtime "github.com/zqkgo/goim-enhanced/pkg/time"
)

var (
	confPath  string
	region    string
	zone      string
	deployEnv string
	host      string
	// Conf config
	Conf *Config
)

func init() {
	var (
		defHost, _ = os.Hostname()
	)
	flag.StringVar(&confPath, "conf", "stat-example.toml", "default config path")
	flag.StringVar(&region, "region", os.Getenv("REGION"), "avaliable region. or use REGION env variable, value: sh etc.")
	flag.StringVar(&zone, "zone", os.Getenv("ZONE"), "avaliable zone. or use ZONE env variable, value: sh001/sh002 etc.")
	flag.StringVar(&deployEnv, "deploy.env", os.Getenv("DEPLOY_ENV"), "deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.")
	flag.StringVar(&host, "host", defHost, "machine hostname. or use default machine hostname.")
}

// Init init config.
func Init() (err error) {
	Conf = Default()
	_, err = toml.DecodeFile(confPath, &Conf)
	log.Infof("Conf: %v", Conf)
	return
}

// Default new a config with specified defualt value.
func Default() *Config {
	return &Config{
		Env:       &Env{Region: region, Zone: zone, DeployEnv: deployEnv, Host: host},
		Discovery: &naming.Config{Region: region, Zone: zone, Env: deployEnv, Host: host},
	}
}

// Config is stat config.
type Config struct {
	Env       *Env
	Collector *Collector
	Redis     *Redis
	Discovery *naming.Config
}

// Env is env config.
type Env struct {
	Region    string
	Zone      string
	DeployEnv string
	Host      string
}

// Collector is collector config.
type Collector struct {
	Itvl xtime.Duration
}

// Redis is redis config.
type Redis struct {
	Network      string
	Addr         string
	Auth         string
	Active       int
	Idle         int
	DialTimeout  xtime.Duration
	ReadTimeout  xtime.Duration
	WriteTimeout xtime.Duration
	IdleTimeout  xtime.Duration
	Expire       xtime.Duration
}
