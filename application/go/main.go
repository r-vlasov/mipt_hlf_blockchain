package main

import (
	"fmt"
	"io/ioutil"
	"time"
    	"encoding/json"
	"log"
	"strings"
	"os"
	"path/filepath"
	"strconv"
	"github.com/chzyer/readline"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)


type Resident struct {
        Firstname       string 	`json:"firstname"`
        Secondname      string 	`json:"secondname"`
        City            string 	`json:"city"`
        Address         string 	`json:"address"`
        Mobile          string	`json:"mobile"`
        Status          string	`json:"status"`
}

type ResidentHistory struct { 
	Content		*Resident	`json:"content"`
	TxId		string	 	`json:"txid"`
	Timestamp	string		`json:"ts"`
}


func connectionInitialize(userName string, channelName string, contractName string) (*gateway.Contract) {
	log.Println("<<< Connection and initialization started")
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		log.Fatalf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
	}
	wallet, err := gateway.NewFileSystemWallet("wallet")
	if err != nil {
		log.Fatalf("Failed to create wallet: %v", err)
	}
	if !wallet.Exists(userName) {
		err = populateWallet(wallet)
		if err != nil {
			log.Fatalf("Failed to populate wallet contents: %v", err)
		}
	}
	ccpPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"connection-org1.yaml",
	)

	gw, err := gateway.Connect(
		gateway.WithConfig(config.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if err != nil {
		log.Fatalf("Failed to connect to gateway: %v", err)
	}

	network, err := gw.GetNetwork(channelName)
	if err != nil {
		log.Fatalf("Failed to get network: %v", err)
	}
	contract := network.GetContract(contractName)

	log.Println(">>> Connection and initialization successfully finished")
	
	defer gw.Close()
	return contract
}

func populateWallet(wallet *gateway.Wallet) error {
	credPath := filepath.Join(
		"..",
		"..",
		"test-network",
		"organizations",
		"peerOrganizations",
		"org1.example.com",
		"users",
		"User1@org1.example.com",
		"msp",
	)

	certPath := filepath.Join(credPath, "signcerts", "User1@org1.example.com-cert.pem")
	// read the certificate pem
	cert, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return err
	}

	keyDir := filepath.Join(credPath, "keystore")
	// there's a single file in this dir containing the private key
	files, err := ioutil.ReadDir(keyDir)
	if err != nil {
		return err
	}
	if len(files) != 1 {
		return fmt.Errorf("keystore folder should have contain one file")
	}
	keyPath := filepath.Join(keyDir, files[0].Name())
	key, err := ioutil.ReadFile(filepath.Clean(keyPath))
	if err != nil {
		return err
	}

	identity := gateway.NewX509Identity("Org1MSP", string(cert), string(key))

	return wallet.Put("appUser", identity)
}


var greetingIn string = "#  "
var greetingOut string = "$ "
var greetingHelp string = "! "

func gotStringFromStdin(greeting string, va *string) {
	r1, err := readline.New(greeting)
	if err != nil {
		log.Fatalf("io error")
	}
	line, err := r1.Readline()
	if err != nil {
		log.Fatalf("io error")
	}
	*va = line
	defer r1.Close()
}

func pushStringToStdout(greeting string, va interface{}) {
	fmt.Printf("%s %s\n", greeting, va)
}

func parseStringDatetimeToUnixtime(dt string) (time.Time, error) {
	seconds := strings.Split(strings.Split(dt, " ")[0], ":")[1]
	unixSec, err := strconv.ParseInt(seconds, 10, 64)
	if err != nil {
		log.Printf("Failed to parse responsed datetime : %v\n", err)
		return time.Time{}, err
	}
	return time.Unix(unixSec, 0), nil
}



func handle_query(contract *gateway.Contract) {
	var id string
	gotStringFromStdin(greetingIn + "Input person 'ID' you searched: ", &id)
	jsPers, err := contract.SubmitTransaction("QueryResident", id)
	if err != nil {
        	log.Printf("Failed to query person data: %v", err)
		return
        }

	var pers Resident
	err = json.Unmarshal(jsPers, &pers)
	if err != nil {
		log.Printf("Error parse json : %s", err)
		return
	}
	pushStringToStdout(greetingOut, pers)
}

