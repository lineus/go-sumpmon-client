package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lineus/go-loadaws"
	awsNotify "github.com/lineus/go-notify"
	"github.com/lineus/go-sqlitelogs"
	"github.com/lineus/go-sumpmon"
)

const live bool = false

var config loadaws.Config

func connect(logger sqlitelogs.SqliteLogger) error {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp4", "192.168.1.15:1776")
	if err != nil {
		return err
	}

	defer conn.Close()

	_, err = conn.Write([]byte("ping"))
	if err != nil {
		return err
	}

	res := make([]byte, 4)

	_, err = conn.Read(res)
	if err != nil {
		return err
	}

	_, err = logger.SaveLog("txBeacon", string(res))
	return err
}

func goodPongLast60Mins(logger sqlitelogs.SqliteLogger) bool {
	goodPong := false
	logs, err := logger.GetLogsBetween(time.Now().Add(-60*time.Minute), time.Now())
	if err != nil {
		log.Fatal("Failed To Get Recent Logs: ", err)
	}

	for _, v := range logs {
		if v.Action == "txBeacon" && v.Result == "pong" {
			goodPong = true
		}
	}

	return goodPong
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	runs := 0
	logger, err := sumpmon.Init("/tmp/test.db")
	if err != nil {
		log.Fatal("Failed To Initialize Logger: ", err)
	}

	config, err = loadaws.FromJSON()
	if err != nil {
		log.Fatal("loadAWSEnv Fail: ", err)
	}

	if live {
		defer awsNotify.Send(config.AWS.SNS.ARN, "SumpMonitor Off")
	}

	go func() {
		for {
			if !logger.Alive() && runs > 5 {
				notify("Sump Logger Not Alive", true)
			}

			if !goodPongLast60Mins(logger) && runs > 5 {
				notify("No Good Pongs For 1 Hour", true)
			}
			err := connect(logger)
			if err != nil {
				notify("Cant Connect To Sump Monitor", true)
			}
			runs++
			time.Sleep(10 * time.Second)
		}
	}()
	<-sigs
	notify("You Haved Murdered My Monitoring!", true)

}

func notify(s string, e bool) {
	if live {
		awsNotify.Send(config.AWS.SNS.ARN, s)
	} else {
		fmt.Println(s)
	}
	if e {
		os.Exit(1)
	}
}
