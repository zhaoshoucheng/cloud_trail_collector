package main

import (
	"cloud_trail_collector/collector"
	"cloud_trail_collector/config"
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
)

var cfgfile = flag.String("c", "/Users/zhaoshoucheng/data/src/mysrc/CloudTrailCollector/config.toml",
	"configuration file, default to config.toml")

func InitSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP,
		syscall.SIGUSR1, syscall.SIGUSR2)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT, syscall.SIGHUP:
			os.Exit(1)
			return
		case syscall.SIGUSR2:
		case syscall.SIGUSR1:
			//return
		default:
			//return
		}
	}
}

func main() {
	flag.Parse()

	// config file parse
	_, err := config.NewConfig(*cfgfile)
	if err != nil {
		panic(err)
	}
	err = collector.MakePipeLines()
	if err != nil {
		panic(err)
	}
	err = collector.StartAllPipelines(context.Background())
	if err != nil {
		panic(err)
	}

	InitSignal()
}
