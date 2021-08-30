package workers

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
)

const (
	host        = "ws-feed.pro.coinbase.com"
	scheme      = "wss"
	channelName = "matches"
)

type feedData struct {
	ProductId string
	Price     float64
	Volume    float64
}

type channel struct {
	Name        string
	Product_ids []string
}

type channelReq struct {
	Type     string
	Channels []channel
}

func StartCoinbase() (*websocket.Conn, chan feedData, error) {
	// feedChannel to receive messages from feed
	feedChannel := make(chan feedData, 100)

	// open websocket
	log.Printf("connecting to %s", host)
	var addr = flag.String("addr", host, "https service address")
	u := url.URL{Scheme: scheme, Host: *addr, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return c, feedChannel, err
	}

	// subscribe to the channel
	var req channelReq
	req.Type = "subscribe"
	req.Channels = []channel{channel{channelName, []string{"BTC-USD", "ETH-USD", "ETH-BTC"}}}
	err = c.WriteJSON(&req)
	if err != nil {
		return c, feedChannel, err
	}

	// watch feed data and read the messages
	// go feedWatcher(c, feedChannel)
	go func() {

		defer c.Close()

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			var result map[string]interface{}
			json.Unmarshal([]byte(message), &result)

			if result["type"] == "match" {
				priceStr, _ := result["price"].(string)
				sizeStr, _ := result["size"].(string)

				price, priceErr := strconv.ParseFloat(priceStr, 64)
				if priceErr != nil {
					price = 0
				}

				size, sizeErr := strconv.ParseFloat(sizeStr, 64)
				if sizeErr != nil {
					size = 0
				}

				var fd feedData
				fd.ProductId = result["product_id"].(string)
				fd.Price = price
				fd.Volume = size

				feedChannel <- fd
			}
		}
	}()

	return c, feedChannel, nil
}

func StopCoinbase(c *websocket.Conn) error {
	if c == nil {
		return errors.New("No active connection supplied")
	}
	// subscribe to the channel
	var req channelReq
	req.Type = "unsubscribe"
	req.Channels = []channel{channel{channelName, []string{"BTC-USD", "ETH-USD", "ETH-BTC"}}}

	err := c.WriteJSON(&req)
	if err != nil {
		return err
	}
	return nil
}
