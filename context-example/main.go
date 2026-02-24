package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go handleSignals(cancel)

	if err := startServer(ctx); err != nil {
		log.Fatal(err)
	}
}

func handleSignals(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, os.Interrupt)
	for {
		sig := <-sigCh
		switch sig {
		case os.Interrupt:
			cancel()
			return
		}
	}
}

func startServer(ctx context.Context) error {
	laddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}

	defer l.Close()

	for {
		select {
		case <-ctx.Done():
			log.Println("server stopped")
			return nil
		default:
			if err := l.SetDeadline(time.Now().Add(time.Second)); err != nil {
				return err
			}

			_, err := l.Accept()
			if err != nil {
				if os.IsTimeout(err) {
					continue
				}
			}

			log.Println("new client connected")
		}
		
	}
}