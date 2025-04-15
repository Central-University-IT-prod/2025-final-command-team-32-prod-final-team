package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"solution/pkg/logger"
	"sync"
	"time"

	"solution/internal/app"
)

//	@title		MetaCinema
//	@version	1.0
//	@host		https://prod-team-32-n26k57br.REDACTED
//	@BasePath	/api/v1

// @securitydefinitions.apikey	Bearer
// @in							header
// @name						Authorization
// @descrtiption				"access token 'Bearer {token}'"
func main() {
	var (
		shutDownGroup sync.WaitGroup
		ctx           = context.Background()
		signalCh      = make(chan os.Signal, 1)
	)
	signal.Notify(signalCh, os.Interrupt)

	ctx = logger.CtxWithLogger(ctx)
	application := app.NewApp(ctx)
	// graceful shutdown handler
	go InterruptHandler(application, signalCh, &shutDownGroup)
	application.Start()
	// connections closing on shutdown
	application.CloseConnections(ctx)
}

func InterruptHandler(app *app.App, signalCh chan os.Signal, group *sync.WaitGroup) {
	<-signalCh
	fmt.Printf("\n**GRACEFULLY SHUTTING DOWN**\n\n")
	group.Add(1)
	defer group.Done()
	app.Fiber.ShutdownWithTimeout(15 * time.Second)
}
