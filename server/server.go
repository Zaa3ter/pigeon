package server

import (
	// "bufio"
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var users = make(map[string]*net.Conn)

func Listen(port string) {
	l, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("server start Listen on ", l.Addr())
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
		}
		go login(conn)
	}
}

func login(conn net.Conn) {
	input := make([]byte, 10)

	for {
		n, err := conn.Read(input)
		if err != nil {
			conn.Write([]byte{3})
			log.Println(err)
			continue
		}

		username := string(input[:n-1])
		if n > 0 {
			_, ok := users[username]
			if ok {
				conn.Write([]byte{2})
				continue
			}
		}

		users[username] = &conn
		conn.Write([]byte{1})
		log.Println("+++ added user " + username)
		go messageReader(username)
		return
	}
}

func messageReader(username string) {
	println("massage reader")
	conn := *users[username]
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			delete(users, username)
			return
		}

		msg = strings.TrimRight(msg, "\n")
		go hub(username, msg)
	}
}

func hub(username, msg string) {
	fullmsg := "\x1b[32m" + username + ":\x1b[0m " + msg
	for resever := range users {
		if resever == username {
			continue
		}
		fmt.Println("sending to ", resever)
		_, err := (*users[resever]).Write([]byte(fullmsg))
		if err != nil {
			log.Println(err)
			delete(users, resever)
		}
	}
}
