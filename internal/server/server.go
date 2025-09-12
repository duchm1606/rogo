package server

import (
	"duchm1606/rogo/internal/config"
	"duchm1606/rogo/internal/constant"
	"duchm1606/rogo/internal/core"
	"duchm1606/rogo/internal/core/io_multiplexing"
	"io"
	"log"
	"net"
	"syscall"
	"time"
)

func readCommand(fd int) (*core.Command, error) {
	var buf = make([]byte, 512)

	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, io.EOF
	}
	return core.ParseCommand(buf)
}

// func respond(data string, fd int) error {
// 	if _, err := syscall.Write(fd, []byte(data)); err != nil {
// 		return err
// 	}
// 	return nil
// }

func RunIoMultiplexingServer() {
	log.Println("starting an I/O Multiplexing TCP server on", config.Port)
	listener, err := net.Listen(config.Protocol, config.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()


	// Get the file descriptor from the listener
	tcpListener, ok := listener.(*net.TCPListener)
	if !ok {
		log.Fatal("listener is not a TCPListener")
	}
	listenerFile, err := tcpListener.File()
	if err != nil {
		log.Fatal(err)
	}
	defer listenerFile.Close()

	serverFd := int(listenerFile.Fd())

	// Create an ioMultiplexing server instance
	ioMultiplexingServer, err := io_multiplexing.CreateIOMultiplexer()
	if err != nil {
		log.Fatal(err)
	}

	// Monitor "read" events on the Server FD
	if err = ioMultiplexingServer.Monitor(io_multiplexing.Event{
		Fd: serverFd,
		Op: io_multiplexing.OpRead,
	}); err != nil {
		log.Fatal(err)
	}

	var events = make([]io_multiplexing.Event, config.MaxConnection)
	var lastActiveExpireExecTime = time.Now()


	for {
		// Check last execution time and call if it is more than 100ms ago.
		if time.Now().After(lastActiveExpireExecTime.Add(constant.ActiveExpireFrequency)) {
			core.ActiveDeleteExpiredKeys()
			lastActiveExpireExecTime = time.Now()
		}
		// wait for file descriptor events ready. IO blocking here
		events, err = ioMultiplexingServer.Wait()
		if err != nil {
			log.Fatal(err)
		}

		for i, event := range events {
			// New connection is requesting to serverFd (listener)
			if event.Fd == serverFd {
				log.Println("New client is trying to connect")
				// accept the connection
				conn, _, err := syscall.Accept(serverFd)
				if err != nil {
					log.Println("err", err)
					continue
				}

				// Create a new goroutine to handle the connection
				log.Printf("set up a new connection")

				// ask epoll to monitor this connection
				// Monitor "read" events on the Client FD
				if err = ioMultiplexingServer.Monitor(io_multiplexing.Event{
					Fd: conn,
					Op: io_multiplexing.OpRead,
				}); err != nil {
					log.Fatal(err)
				}
			} else {
				// handle the connection
				cmd, err := readCommand(event.Fd)
				if err != nil {
					if err == syscall.ECONNRESET || err == io.EOF {
						log.Println("client closed the connection")
						_ = syscall.Close(event.Fd)
						continue
					}
					log.Println("err", err)
					continue
				}
				if err = core.ExecuteAndResponse(cmd, events[i].Fd); err != nil {
					log.Println("err write:", err)
				}
			}
		}
	}
}