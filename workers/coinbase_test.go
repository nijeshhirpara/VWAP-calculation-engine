package workers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test_SendToFeedChannel(t *testing.T) {
	// Open our jsonFile
	jsonFile, err := os.Open("../test/test_feed.json")
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
		if _, err := SendToFeedChannel([]byte(feed.(string))); err != nil {
			t.Error("Error sending to FeedChannel: ", err)
		}
	}
}
