package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Insufficient arguments: [host] [port]")
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", os.Args[1], os.Args[2]))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// read server message
	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", reply)

	// read name
	n, err := bufio.NewReader(os.Stdin).Read(reply)
	if err != nil {
		log.Fatal(err)
	}

	// write name on connection
	_, err = conn.Write(reply[:n])
	if err != nil {
		log.Fatal(err)
	}

	// read server messages
	n, err = conn.Read(reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", reply[:n])

	wg := sync.WaitGroup{}

	// write on connection
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		for {
			// read name
			n, err := bufio.NewReader(os.Stdin).Read(reply)
			if err != nil {
				log.Fatal(err)
			}

			// write message on connection
			_, err = conn.Write(reply[:n])
			if err != nil {
				log.Fatal(err)
			}
		}
	}(&wg)

	// read from connection
	for {
		// read messages
		n, err = conn.Read(reply)
		if err != nil {
			log.Fatal(err)
			break
		}

		fmt.Printf("%s", reply[:n])
	}

	wg.Wait()

}
