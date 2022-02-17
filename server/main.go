package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var conns []net.Conn

func main() {
	if len(os.Args) != 3 {
		log.Fatalln("Insufficient arguments: [host] [port]")
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", os.Args[1], os.Args[2]))
	if err != nil {
		log.Fatalln(err)
	}
	defer ln.Close()

	fmt.Printf("Listening on %s\n", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go handleRequest(conn)
	}
}

// handleRequest Handle client request
//  @param1 (conn): connection between client and server
//
//  @return1 (err): error variable
func handleRequest(conn net.Conn) (err error) {
	defer conn.Close()

	mess := fmt.Sprintln(" - Welcome to chat - ")
	mess += fmt.Sprint("Enter your name: ")

	// write message
	_, err = conn.Write([]byte(mess))
	if err != nil {
		log.Fatalln(err)
	}

	reply := make([]byte, 1024)

	// read user name
	res := bufio.NewReader(conn)
	n, err := res.Read(reply)
	if err != nil {
		log.Fatal(err)
	}

	name := reply[:n-1]

	conns = append(conns, conn)
	fmt.Printf("%s connected\n", name)

	mess = fmt.Sprintf(" - %s connected - \n", name)
	mess += fmt.Sprintf(" - %d connected users - \n", len(conns))

	// write message to all connections
	for _, element := range conns {
		_, err = element.Write([]byte(mess))
		if err != nil {
			log.Fatalln(err)
		}
	}

	for {
		reply = make([]byte, 1024)

		// read text to chat
		n, err = res.Read(reply)
		if err != nil {
			if err == io.EOF {

				// remove connection from chat
				for n, element := range conns {
					if conn == element {
						conns = append(conns[:n], conns[n+1:]...)
					}

					mess = fmt.Sprintf(" - Bye %s - \n", name)
					mess += fmt.Sprintf(" - %d connected users - \n", len(conns)-1)

					_, err = element.Write([]byte(mess))
					if err != nil {
						log.Fatalln(err)
					}
				}

				fmt.Printf("%s offline\n", name)

				break
			} else {
				log.Fatal(err)
			}
		}

		// write message to all connections
		for _, element := range conns {
			_, err = element.Write([]byte(fmt.Sprintf("%s: %s", name, reply[:n])))
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	return
}
