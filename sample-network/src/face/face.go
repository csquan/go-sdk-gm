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
	"bytes"
	"fmt"
	"time"
	"encoding/json"
        "strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("face")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type secoo_face struct {
	Secoo_userid_hash    string `json:"secoo_userid"`
	Auth_id              string `json:"auth_id"`
	Third_requestid      string `json:"third_request_id"`
	Third_resultid       string `json:"third_result_id"`
} 

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response  {
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### face Invoke ###########")

	function, args := stub.GetFunctionAndParameters()
	if function == "createFace" {
		return t.createFace(stub, args)
	}
	if function == "queryOneFace" {
	        return t.queryOneFace(stub, args)
	}
	if function == "queryAllFace" {
	        return t.queryAllFace(stub, args)
	}

	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0]))
}

func (t *SimpleChaincode) createFace(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}


	var face = secoo_face{Secoo_userid_hash: args[0],Auth_id: args[1],Third_requestid: args[2], Third_resultid: args[3]}

        faceAsBytes, _ := json.Marshal(face)

	//make key = args[0] + time
	key := args[0] + strconv.FormatInt(time.Now().Unix(),10)
	logger.Info("########### face key:%s",key)
	// Write the state back to the ledger
	err := stub.PutState(key, faceAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

        return shim.Success(nil);
}

func (s *SimpleChaincode) queryOneFace(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	if len(args) != 1 {
	      return shim.Error("Incorrect number of arguments. Expecting 1")
        }
        facebytes, err := stub.GetState(args[0])
        if err != nil {
	        jsonResp := "{\"Error\":\"Failed to get state for " + args[0]+ "\"}"
		return shim.Error(jsonResp)
	}
	jsonResp := "{\"face\":\"" + string(facebytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)

	return shim.Success(facebytes)
}

func (s *SimpleChaincode) queryAllFace(stub shim.ChaincodeStubInterface,args []string) pb.Response {
	if len(args) != 2 {
                 return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	startKey := args[0]
        endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		 return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
        buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		 queryResponse, err := resultsIterator.Next()
		 if err != nil {
			return shim.Error(err.Error())
		 }
																									        if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")																				        }
		 buffer.WriteString("{\"Key\":")
		 buffer.WriteString("\"")
		 buffer.WriteString(queryResponse.Key)
		 buffer.WriteString("\"")
		 buffer.WriteString(", \"Record\":")
		 buffer.WriteString(string(queryResponse.Value))
		 buffer.WriteString("}")
		 bArrayMemberAlreadyWritten = true
	 }
	 buffer.WriteString("]")
	 fmt.Printf("- queryAllCars:\n%s\n", buffer.String())
         return shim.Success(buffer.Bytes())
 }

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
