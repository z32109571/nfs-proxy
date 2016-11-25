package nfsproxy

import (
	"fmt"
	"log"
	//"nfs/rpc"
	"logrus"
	"flag"
	"os"
	"os/signal"
	"syscall"
)
var (
	pConfig    ProxyConfig
	pLog       *logrus.Logger
	configFile = flag.String("c", "conf/conf.ini", "config file,default conf/conf.ini")
)

func onExitSignal() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
	
L:
	for {
		sig := <-sigChan
		switch sig {
		case syscall.SIGINT:
			log.Fatal("Reopen log file")
		case syscall.SIGTERM:
			log.Fatal("Catch SIGTERM singal, exit.")
			break L
		}
	}
}

func main() {
	flag.Parse()
	fmt.Println("Start Proxy...")

	if parseConfigFile(*configFile) != nil {
		return
	}
	// init logger server
	initLogger()
	go onExitSignal()

	// init status service
	//initStats()
	// init proxy service
	initProxy()
}
