package main

import (
    "os"
	"fmt"
	"log"
	"strconv"
	"math/big"
	"math/rand"
	"net/http"
	"io/ioutil"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
)
func checkBalance(address string, ch chan int) {
	//queryComp := "https://blockexplorer.com/api/addr/" + address + "/balance"
	queryComp := "https://blockchain.info/q/addressbalance/" + address
	resp, err := http.Get(queryComp)
	if err != nil {
		log.Fatalf("Checking balance (uncomp): %s\n", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bodystring := string(body)
	balance, _ := strconv.Atoi(bodystring)
	ch <- balance
}

func generateSeedAddress() []byte {
	padded := make([]byte, 32)
	for i := 0; i < 32; i++ {
		padded[i] = byte(rand.Intn(256))
	}
	return padded
}

func main() {
    ch := make(chan int)
	file, err := os.OpenFile("result.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	
	// Initialise big numbers with small numbers
	count, one := big.NewInt(0), big.NewInt(1)

	// Create a slice to pad our count to 32 bytes
	//padded := make([]byte, 32)
	
	// Loop forever because we're never going to hit the end anyway
	for {
		// Increment our counter
		count.Add(count, one)

		// Copy count value's bytes to padded slice
		copy(generateSeedAddress()[32-len(count.Bytes()):], count.Bytes())

		// Get public key
		privkey, public := btcec.PrivKeyFromBytes(btcec.S256(), generateSeedAddress())
        wif, err := btcutil.NewWIF(privkey, &chaincfg.MainNetParams, false)
		// Get addresses
		uaddr, _ := btcutil.NewAddressPubKey(public.SerializeUncompressed(), &chaincfg.MainNetParams)

		var uncom_balance int
		if err != nil {
			log.Fatalf("Checking balance (comp): %s\n", err)
		}
		go checkBalance(uaddr.EncodeAddress(), ch)
		uncom_balance = <-ch

		fmt.Println(wif.String() + " " +  uaddr.EncodeAddress() + " " + strconv.Itoa(uncom_balance))
		if uncom_balance > 0 {
			var balance int
			balance = uncom_balance
			file.WriteString(wif.String() + " " +  uaddr.EncodeAddress() + " " + strconv.Itoa(balance) + "\n")
		}
	}	
}
