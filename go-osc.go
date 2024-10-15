package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hypebeast/go-osc/osc"
	"github.com/nats-io/nats.go"
)

type Msg struct {
	Timestamp              time.Time `json:"Rimestamp"`
	Power                  int       `json:"Power"`
	AccumulatedConsumption float64   `json:"AccumulatedConsumption"`
	AccumulatedCost        float64   `json:"AccumulatedCost"`
	Currency               string    `json:"Currency"`
	MinPower               int       `json:"MinPower"`
	AveragePower           float64   `json:"AveragePower"`
	MaxPower               float64   `json:"MaxPower"`
}

func main() {
	log.Println("Starting GO-Osc program")
	client := osc.NewClient("localhost", 57120)

	log.Println("Connect to NATS")
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()

	sub, _ := nc.SubscribeSync("tibber.>")
	for {
		jsonMsg, err := sub.NextMsg(5 * time.Second)
		if err == nil {
			var msg Msg
			err := json.Unmarshal(jsonMsg.Data, &msg)
			if err != nil {
				log.Fatalf("Error unmarshalling JSON: %v", err)
			}
			fmt.Println("Sent message to osc with value:", msg.Power)
			sendOsc(client, msg.Power)
		} else {
			fmt.Println("NextMsg timed out.")
		}
	}
}

func sendOsc(client *osc.Client, value int) {
	msg := osc.NewMessage("/tibber")
	msg.Append(int32(value))
	err := client.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}
