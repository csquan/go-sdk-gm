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
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("note")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Note struct {
        NoteID               string `json:"noteID"`
        NoteHash             string `json:"noteHash"`
        CipperText           string `json:"cipperText"`
}



func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response  {
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### note Invoke ###########")

	function, args := stub.GetFunctionAndParameters()
	if function == "createNote" {
		return t.createNote(stub, args)
	}
	if function == "modifyNote" {
	        return t.modifyNote(stub, args)
	}
        if function == "getNote" {
                return t.getNote(stub, args)
        }

	logger.Errorf("Unknown action, check the first argument, must be one of 'createNote', 'modifyNote'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'creatNote', 'modifyNote','getNote'. But got: %v", args[0]))
}

func (t *SimpleChaincode) createNote(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	_, err := stub.GetState(args[0])
	if err == nil {
		jsonResp := "{\"Error\":\"OrderNum " + args[0]+ " have already exist\"}"
		return shim.Error(jsonResp)
	}

	var note = Note{NoteID: args[0],NoteHash: args[1],CipperText: args[2]}

        noteAsBytes, _ := json.Marshal(cargo)

	err = stub.PutState(args[0], noteAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

        return shim.Success(nil);
}

func (s *SimpleChaincode) modifyNote(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	if len(args) != 3 {
	      return shim.Error("Incorrect number of arguments. Expecting 3")
        }
        notebytes, err := stub.GetState(args[0])
        if err != nil {
	        jsonResp := "{\"Error\":\"Failed to get state for " + args[0]+ "\"}"
		return shim.Error(jsonResp)
	}

	note := Note{}

	err = json.Unmarshal(Notebytes,&note)

        if err != nil {
		return shim.Error(err.Error())
	}
	note.NoteHash = args[1]
        note.CipperText = args[2]

	notechangedAsBytes, _ := json.Marshal(note)

	err = stub.PutState(args[0], notechangedAsBytes)

	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (s *SimpleChaincode) getNote(stub shim.ChaincodeStubInterface,args []string) pb.Response {

        if len(args) != 1 {
              return shim.Error("Incorrect number of arguments. Expecting 1")
        }
        notebytes, err := stub.GetState(args[0])
        if err != nil {
                jsonResp := "{\"Error\":\"Failed to get state for " + args[0]+ "\"}"
                return shim.Error(jsonResp)
        }
        jsonResp := "{\"\"note:\"" + string(notebytes) + "\"}"

        logger.Info("jsonResp")
        logger.Info(jsonResp)

        return shim.Success([]byte(jsonResp))
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
