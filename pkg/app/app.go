package app

import (
	"bufio"
	"fmt"
	"github.com/dthung1602/arc/pkg/core"
	"github.com/dthung1602/arc/pkg/resp"
	"io"
	"log"
	"net"
)

type App struct {
	Config                Config
	Parser                resp.Parser
	CommandHandlerFactory core.CommandHandlerFactory
	stopChan              chan bool
	stopped               bool
}

func NewApp() (*App, error) {
	app := App{
		Config:                Config{},
		Parser:                resp.Parser{},
		CommandHandlerFactory: core.CommandHandlerFactoryImpl,
		stopChan:              make(chan bool),
		stopped:               false,
	}
	if err := app.Config.Read(); err != nil {
		return nil, err
	}
	return &app, nil
}

func (app *App) Serve() {
	port := app.Config.GetInt("port")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	log.Printf("Arc starts serving on port %d\n", port)

	go func() {
		<-app.stopChan
		if err := listener.Close(); err != nil {
			panic(err)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if app.stopped {
				return
			}
			// TODO do some recovery
			panic(err)
		}
		log.Printf("Accepted new connection from %v\n", conn.RemoteAddr())
		go app.handle(conn)
	}
}

func (app *App) Stop() {
	app.stopped = true
	app.stopChan <- true
}

func (app *App) handle(conn net.Conn) {
	log.Println("Processing ...")

	r := bufio.NewReaderSize(conn, app.Config.GetInt("buffersize"))
	defer func() {
		conn.Close()
		log.Println("Closed connection")
	}()

	for {
		request, err := app.Parser.ParseArray(r)
		if err != nil {
			if err == io.EOF {
				log.Println("Done processing")
				break
			}
			writeError(conn, err)
			continue
		}
		log.Printf("Received: %v\n", request)

		handler, err := app.CommandHandlerFactory(request)
		if err != nil {
			writeError(conn, err)
			continue
		}
		response, err := handler.Handle(request)
		if err != nil {
			writeError(conn, err)
			continue
		}

		writeResponse(conn, response)
	}
}

func writeResponse(conn net.Conn, res resp.Resp) {
	log.Printf("Responsed: %s", res.String())
	_, err := conn.Write(res.Resp())
	if err != nil {
		panic(err)
	}
}

func writeError(conn net.Conn, err error) {
	log.Printf("Error occured: %v\n", err)
	writeResponse(conn, resp.NewSimpleError(err))
}
