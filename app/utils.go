package main

import (
	"fmt"
	"net"
)

func nilChecker(err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func sendMessage(conn net.Conn, message string) {
	conn.Write([]byte(fmt.Sprintf("+%v\r\n", message)))
}

// func sendEmptyMessage(conn net.Conn) {
// 	conn.Write([]byte("$0\r\n\r\n"))
// }

// func sendNilMessage(conn net.Conn) {
// 	conn.Write([]byte("*-1\r\n"))
// }

func sendError(conn net.Conn, message string) {
	conn.Write([]byte(fmt.Sprintf("-%v\r\n", message)))
}
