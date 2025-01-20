package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Обработка аргументов
	timeoutFlag := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go-telnet [--timeout=10s] host port")
		return
	}
	host := flag.Arg(0)
	port := flag.Arg(1)
	address := net.JoinHostPort(host, port)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Создание клиента
	tClient := NewTelnetClient(address, *timeoutFlag, os.Stdin, os.Stdout)

	if err := tClient.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer tClient.Close()

	go func() {
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := tClient.Receive()
				if err != nil {
					switch {
					case errors.Is(err, syscall.ECONNRESET):
						fmt.Fprintln(os.Stderr, "...Connection closed by peer")
					case errors.Is(err, io.EOF):
						fmt.Fprintln(os.Stderr, "...EOF")
					default:
						fmt.Fprintf(os.Stderr, "Receive error: %v\n", err)
					}
					return
				}
			}
		}
	}()

	go func() {
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := tClient.Send()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Send error: %v\n", err)
					return
				}
			}
		}
	}()

	<-ctx.Done()
	fmt.Fprintln(os.Stderr, "...Terminating")
}
