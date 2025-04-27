package main

import (
	"flag"
	"fmt"
	"github.com/gin-contrib/pprof"
	"os"
	"os/signal"
	"personal-secretary-user-ap/internal/api"
	"personal-secretary-user-ap/internal/application"
	"strconv"
	"syscall"
)

func main() {
	envFile := flag.String("env_file", ".env", "Path to environment file")
	flag.Parse()

	app, err := application.GetInstance(envFile)
	if nil != err {
		msg := fmt.Sprintf("Failed to initialize application: %s", err.Error())
		if nil != app && nil != app.GetLogger() {
			app.GetLogger().Emergency(msg)
		} else {
			_, err := os.Stderr.WriteString(msg + "\n")
			if err != nil {
				return
			}
		}
		os.Exit(1)
	}
	defer app.Close()

	router := api.GetRouter()
	pprof.Register(router)

	go func() {
		addr := "0.0.0.0:" + strconv.Itoa(app.GetConfig().Api.Port)
		err = router.Run(addr)
		if nil != err {
			// TODO: log
			//loggerService.GetLogger().Fatal(logTag, "Can't run server: "+err.Error())
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
}
