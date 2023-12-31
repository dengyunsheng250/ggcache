package main

import (
	"context"
	"fmt"
	"ggcache/cache"
	"log"
	"net"
)

type ServerOpts struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	ServerOpts
	cache cache.Cacher
}

func NewServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return fmt.Errorf("listen error: %s", err)
	}
	log.Printf("server starting on port [%s]\n", s.ListenAddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %s\n", err)
			continue // 如果break,其他人就无法连接了
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("cannot read error: %s\n", err)
			break
		}

		// msg := buf[:n]
		// fmt.Println(string(msg))

		go s.handleCmd(conn, buf[:n])
	}
}

func (s *Server) handleCmd(conn net.Conn, rawCmd []byte) {
	msg, err := parseMessage(rawCmd)
	if err != nil {
		fmt.Println("failed to parse command:", err)
		// respond
		conn.Write([]byte(err.Error()))
		return
	}
	switch msg.Cmd {
	case CMDSet:
		err = s.handleSetCmd(conn, msg)
	case CMDGet:
		err = s.handleGetCmd(conn, msg)
	}
	if err != nil {
		fmt.Println("failed to handle command:", err)
		// respond
		conn.Write([]byte(err.Error()))
		return
	}
}

func (s *Server) handleGetCmd(conn net.Conn, msg *Message) error {
	s.cache.Get(msg.Key)
	return nil
}

func (s *Server) handleSetCmd(conn net.Conn, msg *Message) error {
	// fmt.Println("handling the set command: ", msg)
	_ = s.cache.Set(msg.Key, msg.Value, msg.TTL)
	return nil
}

func (s *Server) sendToFollowers(ctx context.Context, msg *Message) error {
	return nil
}
