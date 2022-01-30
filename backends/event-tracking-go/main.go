package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var rpcEndpoint = "https://rinkeby.infura.io/v3/c75f2ce78a4a4b64aa1e9c20316fda3e"

// ERC721 contract for tracking
// For this demo use contract in "contracts/GameCharacter.sol"
var contractAddress = common.HexToAddress("0x976d91f3ec4017db4b17b19a35e4862b5f7ed8c4")

// Blocknumber of starting to track events
var startBlock = int64(10018020)

func convertToType(val interface{}, t abi.Type) (interface{}, error) {
	switch t.T {
	case abi.StringTy:
		return string(val.(common.Hash).String()), nil
	case abi.AddressTy:
		return common.BytesToAddress(val.(common.Hash).Bytes()), nil
	case abi.HashTy:
		return common.BytesToHash(val.(common.Hash).Bytes()), nil
	default:
		return nil, fmt.Errorf("abi: unknown type %v", t.T)
	}
}

func main() {
	eventSignatures := map[string]common.Hash{}

	client, err := ethclient.Dial(rpcEndpoint)

	if err != nil {
		log.Fatal(err)
	}

	rdr, err := os.ReadFile("GameCharacter.abi.json")

	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(rdr)))
	if err != nil {
		log.Fatal(err)
	}

	for _, eventObj := range contractAbi.Events {
		eventSignatures[eventObj.Name] = crypto.Keccak256Hash([]byte(eventObj.Sig))
	}

	// Create a for loop, which is breakable by CTRL+C
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	for {
		headBlock, err := client.HeaderByNumber(context.Background(), nil)

		lastBlock := headBlock.Number.Int64()

		if lastBlock > startBlock+5000 {
			lastBlock = startBlock + 5000
		}
		fmt.Printf("Getting events between, %d - %d blocks\n", startBlock, lastBlock)

		query := ethereum.FilterQuery{
			FromBlock: big.NewInt(startBlock),
			ToBlock:   big.NewInt(lastBlock),
			Addresses: []common.Address{
				contractAddress,
			},
		}

		logs, err := client.FilterLogs(context.Background(), query)

		if err != nil {
			log.Print(err)
			fmt.Println("retry in 3sec...")
			time.Sleep(3 * time.Second)
		} else {
			for _, vLog := range logs {

				eventData := map[string]interface{}{}
				topicData := map[string]interface{}{}
				activeEvent := abi.Event{}

				for _, eventObj := range contractAbi.Events {
					if eventSignatures[eventObj.Name].Hex() == vLog.Topics[0].Hex() {
						activeEvent = eventObj
					}
				}

				contractAbi.UnpackIntoMap(eventData, string(activeEvent.Name), vLog.Data)
				abi.ParseTopicsIntoMap(topicData, activeEvent.Inputs, vLog.Topics)

				i := 0
				for _, input := range activeEvent.Inputs {

					if input.Indexed {
						i += 1
						eventData[input.Name], err = convertToType(vLog.Topics[i], input.Type)
					}

				}
				eventData["txid"] = vLog.TxHash
				eventData["blocknumber"] = vLog.BlockNumber

				fmt.Printf("%s event fired on block %d, with parameters ( %v )\n", activeEvent.Name, vLog.BlockNumber, eventData)
			}

			startBlock = lastBlock

			fmt.Println("sleeping 10sec...")
			time.Sleep(10 * time.Second)
		}
	}
}
