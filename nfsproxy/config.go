package nfsproxy

import (
	//"fmt"
	"goini"
	"strconv"
)

// ProxyConfig Type
type ProxyConfig struct {
	Bind         string    
	WaitQueueLen int       
	MaxConn      int       
	Timeout      int       
	FailOver     int       
	Backend      []string  
	Log          LogConfig 
	Stats        string    
}

// LogConfig Type
type LogConfig struct {
	Level string 
	Path  string 
}

func parseConfigFile(filepath string) error {
    var err error
	conf := goini.SetConfig(filepath)
	pConfig.Log.Path = conf.GetValue("log", "save_path")
	pConfig.Log.Level = conf.GetValue("log", "level")
	pConfig.Bind = conf.GetValue("setting", "bind")
	pConfig.WaitQueueLen,err = strconv.Atoi(conf.GetValue("setting", "wait_queue_len"))
	if err != nil{
		return err
	}
	pConfig.MaxConn,err = strconv.Atoi(conf.GetValue("setting", "max_conn"))
	if err != nil{
		return err
	}
	pConfig.Timeout,err = strconv.Atoi(conf.GetValue("setting", "timeout"))
	if err != nil{
		return err
	}
	pConfig.FailOver,err = strconv.Atoi(conf.GetValue("setting", "failover"))
	if err != nil{
		return err
	}
	pConfig.Stats = conf.GetValue("setting", "stats")
	pConfig.Backend = append(pConfig.Backend,conf.GetValue("node", "server"))
	return nil
}
