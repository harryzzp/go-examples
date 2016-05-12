package main

import (
	"flag"
	"github.com/nats-io/nats"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	//gnatsd -m 8222
	flag.Parse()

	//natsConnection, _ := nats.Connect(nats.DefaultURL)

	//natsConnection, _ := nats.Connect("nats://192.168.99.100:4222")

	slc := []string{"nats://192.168.99.100:4222", "nats://192.168.99.101:5222"}

	opts := nats.Options{
		Servers:        slc,
		AllowReconnect: true,
		MaxReconnect:   10,
		ReconnectWait:  5 * time.Second,
		Timeout:        1 * time.Second,
	}
	natsConnection, _ := opts.Connect()

	defer natsConnection.Close()
	log.Println("Connected to NATS server: " + nats.DefaultURL)

	start := time.Now()

	cnt := 0

	var msg *nats.Msg
	var str string

	for i := 0; i < 300; i++ {
		cnt++

		str = "measurementName,tag1key=" + randSeq(320) +
			" fieldname=1 " + strconv.FormatInt(time.Now().UnixNano(), 10)

		msg = &nats.Msg{Subject: "telegraf", Reply: "bar", Data: []byte(str)}

		natsConnection.PublishMsg(msg)

		log.Println(cnt)
	}

	elapsed := time.Since(start)
	log.Printf("Time took %s", elapsed)
}