/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("CertStore")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Cert struct {
	ID            string `json:"id"`
	Data          string `json:"data"`
	IsDelete      string `json:"isDelete"`
}


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### cert Invoke ###########")

	function, args := stub.GetFunctionAndParameters()
	if function == "addInfo" {
		return t.addInfo(stub, args)
	}
	if function == "updateInfo" {
		return t.updateInfo(stub, args)
	}
	if function == "getInfo" {
		return t.getInfo(stub, args)
	}
	if function == "getHistoryInfo" {
		return t.getHistoryInfo(stub, args)
	}
	
	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0]))
}

func (t *SimpleChaincode) addInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	var cert = Cert{ID:args[0],Data:args[1]}

	certAsBytes, _ := json.Marshal(cert)

	// Write the state back to the ledger
	err := stub.PutState(args[0], certAsBytes)
	if err != nil {
		jsonResp := "{\"success\":\"false\"+\"details\":\"" + err.Error() + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"success\":\"true\"}"

	return shim.Success([]byte(jsonResp))
}


func (s *SimpleChaincode) updateInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	creditbytes, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}

	cert := Cert{}

	err = json.Unmarshal(creditbytes,cert)

	if err != nil {
		jsonResp := "{\"success\":\"false\"+\"details\":\"Failed to Unmarshal\"}"
		return shim.Error(jsonResp)
	}

	cert.Data = args[1]

	certAsBytes, _ := json.Marshal(cert)

	// Write the state back to the ledger
	err = stub.PutState(args[0], certAsBytes)
	if err != nil {
		jsonResp := "{\"success\":\"false\"+\"details\":\"" + err.Error() + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"success\":\"true\"}"

	return shim.Success([]byte(jsonResp))
}


func (s *SimpleChaincode) getInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	certbytes, err := stub.GetState(args[0])
	if err != nil {
		jsonResp := "{\"success\":\"false\"+\"datails\":\"Failed to get state for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}
	jsonResp := "{\"success\":\"true\" + \"details\":\"" + string(certbytes) + "\"}"

	return shim.Success([]byte(jsonResp))
}

func (s *SimpleChaincode) getHistoryInfo(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	historyIterator, err := stub.GetHistoryForKey(args[0])
	defer historyIterator.Close()
	if err != nil {
		jsonResp := "{\"success\":\"false\"+\"datails\":\"Failed to get history state for " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}
	fmt.Println("-----start historyIterator-----")
	str := ""
	for historyIterator.HasNext() {
		item, _ := historyIterator.Next()
		str = str + item.TxId + string(item.Value) + item.Timestamp.String() + strconv.FormatBool(item.IsDelete)
	}

	jsonResp := "{\"success\":\"true\" + \"details\":\"" + str + "\"}"

	return shim.Success([]byte(jsonResp))
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
