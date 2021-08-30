package main

import (
	"VWAPEngine/workers"
	"container/list"
	"fmt"
	// "sort"
)

type VWAP struct {
	value       float64
	totalPV     float64
	totalVolume float64
	queue       *list.List
}

var FeedList map[string]VWAP

func init() {
	FeedList = make(map[string]VWAP)
}

func appendToVWAPList(fd workers.FeedData) {
	// check if key exists, if not create a new list for the key
	if _, ok := FeedList[fd.ProductId]; !ok {
		FeedList[fd.ProductId] = VWAP{0, 0, 0, list.New()}
	}

	VWAPFeed := FeedList[fd.ProductId]

	// Calculate the VWAP per trading pair using a sliding window of 200 data points.
	// Meaning, when a new data point arrives through the websocket feed the oldest data point will fall off
	// and the new one will be added such that no more than 200 data points are included in the calculation.
	if VWAPFeed.queue.Len() >= 200 {
		// Dequeue
		frontFeed := VWAPFeed.queue.Front()
		frontFD := frontFeed.Value.(workers.FeedData)

		// remove from VWAP calculation
		VWAPFeed.removeFeed(frontFD)

		// This will remove the allocated memory and avoid memory leaks
		VWAPFeed.queue.Remove(frontFeed)
	}

	// enqueue
	VWAPFeed.queue.PushBack(fd)

	// add to VWAP calculation
	VWAPFeed.addFeed(fd)

	FeedList[fd.ProductId] = VWAPFeed
}

// removeFeed removes feed from VWAP calculation
func (v *VWAP) removeFeed(fd workers.FeedData) {
	v.totalPV = v.totalPV - (fd.Price * fd.Volume)
	v.totalVolume = v.totalVolume - fd.Volume
}

// addFeed adds newly arrived feed to VWAP calculation
func (v *VWAP) addFeed(fd workers.FeedData) {
	v.totalPV = v.totalPV + (fd.Price * fd.Volume)
	v.totalVolume = v.totalVolume + fd.Volume
	v.value = v.totalPV / v.totalVolume
}

// printVWAP prints products with their realtime VWAP values
func printVWAP() {
	result := ""
	for _, product := range workers.ProductIds {
		if _, ok := FeedList[product]; ok {
			result += fmt.Sprintf("%s: %.3f [DataPoints: %d] \t ", product, FeedList[product].value, FeedList[product].queue.Len())
		}
	}
	fmt.Printf("\r%s", result)
}
