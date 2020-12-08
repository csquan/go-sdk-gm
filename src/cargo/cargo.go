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
	"fmt"
	"encoding/json"


	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("cargo")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Cargo struct {
	OrderNum	     string `json:"orderNum"`
	Shipper              string `json:"shipper"`
	Carrier              string `json:"carrier"`
	Buyer		     string `json:"buyer"`
	GoodsName            string `json:"goodsName"`
	Num                  string `json:"num"`
	Weight               string `json:"weight"`
	Price                string `json:"price"`
	Status               string `json:"status"`
}


func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response  {
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### cargo Invoke ###########")

	function, args := stub.GetFunctionAndParameters()
	if function == "createCargo" {
		return t.createCargo(stub, args)
	}
	if function == "flagCargoStatus" {
	        return t.flagCargoStatus(stub, args)
	}

	logger.Errorf("Unknown action, check the first argument, must be one of 'createCargo', 'flagCargoStatus'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'createCargo', 'flagCargoStatus'. But got: %v", args[0]))
}

func (t *SimpleChaincode) createCargo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke

	if len(args) != 9 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

	var cargo = Cargo{OrderNum: args[0],Shipper: args[1],Carrier: args[2], Buyer: args[3],
	GoodsName: args[4],Num: args[5],Weight: args[6], Price: args[7],Status: args[8]}

        cargoAsBytes, _ := json.Marshal(cargo)

	err := stub.PutState(args[0], cargoAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

        return shim.Success(nil);
}

func (s *SimpleChaincode) flagCargoStatus(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	if len(args) != 2 {
	      return shim.Error("Incorrect number of arguments. Expecting 2")
        }
        cargobytes, err := stub.GetState(args[0])
        if err != nil {
	        jsonResp := "{\"Error\":\"Failed to get state for " + args[0]+ "\"}"
		return shim.Error(jsonResp)
	}
	cargo := Cargo{}

	err = json.Unmarshal(cargobytes,&cargo)

        if err != nil {
		return shim.Error(err.Error())
	}

        cargo.Status = args[1]

	cargochangedAsBytes, _ := json.Marshal(cargo)

	err = stub.PutState(args[0], cargochangedAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
