package main

import (
	"encoding/json"
	"fmt"
	"bytes"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing the population in the country
type SmartContract struct {
	contractapi.Contract
}

// Resident describes basic details
type Resident struct {
	Firstname   	string `json:"firstname"`
	Secondname  	string `json:"secondname"`
	City 		string `json:"city"`
	Address		string `json:"address"`
	Mobile		string `json:"mobile"`
	Status		string `json:"status"`
}

type ResidentHistory struct {
	Content		*Resident	`json:"content"`
	TxId		string	 	`json:"txid"`
	Timestamp	string		`json:"ts"`
}


// InitLedger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	fmt.Println(" ==== Init contract ====")
	return nil
}

// QueryResident returns struct Resident if exists
func (s *SmartContract) QueryResident(ctx contractapi.TransactionContextInterface, pass string) (*Resident, error) {
	personAsBytes, err := ctx.GetStub().GetState(pass)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if personAsBytes == nil {
		return nil, fmt.Errorf("Person with pass = %s does not exist.", pass)
	}
	person := new(Resident)
	// we may not check, because the json stored in the ledger
	_ = json.Unmarshal(personAsBytes, person)
	return person, nil
}


// Insert resident to ledger
func (s *SmartContract) InsertResident(ctx contractapi.TransactionContextInterface, pass string, fn string, sn string, city string, addr string, mob string, stat string) (error) {
	personAsBytes, err := ctx.GetStub().GetState(pass)
	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if personAsBytes != nil {
		return fmt.Errorf("Resident already exists with such passID = %s", pass)
	}
	person := Resident {
		Firstname: 	fn,
		Secondname:	sn,
		City:		city,
		Address:	addr,
		Mobile:		mob,
		Status:		stat,
	}
	personAsBytes, _ = json.Marshal(person)
	return ctx.GetStub().PutState(pass, personAsBytes)
}


// UpdateResident updates the fields of resident with given id in world state
func (s *SmartContract) UpdateResident(ctx contractapi.TransactionContextInterface, 
		pass		string,
		firstname 	string, 
		secondname 	string, 
		city 		string,
		address 	string,
		mobile		string,
		status		string,
	) error {
	storedPersonAsBytes, err := ctx.GetStub().GetState(pass)
	if err != nil {
		return fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if storedPersonAsBytes == nil {
		return fmt.Errorf("Person with pass = %s does not exist.", pass)
	}
	person := Resident{
		Firstname:  	firstname,
		Secondname:  	secondname,
		City: 		city,
		Address:  	address,
		Mobile:  	mobile,
		Status:  	status,
	}
	loadedPersonAsBytes, _ := json.Marshal(person)
	if bytes.Compare(storedPersonAsBytes, loadedPersonAsBytes) == 0 {
		return fmt.Errorf("Failed to update resident in the world state. Resident with pass = %s has same fields", pass)
	}
	return ctx.GetStub().PutState(pass, loadedPersonAsBytes)
}

// Get credential changing log
func (s *SmartContract) ReadHistoryResident(ctx contractapi.TransactionContextInterface, pass string) ([]ResidentHistory, error) {
	storedPersonAsBytes, err := ctx.GetStub().GetState(pass)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	if storedPersonAsBytes == nil {
		return nil, fmt.Errorf("Person with pass = %s does not exist.", pass)
	}
	// get history
	log, err := ctx.GetStub().GetHistoryForKey(pass)
	if err != nil {
		return nil, fmt.Errorf("Failed to get history by passID = %s. %s", pass, err.Error())
	}
	fmt.Println(log, err)
	residentHistory := make([]ResidentHistory, 0)
	for log.HasNext() {
		epoch, err := log.Next()
		fmt.Println(epoch)
		if err != nil {
			// internal error while reading next elem
			return nil, fmt.Errorf("Can't read all elements in history log passID = %s. %s", pass, err.Error())
		}
		item := new(ResidentHistory)
		_ = json.Unmarshal(epoch.GetValue(), &item.Content)
		item.Timestamp = epoch.GetTimestamp().String()
		item.TxId = epoch.TxId
		residentHistory = append(residentHistory, *item)
	}
	return residentHistory, nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
