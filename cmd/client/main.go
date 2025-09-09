package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

// Клиент для общения через консоль
func main() {
	fmt.Print("Введите IP сервера и порт (пример localhost:9000): ")
	var addr string
	fmt.Scanln(&addr)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	go func() {
		serverReader := bufio.NewReader(conn)
		for {
			msg, err := serverReader.ReadString('\n')
			if err != nil {
				fmt.Println("Соединение закрыто сервером")
				os.Exit(0)
			}
			fmt.Print(msg)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		fmt.Fprintln(conn, text)
	}
}
