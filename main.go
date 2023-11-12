package main

import (
	"fmt"
	"ggcache/cache"
	"log"
	"net"
	"time"
)

func main() {
	opts := ServerOpts{
		ListenAddr: ":3000",
		IsLeader:   true,
	}
	go func() {
		time.Sleep(time.Second * 2)
		conn, err := net.Dial("tcp", ":3000") // dial just for testing
		if err != nil {
			log.Fatal(err)
		}
		conn.Write([]byte("SET Foo Bar 3400000000"))

		time.Sleep(time.Second * 2)
		conn.Write([]byte("GET Foo"))

		buf := make([]byte, 1000)
		n, _ := conn.Read(buf)
		fmt.Println(string(buf[:n]))
	}()

	server := NewServer(opts, cache.New())
	server.Start()
}
