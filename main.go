package main

import (
	"log"
	"net"
	"time"
)

type Job struct {
	conn net.Conn
}

type Worker struct {
	id 			int
	jobChannel 	chan Job
}

type Pool struct {
	jobQueue chan Job
	workers []*Worker
}

func handleConnection(conn net.Conn) {
	defer conn.Close()	
	var buf []byte = make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	// Process req simulation
	time.Sleep(time.Second * 1)
	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 12\r\n\r\nHello world!"))
}

func NewWorker(id int, jobChannel chan Job) *Worker {
	return &Worker{
		id: id,
		jobChannel: jobChannel,
	}
}

func (w *Worker) Start() {
	go func() {
		for job := range w.jobChannel {
			log.Printf("Worker %d is handling job from %s", w.id, job.conn.RemoteAddr())
			handleConnection(job.conn)
		}
	}()
}

func NewPool(numWorkers int) *Pool {
	return &Pool{
		jobQueue: make(chan Job),
		workers: make([]*Worker, numWorkers),
	}
}

func (p *Pool) AddJob(conn net.Conn) {
	p.jobQueue <- Job{conn: conn}
}

func (p *Pool) Start() {
	for i := 0; i < len(p.workers); i++ {
		worker := NewWorker(i, p.jobQueue)
		p.workers[i] = worker
		worker.Start()
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	pool := NewPool(2)
	pool.Start()

	for {	
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("socket = ", conn)
		// go handleConnection(conn)
		pool.AddJob(conn)
	}
}