# CMPE273_StockTradingEngine

A virtual stock trading system for whoever wants to learn how to invest in stocks.

This project makes use of Yahoo Finance API. The system has two features:
  1. Buying Stocks
  2. Checking portfolio (Gain/Loss)

## Usage

### Install

```
go get github.com/PrasannaGajbhiye/CMPE273_StockTradingEngine
```

Start the  server:

```
cd CMPE273_StockTradingEngine
go run server.go
```

### Start the client 
#### Buying Stocks
```
cd CMPE273_StockTradingEngine
go run client.go "GOOG:100%" 2000
```
Following will be the response for the above request:
```
TradeId: XXXXXX
GOOG:X:$XXX.XX
UnvestedAmount: $XXX.XX
```

#### Checking Portfolio

```
cd CMPE273_StockTradingEngine
go run client.go XXXXXX
```
Following will be the response for the above request:
```
GOOG:X:$XXX.XX
Current Market Value: $XXX.XX
Unvested Amount: $XXX.XX
```
