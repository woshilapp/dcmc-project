package main

import (
	"fmt"
	"syscall"

	"os"
	"os/signal"

	"github.com/woshilapp/dcmc-project/server/network"
)

func main() {
	fmt.Println("Hello, world server!")

	listener, err := network.ListenServer(":7789")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Listening on Port tcp:7789")

	errchan := make(chan error, 100)

	go network.ListenUDP(":7789", errchan)
	fmt.Println("Listening on Port udp:7789")

	go network.AccpetConn(listener, errchan)

	go func() {
		for {
			err := <-errchan
			fmt.Println("[ERROR]", err)
		}
	}()

	//stop it
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		<-exitChan

		fmt.Println("Interrupt")
		break
	}
}
