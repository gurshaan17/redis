package server

import (
	"bufio"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

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

	lsnr, err := net.Listen(
		"tcp",
		config.Host+":"+strconv.Itoa(config.Port),
	)
	if err != nil {
		panic(err)
	}
	defer lsnr.Close()

	for {
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

		go handleClient(c, &conClients)
	}
}

func handleClient(c net.Conn, conClients *int) {
	defer func() {
		c.Close()
		*conClients--
		log.Println("client disconnected:", c.RemoteAddr())
	}()

	reader := bufio.NewReader(c)
	buffer := make([]byte, 0)

	for {
		tmp := make([]byte, 1024)
		n, err := reader.Read(tmp)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println("read error:", err)
			return
		}

		// append new data to buffer
		buffer = append(buffer, tmp[:n]...)

		for {
			value, delta, err := core.DecodeOne(buffer)
			if err != nil {
				// not enough data yet → wait for more
				break
			}

			buffer = buffer[delta:]

			arr, ok := value.([]interface{})
			if !ok || len(arr) == 0 {
				c.Write([]byte("-ERR invalid command\r\n"))
				continue
			}

			tokens := make([]string, len(arr))
			for i := range arr {
				s, ok := arr[i].(string)
				if !ok {
					c.Write([]byte("-ERR invalid argument type\r\n"))
					continue
				}
				tokens[i] = s
			}

			cmd := &core.RedisCmd{
				Cmd:  strings.ToUpper(tokens[0]),
				Args: tokens[1:],
			}

			if err := core.EvalAndRespond(cmd, c); err != nil {
				c.Write([]byte("-" + err.Error() + "\r\n"))
			}
		}
	}
}
