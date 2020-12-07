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

var logger = shim.NewLogger("insurance")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Insurance struct {
	InsuranceNum         string `json:"insuranceNum"`
	OrderNum             string `json:"orderNum"`
	Shipper              string `json:"shipper"`
	Carrier              string `json:"carrier"`
	InsuranceCompany     string `json:"insuranceCompany"`
	GoodsName            string `json:“goodsName"`
	Num                  string `json:"num"`
	Weight               string `json:"weight"`
	Premium              string `json:"premium"`
	Status               string `json:"status"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response  {
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### insurance Invoke ###########")

	function, args := stub.GetFunctionAndParameters()
	if function == "createInsurance" {
		return t.createInsurance(stub, args)
	}
	if function == "flagInsuranceStatus" {
	        return t.flagInsuranceStatus(stub, args)
	}

	logger.Errorf("Unknown action, check the first argument, must be one of 'createInsurance', 'flagInsuranceStatus'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'createInsurance', 'flagInsuranceStatus'. But got: %v", args[0]))
}

func (t *SimpleChaincode) createInsurance(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke

	if len(args) != 10 {
		return shim.Error("Incorrect number of arguments. Expecting 10")
	}

	var insurance = Insurance{InsuranceNum: args[0],OrderNum: args[1],Shipper: args[2], Carrier: args[3],
	InsuranceCompany: args[4],GoodsName: args[5],Num: args[6], Weight: args[7],Premium: args[8],Status:args[9]}

        insuranceAsBytes, _ := json.Marshali(insurance)

	err := stub.PutState(args[0], insuranceAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

        return shim.Success(nil);
}

func (s *SimpleChaincode) flagInsuranceStatus(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	if len(args) != 2 {
	      return shim.Error("Incorrect number of arguments. Expecting 2")
        }
        insurancebytes, err := stub.GetState(args[0])
        if err != nil {
	        jsonResp := "{\"Error\":\"Failed to get state for " + args[0]+ "\"}"
		return shim.Error(jsonResp)
	}
	insurance := Insurance{}

	err = json.Unmarshal(insurancebytes,&insurance)

        if err != nil {
		return shim.Error(err.Error())
	}

        insurance.Status = args[1]

	insurancechangedAsBytes, _ := json.Marshal(insurance)


	err := stub.PutState(args[0], insurancechangedAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success()
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
