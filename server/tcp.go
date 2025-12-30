package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/gurshaan17/redis/config"
	"github.com/gurshaan17/redis/core"
)

func RunSyncTCPServer() {
	log.Println(
		"starting a synchronous TCP server on",
		config.Host,
		config.Port,
	)

	conClients := 0

	// listen on host:port
	lsnr, err := net.Listen(
		"tcp",
		config.Host+":"+strconv.Itoa(config.Port),
	)
	if err != nil {
		panic(err)
	}
	defer lsnr.Close()

	for {
		// blocking call: wait for a new client
		c, err := lsnr.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		conClients++
		log.Println(
			"client connected:",
			c.RemoteAddr(),
			"concurrent clients:",
			conClients,
		)

		// synchronous: handle client in same goroutine
		handleClient(c, &conClients)
	}
}

func handleClient(c net.Conn, conClients *int) {
	defer func() {
		c.Close()
		*conClients--
		log.Println(
			"client disconnected:",
			c.RemoteAddr(),
			"concurrent clients:",
			*conClients,
		)
	}()

	for {
		cmd, err := readCommand(c)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println("read error:", err)
			return
		}

		// 🔹 Decode RESP here
		value, err := core.Decode([]byte(cmd))
		if err != nil {
			log.Println("RESP decode error:", err)
			c.Write([]byte("-ERR invalid command\r\n"))
			continue
		}

		log.Printf("decoded value: %#v\n", value)

		// 🔹 For now: respond with OK (later: command handling)
		c.Write([]byte("+OK\r\n"))
	}

}

func readCommand(c net.Conn) (string, error) {
	// max 512 bytes per read
	buf := make([]byte, 512)

	n, err := c.Read(buf)
	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func respond(cmd string, c net.Conn) error {
	_, err := c.Write([]byte(cmd))
	return err
}
