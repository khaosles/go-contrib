package addr

import (
	"flag"
	"fmt"

	"github.com/khaosles/go-contrib/core/config"
)

/*
   @File: addr.go
   @Author: khaosles
   @Time: 2023/8/17 00:44
   @Desc:
*/

func GetAddr(serverName string) string {
	var host string
	var port string
	flag.StringVar(&host, "host", "", "host")
	flag.StringVar(&port, "port", "", "port")
	flag.Parse()
	if port == "" {
		port = config.Viper.GetString(serverName + ".port")
	}
	if host == "" {
		host = config.Viper.GetString(serverName + ".host")
	}
	return fmt.Sprintf("%s:%s", host, port)
}
