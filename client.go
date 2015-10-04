package main

import (
	"os"
	"net"
	"net/rpc/jsonrpc"
	"log"
	"fmt"
	"strings"
	"strconv"
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

func main(){
	cmdArgs:=os.Args[1:]
	//	Start Client
	if(len(cmdArgs)==0){
		fmt.Println("No command line arguments provided.")
	}else{
		StartClient(cmdArgs)	
	}
}

func StartClient(indCmdArgs []string)  {
	// connect to the server
	//stocksDesc string,budget float64
	connect, err := net.Dial("tcp", "localhost:8222")
	if err != nil {
		panic(err)
	}
	defer connect.Close()
	
	client := jsonrpc.NewClient(connect)
	
	if(len(indCmdArgs)==2){
		f,_ := strconv.ParseFloat(indCmdArgs[1],64)
		var replyStock ReplyBoughtStocks
		var argsStockStruct ArgsBuyingStocks
		var argsStock *ArgsBuyingStocks
		argsStockStruct.TransBudget = float32(f)
		individualStocks := strings.Split(indCmdArgs[0],",")
		
		k:=0
		totalStocksShareValue:=0.0
		for i:=0; i < len(individualStocks); i++ {
			individualStockInfo:=strings.Split(individualStocks[i],":")
			for j:=0; j< len(individualStockInfo);j++{
				if(j%2==0){
					argsStockStruct.StockList=append(argsStockStruct.StockList,individualStockInfo[j])
				}else{
					argsStockStruct.StockShare=append(argsStockStruct.StockShare,individualStockInfo[j])
					vl,_:=strconv.ParseFloat(strings.Replace(individualStockInfo[j],"%","",-1),64)
					totalStocksShareValue = totalStocksShareValue +  vl
				}
				k++
			}
		}
		
		if(totalStocksShareValue!=100.00){
			log.Fatal("Invalid Inputs : ","Stocks share percentages.")
		}else{
			argsStock = &argsStockStruct
		
			err = client.Call("StockEngine.GenerateTransId",argsStock,&replyStock)
			if err!=nil{
				log.Fatal("stringed error:", err)
			}
			
			fmt.Printf("\nTradeId: %d\n",replyStock.TradeId)
			
			for j:=0;j<len(replyStock.StocksList);j++{
				if(j!=len(replyStock.StocksList)-1){
					fmt.Printf("%s:%d:$%.2f,",replyStock.StocksList[j],replyStock.StocksCount[j], replyStock.StocksPrice[j])	
				}else{
					fmt.Printf("%s:%d:$%.2f\n",replyStock.StocksList[j],replyStock.StocksCount[j], replyStock.StocksPrice[j])	
				}
			}
			fmt.Printf("UnvestedAmount: $%.2f\n",replyStock.UnvestedAmount)
			fmt.Printf("\n")
			
		}
	}else if(len(indCmdArgs)==1){
		
		reqTradeId,_ := strconv.ParseInt(indCmdArgs[0],10,64)
		
		var argsGainLoss *ArgsGainLoss
		var replyGainLoss ReplyGainLoss 
		argsGainLoss = &ArgsGainLoss{reqTradeId}
		
		err = client.Call("StockEngine.GetTradeDetails",argsGainLoss,&replyGainLoss)
		if err!=nil{
			log.Fatal("stringed error:", err)
		}
		
		fmt.Printf("\n")
		for j:=0;j<len(replyGainLoss.StocksList);j++{
			if(j!=len(replyGainLoss.StocksList)-1){
				fmt.Printf("%s:%d:%s$%.2f,",replyGainLoss.StocksList[j],replyGainLoss.StocksCount[j],replyGainLoss.StocksGainLoss[j],replyGainLoss.StocksPrice[j])	
			}else{
				fmt.Printf("%s:%d:%s$%.2f\n",replyGainLoss.StocksList[j],replyGainLoss.StocksCount[j],replyGainLoss.StocksGainLoss[j], replyGainLoss.StocksPrice[j])	
			}
		}
		fmt.Printf("Current Market Value: $%.2f\n",replyGainLoss.CurrentMarketValue)
		fmt.Printf("Unvested Amount: $%.2f\n",replyGainLoss.UnvestedAmount)
		fmt.Printf("\n")
		
				
	}
}