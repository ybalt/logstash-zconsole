package main

import (
	"fmt"
	"net/http"
	"os"
	zmq "github.com/pebbe/zmq4"
	"strconv"
	"time"
	"encoding/json"
	"log"
)

type Colorizer map[string]int

const Version string = "v0.1"

// returns up to 14 color escape codes (then repeats) for each unique key
func (c Colorizer) Get(key string) string {
	i, exists := c[key]
	if !exists {
		c[key] = len(c)
		i = c[key]
	}
	bright := "1;"
	if i % 14 > 6 {
		bright = ""
	}
	return "\x1b[" + bright + "3" + strconv.Itoa(7 - (i % 7)) + "m"
}

func handler(w http.ResponseWriter, r *http.Request) {
	var colors Colorizer
	var msg map[string]interface{}
	queue := make(chan string, 100)  //buffer only 100 messages, need to be in config
	go subscriber_task(queue)
	time.Sleep(time.Second) //get some time to socket prepare operations

	colors = make(Colorizer)

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(r.Host))
	for elem := range queue {
		fmt.Printf("queue: %s\n",elem)
		err := json.Unmarshal([]byte(elem), &msg)
		if err != nil {
			fmt.Println("error:", err)
		}
		container_name := msg["container_name"].(string)
		message := msg["message"].(string)
		time := msg["@timestamp"].(string)
		format := "%s%" + strconv.Itoa(len(container_name)) + "s|%s|%s\x1b[0m\n"
		output := fmt.Sprintf(format, colors.Get(container_name),
					time,container_name,message)
		w.Write([]byte(output))
		w.(http.Flusher).Flush()
	}

}

//func publisher_task() {
//
//	go func() {
//		publisher, err := zmq.NewSocket(zmq.PUB)
//		if err != nil {
//			fmt.Println("new socket error", err)
//			os.Exit(0)
//		}
//		defer publisher.Close()
//		err1 := publisher.Bind("tcp://127.0.0.1:5561")
//		if err1 != nil {
//			fmt.Println("bind error", err1)
//			os.Exit(0)
//		}
//		for request_nbr := 1; true; request_nbr++ {
//			time.Sleep(time.Second)
//			_, err := publisher.Send(fmt.Sprintf("request #%d", request_nbr), 0)
//			if err != nil {
//				fmt.Println("send err", err)
//				os.Exit(0)
//			}
//		}
//	}()
//
//}

func subscriber_task(queue chan string) {
	go func() {
		subscriber, err := zmq.NewSocket(zmq.SUB)
		if err != nil {
			fmt.Println("new socket error", err)
			os.Exit(0)
		}
		defer subscriber.Close()
		logstash_addr := os.Getenv("LOGSTASH_ADDR")
		if logstash_addr == "" {
			logstash_addr = "tcp://local-logstash:12300"
		}
		fmt.Printf("# using %s as logstash endpoint\n", logstash_addr)
		err1 := subscriber.Connect(logstash_addr)
		if err1 != nil {
			fmt.Println("connect error", err1)
			os.Exit(0)
		}
		err2 := subscriber.SetSubscribe("")
		if err2 != nil {
			fmt.Println("subscibe error", err1)
			os.Exit(0)
		}

		poller := zmq.NewPoller()
		poller.Add(subscriber, zmq.POLLIN)
		for {
			sockets, _ := poller.Poll(-1)
			update, _ := sockets[0].Socket.Recv(0)
			if update != "" {
				queue <- update
			}
		}
	}()

}

func main() {

	fmt.Printf("# logstash-zconsole %s by ybalt\n", Version)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
