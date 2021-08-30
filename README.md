# VWAP calculation engine

The goal of this project is to create a real-time VWAP (volume-weighted average price) calculation engine. It uses the coinbase websocket feed to stream in trade executions and update the VWAP for each trading pair
as updates become available.

#### To Run

````
go run .
````

#### To Test

````
go test ./...
````

#### Expected Output

````
2021/08/30 18:28:37 connecting to ws-feed.pro.coinbase.com

------Volume-Weighted Average Price-------

BTC-USD: 48114.250 [DataPoints: 40] 	 ETH-USD: 3195.450 [DataPoints: 14] 	 ETH-BTC: 0.066 [DataPoints: 2]
````

## Considerations

- Calculated VWAP for three trading pairs for demonstration purposes **[BTC-USD, ETH-USD, ETH-BTC]**
- Calculated the VWAP per trading pair using a sliding window of **200 data points**. Meaning, when a new
data point arrives through the websocket feed the oldest data point will fall off and the new one will be
added such that no more than 200 data points are included in the calculation.
- It streams the resulting VWAP values on each websocket update.
- Calculated VWAP with **3 decimal points**

## Coinbase Websocket - Match channel

#### Request to subscribe

````
{
    "type": "subscribe",
    "channels": [{ "name": "matches", "product_ids": ["BTC-USD", "ETH-USD", "ETH-BTC"] }]
}
````

#### Response

````
{
  "type":"match",
  "trade_id":206227069,
  "maker_order_id":"288fb6c9-fc6f-4ce6-a013-0cfcc87ddacd",
  "taker_order_id":"d414ae8e-3ad5-4cef-a701-d9e18078420a",
  "side":"sell",
  "size":"0.00004761",
  "price":"47979.82",
  "product_id":"BTC-USD",
  "sequence":28750468641,
  "time":"2021-08-30T06:19:25.002709Z"
}
````

## VWAP calculation

````
VWAP= (∑ Price * Volume) / ∑ Volume​
````
Note: Size from the websocket response is considered as volume for the calculation.

Example,

VWAP based on 3 data points
````
{
  "size":"0.00004761",
  "price":"47979.82",
  "product_id":"BTC-USD"
},
{
  "size":"0.15995239",
  "price":"47979.82",
  "product_id":"BTC-USD"
},
{
  "size":"0.253",
  "price":"47980",
  "product_id":"BTC-USD"
}
````
````
// A sum of price-volume product
totalPV = (47979.82 * 0.00004761) + (47979.82 * 0.15995239) + (47980 * 0.253) = 19815.7112

// A sum of volume
totalVolume = 0.00004761 + 0.15995239 + 0.253 = 0.413

VWAP = totalPV/totalVolume = 47979.930266344
````
