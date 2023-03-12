package main

import (
	"fmt"
	"strings"

	// Uncomment this block to pass the first stage
	"net"
)

type Application struct {
	storage         *Storage
	closeConnection *chan net.Conn
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// connections := map[string]net.Conn{}
	// var rm sync.RWMutex
	fmt.Println("Logs from your program will appear here!")

	newConnection := make(chan net.Conn)
	closeConnection := make(chan net.Conn)

	app := Application{
		storage:         NewStorage(),
		closeConnection: &closeConnection,
	}

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	nilChecker(err)

	defer l.Close()

	go func() {
		for {
			conn, err := l.Accept()
			nilChecker(err)
			newConnection <- conn
		}
	}()

	go func() {
		conn := <-closeConnection
		fmt.Printf("Disconnected: %v", conn.RemoteAddr().String())
	}()

	for {
		select {
		case conn := <-newConnection:
			go app.handleConnection(conn)
		}
	}
}

func (app *Application) handleConnection(conn net.Conn) {
	defer conn.Close()
	// defer wg.Done()

	fmt.Println(conn.RemoteAddr())

	for {
		buf := make([]byte, 1024)
		data, err := conn.Read(buf)
		if err != nil {
			sendError(conn, "command not found")
			*app.closeConnection <- conn
			return
		}
		if data == 0 {
			// Connection closed by client.
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			break
		}
		if commands, err := parser(string(buf[:data])); err == nil {
			if command, err := commandParser(commands); err == nil {
				app.handleCommand(conn, command)
			} else {
				sendError(conn, err.Error())
			}
		} else {
			sendError(conn, err.Error())
		}
	}
}

func (app *Application) handleCommand(conn net.Conn, command *Command) {

	switch strings.ToUpper(command.cmd) {
	case "PING":
		sendMessage(conn, "PONG")
		return
	case "ECHO":
		sendMessage(conn, command.arguments[0])
	case "SET":
		fmt.Println(command)
		key, value := command.arguments[0], command.arguments[1]
		app.storage.SetItem(key, &Data{value: value})
		sendMessage(conn, "Ok")
	case "GET":
		fmt.Println(command)
		key := command.arguments[0]
		value, err := app.storage.GetItem(key)
		if err != nil {
			sendError(conn, err.Error())
			return
		}
		sendMessage(conn, value.value)
	default:
		sendError(conn, "command not found")
		return
	}

}
