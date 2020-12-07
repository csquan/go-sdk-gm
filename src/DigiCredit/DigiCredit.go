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
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("DigiCredit")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type DigiCredit_Root struct {
	Sign              string `json:"sign"`
	DigiCreditId      string `json:"digiCreditId"`
	CompanyId         string `json:"companyId"`
	OpenAmount        string `json:"openAmount"`
	LastAmount        string `json:"lastAmount"`
	ChildDigiCreditId string `json:"childDigiCreditId"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
	State             string `json:"state"`
	Operator          string `json:"operator"`
	OperTime          string `json:"operTime"`
	Authorizer        string `json:"authorizer"`
	AuthorizeTime     string `json:"authorizeTime"`
	Remark            string `json:"remark"`
}

type DigiCredit_Child struct {
	Sign               string `json:"sign"`
	DigiCreditId       string `json:"digiCreditId"`
	RootDigiCreditId   string `json:"rootDigiCreditId"`
	ParentDigiCreditId string `json:"parentDigiCreditId"`
	ChildDigiCreditId  string `json:"childDigiCreditId"`
	FromCompanyId      string `json:"fromcompanyId"`
	ToCompanyId        string `json:"tocompanyId"`
	StartDate          string `json:"startDate"`
	EndDate            string `json:"endDate"`
	Amount             string `json:"amount"`
	State              string `json:"state"`
	Remark             string `json:"remark"`
}

type CommpanyInfo struct {
	CompanyId    string `json:"companyId"`
	CompanyType  string `json:"companyType"`
	CompanyCert  string `json:"companyCert"`
	CompanyState string `json:"companyState"`
	TotalAmount  string `json:"totalAmount"`
	Remark       string `json:"remark"`
	DigiCredit   string `json:"digiCredit"`
	DigiCreditID string `json:"digiCreditID"`
	Amount       string `json:"amount"`
}

func (s *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Transaction makes payment of X units from A to B
func (s *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### digiCredit Invoke ###########")

	function, args := stub.GetFunctionAndParameters()
	if function == "companyRegist" {
		return s.companyRegist(stub, args)
	}
	if function == "digiCreditRootIssue" {
		return s.digiCreditRootIssue(stub, args)
	}
	if function == "digiCreditIssue" {
		return s.digiCreditIssue(stub, args)
	}
	if function == "digiCreditTrans" {
		return s.digiCreditTrans(stub, args)
	}
	if function == "editCompanyCert" {
		return s.editCompanyCert(stub, args)
	}
	if function == "editState" {
		return s.editState(stub, args)
	}
	if function == "editEndTime" {
		return s.editEndTime(stub, args)
	}
	if function == "editOpenAmount" {
		return s.editOpenAmount(stub, args)
	}
	if function == "queryDigiCreditRoot" {
		return s.queryDigiCreditRoot(stub, args)
	}
	if function == "queryDigiCredit" {
		return s.queryDigiCredit(stub, args)
	}
	if function == "queryCompany" {
		return s.queryCompany(stub, args)
	}
	if function == "queryCompany" {
		return s.queryCompany(stub, args)
	}
	if function == "queryTransHistory" {
		return s.queryTransHistory(stub, args)
	}

	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0]))
}

func (s *SimpleChaincode) companyRegist(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var company = CommpanyInfo{CompanyId: args[0], CompanyType: args[1], CompanyCert: args[2], Remark: args[3]}

	companyAsBytes, _ := json.Marshal(company)

	logger.Info("companyinfo key:%s", args[0])
	// Write the state back to the ledger
	err := stub.PutState(args[0], companyAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SimpleChaincode) digiCreditRootIssue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke

	if len(args) != 11 {
		return shim.Error("Incorrect number of arguments. Expecting 11")
	}
	var rootcredit = DigiCredit_Root{Sign: args[0], CompanyId: args[1], DigiCreditId: args[2], OpenAmount: args[3], LastAmount: args[3],
		StartDate: args[4], EndDate: args[5], Operator: args[6], OperTime: args[7], Authorizer: args[8], AuthorizeTime: args[9], Remark: args[10],ChildDigiCreditId:"[]",}

	rootcreditAsBytes, _ := json.Marshal(rootcredit)

	logger.Info("digiCreditRootIssue key:%s", args[2])
	// Write the state back to the ledger
	err := stub.PutState(args[2], rootcreditAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SimpleChaincode) digiCreditIssue(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// must be an invoke

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	//1.generate child credit
	var childcredit = DigiCredit_Child{Sign: args[0], FromCompanyId: args[1], ToCompanyId: args[2],
		ParentDigiCreditId: args[3], DigiCreditId: args[4], Amount: args[5], Remark: args[6],ChildDigiCreditId:"[]",}

	childcreditAsBytes, _ := json.Marshal(childcredit)

	logger.Info("begin Issue child digiCredit key:%s", args[4])
	// Write the state back to the ledger
	err := stub.PutState(args[4], childcreditAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	//2.add ChildDigiCreditId to ParentDigiCredit's childarray
	logger.Info("begin to get parent credit base on key:%s", args[3])
	parentcreditbytes, err := stub.GetState(args[3])
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info(string(parentcreditbytes))

	if len(string(parentcreditbytes))== 0{
		return shim.Error("parent credit is null")
	}

	index := strings.Index(string(parentcreditbytes),"parent")

	fmt.Println("++++++parent++++++:")
	fmt.Println(string(parentcreditbytes))

	fmt.Println("++++++index++++++:")
	fmt.Println(index)

	if index >=0{
		return shim.Error("args[3] is childcredit no,not rootcredit no")		
	}else{
		rootcredit := DigiCredit_Root{}

		err = json.Unmarshal(parentcreditbytes,&rootcredit)

		pos := strings.Index(rootcredit.ChildDigiCreditId,"]")
	
		if pos >= 0{
			oldchilds := rootcredit.ChildDigiCreditId[0:pos]
			
			logger.Info("oldchilds")
			logger.Info(oldchilds)
		
			rootcredit.ChildDigiCreditId = oldchilds + args[4] + ",]"

			lastamount ,err:= strconv.Atoi(rootcredit.LastAmount)
			if err != nil {
				return shim.Error(err.Error())
			}

			amount ,err:= strconv.Atoi(args[5])
			if err != nil {
				return shim.Error(err.Error())
			}

			newamount := lastamount - amount

			if lastamount <0{
				return shim.Error("amount is big than available lastamount")
			}

			rootcredit.LastAmount = strconv.Itoa(newamount)

			logger.Info("newchilds")
			logger.Info(rootcredit.ChildDigiCreditId)

			//3.store rootcredit
			newrootcreditAsBytes, _ := json.Marshal(rootcredit)

			logger.Info("begin to put root credit back base on key:%s", args[3])
			// Write the state back to the ledger
			err = stub.PutState(args[3], newrootcreditAsBytes)
			if err != nil {
				return shim.Error(err.Error())
			}
		}else{
			return shim.Error("rootcredit's ChildDigiCreditId has no ]")	
		}
	}
	
	return shim.Success(nil)
}

func (s *SimpleChaincode) digiCreditTrans(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	creditbytes, err := stub.GetState(args[3])
	if err != nil {
		return shim.Error(err.Error())
	}

	if creditbytes == nil{
		return shim.Error(err.Error())
	}

	index := strings.Index(string(creditbytes),"parent")

    if index >=0{		
		oldcredit := DigiCredit_Child{}

		err = json.Unmarshal(creditbytes,&oldcredit)
		if err != nil {
			return shim.Error(err.Error())
		}

		logger.Info("oldcredit")
		logger.Info(oldcredit)


		newamount,err:= strconv.Atoi(args[5])

		if err != nil {
			return shim.Error(err.Error())
		}
		oldamount ,err:= strconv.Atoi(oldcredit.Amount)
		fmt.Println("oldcredit.Amount:")
		fmt.Println(oldcredit.Amount)
		if err != nil {
			return shim.Error(err.Error())
		}

		if newamount > oldamount {
			return shim.Error(err.Error())
		}
		//1.modify oldcredit 
		oldcredit.Amount = strconv.Itoa(oldamount - newamount)

		///add this sub credit to parent credit 
		pos := strings.Index(oldcredit.ChildDigiCreditId,"]")
	
		if pos >= 0{
			oldchilds := oldcredit.ChildDigiCreditId[0:pos]
		
			oldcredit.ChildDigiCreditId = oldchilds + args[4] + ",]"
		}else{
			return shim.Error("parent credit have no ]")
		}
		
		//2.store oldcredit
		oldcreditAsBytes, _ := json.Marshal(oldcredit)

		//old_credit_key = creditid
		oldCreditKey := args[3]
		logger.Info("oldCreditKey key:%s", oldCreditKey)
		// Write the state back to the ledger
		err = stub.PutState(oldCreditKey, oldcreditAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	}else{
		oldcredit := DigiCredit_Root{}

		err = json.Unmarshal(creditbytes,&oldcredit)
		if err != nil {
			return shim.Error(err.Error())
		}

		logger.Info("oldcredit")
		logger.Info(oldcredit)

		newamount,err:= strconv.Atoi(args[5])

		if err != nil {
			return shim.Error(err.Error())
		}
		oldamount ,err:= strconv.Atoi(oldcredit.LastAmount)
		fmt.Println("++++++oldcredit.Amount++++++:")
		fmt.Println(oldcredit.LastAmount)
		if err != nil {
			return shim.Error(err.Error())
		}

		if newamount > oldamount {
			return shim.Error(err.Error())
		}
		//1.modify oldcredit 
		oldcredit.LastAmount = strconv.Itoa(oldamount - newamount)
		
		///add this sub credit to parent credit 
		pos := strings.Index(oldcredit.ChildDigiCreditId,"]")
	
		if pos >= 0{
			oldchilds := oldcredit.ChildDigiCreditId[0:pos]	
			oldcredit.ChildDigiCreditId = oldchilds + args[4] + ",]"
		}else{
			return shim.Error("parent credit have no ]")
		}
		//2.store oldcredit
		oldcreditAsBytes, _ := json.Marshal(oldcredit)

		//old_credit_key = creditid
		oldCreditKey := args[3]
		logger.Info("oldCreditKey key:%s", oldCreditKey)
		// Write the state back to the ledger
		err = stub.PutState(oldCreditKey, oldcreditAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	//3.create sub credit
	var newcredit = DigiCredit_Child{Sign: args[0], FromCompanyId: args[1], ToCompanyId: args[2],
		ParentDigiCreditId: args[3],   DigiCreditId:args[4], Amount: args[5]}
		
	//4.store newcredit
	newcreditAsBytes, _ := json.Marshal(newcredit)

	newCreditKey := args[4]
	logger.Info("########### newcreditAsBytes key:%s", newcreditAsBytes)
	// Write the state back to the ledger
	err = stub.PutState(newCreditKey, newcreditAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


func (s *SimpleChaincode) editCompanyCert(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	companybytes, err := stub.GetState(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}

	companyinfo := CommpanyInfo{}

	err = json.Unmarshal(companybytes,&companyinfo)

	if err != nil {
		return shim.Error(err.Error())
	}

	companyinfo.CompanyCert = args[2]

	//3.store newcompanyinfo
	newcompanyinfoAsBytes, _ := json.Marshal(companyinfo)

	// Write the state back to the ledger
	err = stub.PutState(args[1], newcompanyinfoAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SimpleChaincode) editState(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	creditbytes, err := stub.GetState(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}

	credit := DigiCredit_Root{}

	err = json.Unmarshal(creditbytes,&credit)

	if err != nil {
		return shim.Error(err.Error())
	}

	logger.Info("before change credit's state")
	logger.Info(credit)
	credit.Sign = args[0]
	credit.CompanyId = args[1]
	credit.State = args[3]
    logger.Info("after change credit's state")
	logger.Info(credit)

	//3.store credit
	newcreditAsBytes, _ := json.Marshal(credit)

	// Write the state back to the ledger
	err = stub.PutState(args[2], newcreditAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SimpleChaincode) editEndTime(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	creditbytes, err := stub.GetState(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}

	credit := DigiCredit_Root{}

	err = json.Unmarshal(creditbytes,&credit)

	if err != nil {
		return shim.Error(err.Error())
	}

	credit.Sign = args[0]
	credit.CompanyId = args[1]
	credit.EndDate = args[3]

	//3.store credit
	creditAsBytes, _ := json.Marshal(credit)

	// Write the state back to the ledger
	err = stub.PutState(args[2], creditAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

func (s *SimpleChaincode) editOpenAmount(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}
	creditbytes, err := stub.GetState(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
    //checkout root no
	index := strings.Index(string(creditbytes),"parent")

	if index >=0{
		return shim.Error("args[2] is childcredit no,not rootcredit no")		
	}

	credit := DigiCredit_Root{}

	err = json.Unmarshal(creditbytes,&credit)

	if err != nil {
		return shim.Error(err.Error())
	}
	//checkout args[3] is legal
	amount ,err:= strconv.Atoi(args[3])
	if err != nil {
		return shim.Error(err.Error())
	}

	lastamount ,err:= strconv.Atoi(credit.LastAmount)
	if err != nil {
		return shim.Error(err.Error())
	}

	openamount ,err:= strconv.Atoi(credit.OpenAmount)
	if err != nil {
		return shim.Error(err.Error())
	}

	if openamount - lastamount > amount{
		return shim.Error("newamount should be bigger than openamount or filled patten openamount - lastamount <= amount ")
	}

	credit.Sign = args[0]
	credit.CompanyId = args[1]
	credit.OpenAmount = args[3]

	//3.store credit
	creditAsBytes, _ := json.Marshal(credit)

	// Write the state back to the ledger
	err = stub.PutState(args[2], creditAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


func (s *SimpleChaincode) queryDigiCreditRoot(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	rootcreditbytes, err := stub.GetState(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
	jsonResp := "{\"rootcredit\":\"" + string(rootcreditbytes) + "\"}"

	logger.Info("jsonResp")
	logger.Info(jsonResp)

	return shim.Success([]byte(jsonResp))
}

func (s *SimpleChaincode) queryDigiCredit(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	rootcreditbytes, err := stub.GetState(args[2])
	if err != nil {
		return shim.Error(err.Error())
	}
	jsonResp := "{\"credit\":\"" + string(rootcreditbytes) + "\"}"

	return shim.Success([]byte(jsonResp))
}

func (s *SimpleChaincode) queryCompany(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	companyinfobytes, err := stub.GetState(args[1])
	if err != nil {
		return shim.Error(err.Error())
	}
	jsonResp := "{\"companyinfo\":\"" + string(companyinfobytes) + "\"}"

	return shim.Success([]byte(jsonResp))
}

func (s *SimpleChaincode) queryTransHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	digicreditbytes, err := stub.GetState(args[3])
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("get credit:")
	fmt.Println(string(digicreditbytes))

	childIDs := ""
	ret := ""

	if args[2] == "0" {

		rootcredit := DigiCredit_Root{}

		err = json.Unmarshal(digicreditbytes,&rootcredit)

		if err != nil {
			return shim.Error(err.Error())
		}

		fmt.Println("get root bytes:")
		fmt.Println(string(digicreditbytes))

		index := strings.Index(string(digicreditbytes),"parent")
	
		if index >=0{
			return shim.Error("args[3] is childcredit no,not rootcredit no and your input type is root,two should same")	
		}	

		ret = "digiCreditRoot:" + string(digicreditbytes)

		childIDs = rootcredit.ChildDigiCreditId
    }else if args[2] == "1"{
		childcredit := DigiCredit_Child{}

		err = json.Unmarshal(digicreditbytes,&childcredit)

		if err != nil {
			return shim.Error(err.Error())
		}

		fmt.Println("get child bytes:")
		fmt.Println(string(digicreditbytes))

		index := strings.Index(string(digicreditbytes),"parent")
	
		if index < 0{
			return shim.Error("args[3] is rootcredit no,not childcredit no and your input type is child,two should same")	
		}	
		ret = "digiCredit:" + string(digicreditbytes)
		childIDs = childcredit.ChildDigiCreditId
		
	}else{
		return shim.Error(" args[2]  type is wrong not 1 or 2")
	}
	if childIDs == ""{
		return shim.Error(err.Error())
	}
	
	//childIDs rid off []
    childIDs = strings.Trim(childIDs,"[,]")
	IDs := strings.Split(childIDs, ",")
	fmt.Println("parse child ids:")
	fmt.Println(IDs)
	
	ret = ret + ",childDigiCredit:["

	for _ ,childID := range IDs {
		if childID != ""{
			childidbytes, err := stub.GetState(childID)
			
			if err != nil {
				return shim.Error(err.Error())
			}
			fmt.Println("get child bytes:")
			fmt.Println(string(childidbytes))

			ret = ret + string(childidbytes) +","
		  } 
	}
	fmt.Println("++++++before trim++++++:")
	fmt.Println(ret)
	
	ret = strings.TrimRight(ret,",")
	fmt.Println("++++++after trim+++++++:")
	fmt.Println(ret)

	ret = ret + "]"
	fmt.Println("++++++end+++++++:")
	fmt.Println(ret)

	return shim.Success([]byte(ret))
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
