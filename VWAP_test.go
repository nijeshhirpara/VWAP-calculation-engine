package main

import (
	"VWAPEngine/workers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test_appendToVWAPList(t *testing.T) {
	// Open our jsonFile
	jsonFile, err := os.Open("./test/test_feed.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened test_feed.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	for _, feed := range result["feed"].([]interface{}) {
		fd, err := workers.PrepareFeed([]byte(feed.(string)))
		if err != nil {
			t.Error("Error sending to FeedChannel: ", err)
			return
		}
		appendToVWAPList(fd)
	}

	products := result["VWAP"].(map[string]interface{})
	for product, v := range FeedList {

		if _, ok := products[product]; !ok {
			t.Error("Couldn't validate VWAP calculation: test data tempered")
			return
		}

		res := fmt.Sprintf("%.3f", v.value)
		if products[product] != res {
			t.Errorf("Couldn't validate VWAP calculation: %s - Expected %s, Actual %s", product, products[product], res)
		}
	}
}
