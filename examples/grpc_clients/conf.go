package main

import (
	"flag"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/bilibili/discovery/naming"
)

type Config struct {
	Env           *Env
	DiscoveryConf *naming.Config
}

type Env struct {
	Region    string
	Zone      string
	DeployEnv string
	Host      string
}

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
	flag.StringVar(&confPath, "conf", "clients-example.toml", "default config path")
	flag.StringVar(&region, "region", os.Getenv("REGION"), "avaliable region. or use REGION env variable, value: sh etc.")
	flag.StringVar(&zone, "zone", os.Getenv("ZONE"), "avaliable zone. or use ZONE env variable, value: sh001/sh002 etc.")
	flag.StringVar(&deployEnv, "deploy.env", os.Getenv("DEPLOY_ENV"), "deploy env. or use DEPLOY_ENV env variable, value: dev/fat1/uat/pre/prod etc.")
	flag.StringVar(&host, "host", defHost, "machine hostname. or use default machine hostname.")
}

func Init() (err error) {
	Conf = &Config{
		Env:           &Env{Region: region, Zone: zone, DeployEnv: deployEnv, Host: host},
		DiscoveryConf: &naming.Config{Region: region, Zone: zone, Env: deployEnv, Host: host},
	}
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
