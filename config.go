package main

import (
	"gopkg.in/ini.v1"
	"log"
)

var (
	cfgfile string = "config.ini"
	config  *ini.File
)

func GetConfig(configFile string) (*ini.File, error) {
	var Cfg *ini.File
	Cfg, err := ini.Load(configFile)
	if err != nil {
		return Cfg, err
	}
	return Cfg, err
}

func init() {

	var err error
	config, err = GetConfig(cfgfile)
	if err != nil {
		log.Fatalf("error loading config: %s\n", err)
	}

}
