package app

import (
	"bufio"
	"fmt"
	"github.com/dthung1602/arc/pkg/resp"
	"log"
	"net"
)

type App struct {
	Config   Config
	stopChan chan bool
	stopped  bool
	parser   resp.Parser
}

func NewApp() (*App, error) {
	app := App{
		Config:   Config{},
		stopChan: make(chan bool),
		stopped:  false,
		parser:   resp.Parser{},
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
	resp, err := app.parser.Parse(r)

	if err != nil {
		panic(err)
	}

	fmt.Printf("GOT COMMAND: %v\n", resp)

	conn.Close()
	log.Println("Closed connection")
}
