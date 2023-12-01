package main

import (
	"flag"
	"os"
	"os/signal"
	"strings"
	"syscall"

	dnsproxy "github.com/codingeasygo/dns-proxy"
)

var upper string
var listen string
var cache string
var help bool

func init() {
	flag.StringVar(&upper, "upper", "8.8.4.4:53", "upper dns server")
	flag.StringVar(&listen, "listen", ":53", "listen address")
	flag.StringVar(&cache, "cache", "dns.json", "dns cache store file")
	flag.BoolVar(&help, "h", false, "show help")
	flag.Parse()
}

func main() {
	if help {
		flag.PrintDefaults()
		return
	}
	server := dnsproxy.NewServer()
	server.UpperAddr["*"] = strings.Split(upper, ",")
	server.Listen = listen
	if len(cache) > 0 {
		server.Cache = dnsproxy.NewCache()
		server.Cache.SaveFile = cache
	}
	err := server.Start()
	if err != nil {
		panic(err)
	}

	if server.Cache != nil {
		server.Cache.Start()
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-sigc

	if server.Cache != nil {
		server.Cache.Stop()
	}
	server.Stop()
}
