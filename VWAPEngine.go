package main

import (
	"VWAPEngine/workers"
	"fmt"
	"log"
	"os"
	"os/signal"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	done := make(chan bool, 1)
	defer close(done)

	// start the CoinBase websocket connection and receive feed via feedChannel
	coinbaseConn, feedChannel, err := workers.StartCoinbase()
	if err != nil {
		log.Println("Error starting CoinBase Worker:", err)
		done <- true
	}
	defer close(feedChannel)

	fmt.Print("\n------Volume-Weighted Average Price-------\n\n")

	for {
		select {
		case feed := <-feedChannel:
			appendToVWAPList(feed)
			printVWAP()
		case <-interrupt:
			log.Println("interrupt")
			if err := workers.StopCoinbase(coinbaseConn); err != nil {
				log.Println("Error in smooth stopping of CoinBase Worker:", err)
			}
			return
		case <-done:
			return
		}
	}
}