func handle_create(contract *gateway.Contract) {
        var str string
	pushStringToStdout(greetingHelp, "Input person data in ':'-separated format")	
	pushStringToStdout(greetingHelp, "example: passpordId:firstname:secondname:city:address:mobile:status")
	gotStringFromStdin(greetingIn, &str)
	strSlice := strings.Split(str, ":") // id, fn, sn, c, addr, mob, st

        _, err := contract.SubmitTransaction("InsertResident", strSlice[0], strSlice[1], strSlice[2], strSlice[3], strSlice[4], strSlice[5], strSlice[6])
        if err != nil {
        	log.Printf("Failed to inserting person: %v\n", err)
		return
        }
	pushStringToStdout(greetingOut, "Successfully inserted")
}

func handle_update(contract *gateway.Contract) {
        var str string
	pushStringToStdout(greetingHelp, "Input updated person data in ':'-separated format without space")	
	pushStringToStdout(greetingHelp, "example: passpordId:firstname:secondname:city:address:mobile:status")
	gotStringFromStdin(greetingIn, &str)
	strSlice := strings.Split(str, ":") // id, fn, sn, c, addr, mob, st

        _, err := contract.SubmitTransaction("UpdateResident", strSlice[0], strSlice[1], strSlice[2], strSlice[3], strSlice[4], strSlice[5], strSlice[6])
        if err != nil {
        	log.Printf("Failed to update person: %v\n", err)
		return
	}
	pushStringToStdout(greetingOut, "Successfully updated")
}

func handle_history(contract *gateway.Contract) {
	var id string
	gotStringFromStdin(greetingIn + "Input person 'ID' you searched: ", &id)
	result, err := contract.SubmitTransaction("ReadHistoryResident", id)
	if err != nil {
		log.Printf("Failed to read person history %v\n", err)
		return
	}

	// parse got history
	var data []ResidentHistory
	err = json.Unmarshal(result, &data)
	if err != nil {
		log.Printf("Faild to parse data with person history: %v", err)
		return
	}

	lenData := len(data)
	for i, v := range data {
		numb := strconv.Itoa(lenData - i)
		txid := v.TxId[:5]
		dtime, err := parseStringDatetimeToUnixtime(v.Timestamp)
		if err == nil {
			dtimeStr := dtime.String()
			pushStringToStdout(greetingOut + numb + ". TxId=" + txid + "...|" + dtimeStr + ":", *v.Content)
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Invalid number of arguments (expected 2: channelName and contractName)")
	}
	contract := connectionInitialize("appUser", os.Args[1], os.Args[2])
    	
	var command string
	pushStringToStdout(greetingHelp, "Print ?? for help")
	gotStringFromStdin(greetingIn, &command)

	for {
		switch command {
		case "??":
			pushStringToStdout(greetingHelp, "Existed commands:")
			pushStringToStdout(greetingHelp, "\tcreate\t-\tcreate a new resident record in the blockchain")
			pushStringToStdout(greetingHelp, "\tread\t-\tread an existed resident record from the blockchain")
			pushStringToStdout(greetingHelp, "\tupdate\t-\tupdate an existed resident record in the blockchain")
			pushStringToStdout(greetingHelp, "\thistory\t-\tget a log of changes of existed resident from the blockchain")
			pushStringToStdout(greetingHelp, "\tquit\t-\texit from program")
			pushStringToStdout(greetingHelp, "")
		case "create":
			handle_create(contract)
		case "read": 
			handle_query(contract)
		case "update":
			handle_update(contract)
		case "history":
			handle_history(contract)
		case "quit":
			return
		default:
			pushStringToStdout(greetingHelp, "Unknown command")
		}
		gotStringFromStdin(greetingIn, &command)
	}
}
