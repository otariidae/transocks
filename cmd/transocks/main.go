// transocks server.
package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/cybozu-go/log"
	"github.com/cybozu-go/transocks"
)

type tomlConfig struct {
	Listen   string
	ProxyURL string `toml:"proxy_url"`
	LogLevel string `toml:"log_level"`
	LogFile  string `toml:"log_file"`
}

var (
	configFile = flag.String("f", "/usr/local/etc/transocks.toml",
		"TOML configuration file path")
)

func loadConfig() (*transocks.Config, string, error) {
	tc := new(tomlConfig)
	md, err := toml.DecodeFile(*configFile, tc)
	if err != nil {
		return nil, "", err
	}
	if len(md.Undecoded()) > 0 {
		return nil, "", fmt.Errorf("undecoded key in TOML: %v", md.Undecoded())
	}

	c := transocks.NewConfig()
	c.Listen = tc.Listen

	u, err := url.Parse(tc.ProxyURL)
	if err != nil {
		return nil, "", err
	}
	c.ProxyURL = u

	if err = log.DefaultLogger().SetThresholdByName(tc.LogLevel); err != nil {
		return nil, "", err
	}

	return c, tc.LogFile, nil
}

func main() {
	flag.Parse()

	c, logfile, err := loadConfig()
	if err != nil {
		log.ErrorExit(err)
	}

	if len(logfile) > 0 {
		f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.ErrorExit(err)
		}
		defer f.Close()
		log.DefaultLogger().SetOutput(f)
	}

	srv, err := transocks.NewServer(c)
	if err != nil {
		log.ErrorExit(err)
	}
	log.Info("server starts", nil)

	srv.Serve()

	log.Info("server ends", nil)
}
