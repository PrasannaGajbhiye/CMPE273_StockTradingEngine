package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strconv"
	"strings"
	"math/rand"
	"time"
)

type ArgsBuyingStocks struct {
    StockList []string
	StockShare []string
	TransBudget float32
}

type ReplyBoughtStocks struct {
    TradeId int64
	StocksList []string
	StocksCount []int32
	StocksPrice []float32 
	UnvestedAmount float32
}

type ArgsGainLoss struct{
	TradeId int64
}

type ReplyGainLoss struct{
	StocksList []string
	StocksCount []int32
	StocksPrice []float32
	StocksGainLoss []string
	UnvestedAmount float32
	CurrentMarketValue float32
}

type Response struct {
  List struct {
    Resources []struct {
      Resource struct {
        Fields struct {
          Name    string `json:"name"`
          Price   string `json:"price"`
          Symbol  string `json:"symbol"`
          Ts      string `json:"ts"`
          Type    string `json:"type"`
          UTCTime string `json:"utctime"`
          Volume  string `json:"volume"`
        } `json:"fields"`
      } `json:"resource"`
    } `json:"resources"`
  } `json:"list"`
}

type PreviousTradesList struct {
	tradeDetails []*ReplyBoughtStocks
}

var previousTradeList PreviousTradesList

type StockEngine struct{}

func (s *StockEngine) GetTradeDetails(args *ArgsGainLoss, reply *ReplyGainLoss) error{
	isFound:=false
	for l:=0;l<len(previousTradeList.tradeDetails);l++{
		
		var trade *ReplyBoughtStocks
		trade = previousTradeList.tradeDetails[l]
		
		if(trade.TradeId==args.TradeId){
			isFound=true
			reply.CurrentMarketValue = 0.0
			for i:=0; i< len(trade.StocksList); i++{	
				resp, err := http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+trade.StocksList[i]+"/quote?format=json")
				if err != nil {
					fmt.Println("Error")
				}
				defer resp.Body.Close()
				
				body, _ := ioutil.ReadAll(resp.Body)
				
				var msgRes Response
				_ = json.Unmarshal(body, &msgRes)
				
				reply.StocksList=append(reply.StocksList,msgRes.List.Resources[0].Resource.Fields.Symbol)
				
				
				stkPrice,_:= strconv.ParseFloat(msgRes.List.Resources[0].Resource.Fields.Price,64)
				reply.StocksPrice=append(reply.StocksPrice,float32(stkPrice))
				
				if(float32(stkPrice) == trade.StocksPrice[i]){
					reply.StocksGainLoss=append(reply.StocksGainLoss,"") 
				}else if(float32(stkPrice)<trade.StocksPrice[i]){
					reply.StocksGainLoss=append(reply.StocksGainLoss,"-") 
				}else{
					reply.StocksGainLoss=append(reply.StocksGainLoss,"+") 
				}
				
				reply.StocksCount=append(reply.StocksCount,trade.StocksCount[i])
				
				reply.CurrentMarketValue= reply.CurrentMarketValue + float32(reply.StocksCount[i]) *reply.StocksPrice[i]
			}
			reply.UnvestedAmount= trade.UnvestedAmount
		}
	}
	if(isFound==false){
		log.Fatal("No such transaction available!") 
	}
	return nil
}


func (s *StockEngine) GenerateTransId(args *ArgsBuyingStocks, reply *ReplyBoughtStocks) error{
	shouldBreak:=false
	rand.Seed(time.Now().UTC().UnixNano())
	reply.TradeId = rand.Int63n(999999)
	reply.UnvestedAmount = 0.0
	
	for i:=0; i< len(args.StockList); i++{
		
		resp, err := http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+args.StockList[i]+"/quote?format=json")
		if err != nil {
			fmt.Println("Error")
		}
		defer resp.Body.Close()
		
		body, _ := ioutil.ReadAll(resp.Body)
		
		var msgRes Response
		_ = json.Unmarshal(body, &msgRes)
		
		if(len(msgRes.List.Resources)!=0){
			reply.StocksList=append(reply.StocksList,msgRes.List.Resources[0].Resource.Fields.Symbol)
		
		
			stkPrice,_:= strconv.ParseFloat(msgRes.List.Resources[0].Resource.Fields.Price,64)
			reply.StocksPrice=append(reply.StocksPrice,float32(stkPrice))
			
			stkShare,_:= strconv.ParseFloat(strings.Replace(args.StockShare[i],"%","",-1),64)
			
			indAvlBud := (float32(stkShare) * args.TransBudget)/100.0
		
			reply.StocksCount=append(reply.StocksCount,int32(indAvlBud/float32(reply.StocksPrice[i])))
			
			indUnvestedBalance := indAvlBud - float32(reply.StocksCount[i])* float32(reply.StocksPrice[i])
			
			reply.UnvestedAmount=reply.UnvestedAmount+indUnvestedBalance
		}else{
			shouldBreak = true
			log.Fatal("Invalid Inputs : ","Stocks ticker(s)") 
		}
	}
	
	if(shouldBreak==false){
		previousTradeList.tradeDetails = append(previousTradeList.tradeDetails,reply)
	}else{
		log.Fatal("Invalid Inputs : ","Stocks ticker(s)") 
	}
	
	return nil
}


func main(){
	go StartServer()
	var input string
	fmt.Scanln(&input)
}

func StartServer() {
	
	stockEng :=new(StockEngine)
	
	server := rpc.NewServer()
	
	server.Register(stockEng)

    l, e := net.Listen("tcp", ":8222")
    if e != nil {
		log.Fatal("listen error:", e)
    }

    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go server.ServeCodec(jsonrpc.NewServerCodec(conn))
    }
}