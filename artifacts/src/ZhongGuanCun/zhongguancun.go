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

//UploadCert -- upload cert
type UploadCert struct {
	OutOrderID			string `json:"outOrderID"`
	CredentialType 		string `json:"credentialType"`
	PicType 			string `json:"picType"`
	CredentialNo 		string `json:"credentialType"`
	FileType 			string `json:"fileType"`
	PicBase64 			string `json:"picBase64"`
}


//PersonalRegister -- Personal Register
type PersonalRegister struct {
	OutUserID 			string `json:"outUserID"`
	CredentialType 		string `json:"credentialType"`
	CredentialNo 		string `json:"credentialNo"`
	CredentialExpiry 	string `json:"credentialExpiry"`
	CredentialExpiryEnd string `json:"credentialExpiryEnd"`
	CredentialPicIDUp 	string `json:"credentialPicIDUp"`
	CredentialPicIDDown string `json:"credentialPicIDDown"`
	AccountName 		string `json:"accountName"`
	AddressCode 		string `json:"addressCode"`
	Address 			string `json:"address"`
	WorkAddress 		string `json:"workAddress"`
	Phone 				string `json:"phone"`
	BankName 			string `json:"bankName"`
	BranchBankName 		string `json:"branchBankName"`
	BankCode 			string `json:"bankCode"`
	BankCardNo 			string `json:"bankCardNo"`
	BindBankPhone 		string `json:"bindBankPhone"`
}

//QueryPerson -- query Personal
type QueryPerson struct {
	OutUserID 				string `json:"outUserID"`
}

//CardTie -- Card Tie
type CardTie struct {
	BindCard				string `json:"bindCard"`
	CertNo					string `json:"CertNo"`
	Remark					string `json:"remark"`
}

//CardUnTie -- Card UnTie
type CardUnTie struct {
	UnBindCard 				string `json:"unBindCard"`
}

//AccountWithdraw -- Account Withdraw
type AccountWithdraw struct {
	OutPlatformID 			string `json:"outPlatformID"`
	OrderID 				string `json:"orderID"`
	OutUserID 				string `json:"outUserID"`
	OutAccountType 			string `json:"outAccountType"`
	BankCardNo 				string `json:"bankCardNo"`
	IsShowAll 				string `json:"isShowAll"`
	Amount 					string `json:"amount"`
	Fee 					string `json:"fee"`
	Remark 					string `json:"remark"`
	OrderInfo 				string `json:"orderInfo"`
}

//AccountTransfer -- Account Transfer
type AccountTransfer struct {
	OutPlatformID 		string `json:"outPlatformID"`
	OrderID 			string `json:"orderID"`
	GoodsName 			string `json:"goodsName"`
	OrderAmt 			string `json:"orderAmt"`
	IsBackTrack 		string `json:"isBackTrack"`
	PayerList 			string `json:"payerList"`
	OrderSeqn 			string `json:"orderSeqn"`
	PayUserID 			string `json:"payUserID"`
	PayerAccSeqn 		string `json:"payerAccSeqn"`
	PayerAmt 			string `json:"payerAmt"`
	PaysderTradeType 	string `json:"payerTradeType"`
	PayerRemark 		string `json:"payerRemark"`
	OrderID 			string `json:"orderID"`
	OrderInfo 			string `json:"orderInfo"`
	PayerList 			string `json:"payerList"`
	PayeeList 			string `json:"payeeList"`
	OrderSeqn 			string `json:"orderSeqn"`
	ReceiveUserID 		string `json:"receiveUserID"`
	PayeeAccSeqn 		string `json:"payeeAccSeqn"`
	PayeeAmt 			string `json:"payeeAmt"`
	PayeeTradeType 		string `json:"payeeTradeType"`
	PayeeRemark 		string `json:"payeeRemark"`
	PayeeList 			string `json:"payeeList"`
}

//QueryTransfer -- Quers Transfer
type QueryTransfer struct {
	OutPlatformID 		string `json:"outPlatformID"`
	OrigOrderID 		string `json:"origOrderID"`
}

//QueryPublicAccount -- Query Public Account
type QueryPublicAccount struct {
	OutPlatformID   string `json:"outPlatformID"`
	PayerBankCard   string `json:"payerBankCard"`
	Currency  		string `json:"currency"`
}

//QueryTransHistory - -Query Trans History
type QueryTransHistory struct {
	OutPlatformID 	string `json:"outPlatformID"`
	PrdtNo 			string `json:"prdtNo"`
    PriThirdTradeNo string `json:"priThirdTradeNo"`
	OutUserID 		string `json:"outUserID"`
	OutAccSeqn 		string `json:"outAccSeqn"`
	TradeType 		string `json:"tradeType"`
	AddFlag 		string `json:"addFlag"`
	BegDate 		string `json:"begDate"`
	EndDate 		string `json:"endDate"`
	CurrPage 		string `json:"currPage"`
	EachPageNum 	string `json:"eachPageNum"`
	SortType 		string `json:"sortType"`
}

//OpenEnterpriseAccount -- Open Enterprise Account
type OpenEnterpriseAccount struct {
	OutUserID    			string `json:"outUserID"`
	CredentialType 			string `json:"credentialType"`
	CredentialNo 			string `json:"credentialNo"`
	CredentialPicID 		string `json:"credentialPicID"`
	CredentialExpiry 		string `json:"credentialExpiry"`
	CredentialExpiryEnd 	string `json:"credentialExpiryEnd"`
	AccountName 			string `json:"accountName"`
	Phone 					string `json:"phone"`
	State 					string `json:"state"`
	ActualController 		string `json:"actualController"`
	BankName 				string `json:"bankName"`
	BranchBankName 			string `json:"branchBankName"`
	BankCode 				string `json:"bankCode"`
	BankCardNo 				string `json:"bankCardNo"`
	BindBankPhone 			string `json:"bindBankPhone"`
}

//SmallAmountAuth -- Small Amount Auth
type SmallAmountAuth struct {
	OutUserID    			string `json:"outUserID"`
	transCurrency  			string `json:"transCurrency"`
	maxAmount 				string `json:"maxAmount"`
	payAccountNo 			string `json:"payAccountNo"`
	payAccountName 			string `json:"payAccountName"`
	payType 				string `json:"payType"`
	payeeNo 				string `json:"payeeNo"`
	payeeAccountName 		string `json:"payeeAccountName"`
	payeeType 				string `json:"payeeType"`
	payeeNoType 			string `json:"payeeNoType"`
	depositBankNo 			string `json:"depositBankNo"`
	businessType 			string `json:"businessType"`
}

//QueryTransRecord -- Query Trans Record
type QueryEnterpriseUser struct {
	OutUserID    			string `json:"outUserID"`
}

//RegisterOrModifyResponse -- Register Or Modify Response
type MultiAccountFundRevoke struct {
	OutPlatformID 	string `json:"OutPlatformID"`
	PrdtNo 			string `json:"prdtNo"`
	OrderID 		string `json:"orderId"`
	OutUserID  		string `json:"outUserId"`
	IDType 			string `json:"iDType"`
	IDNo 			string `json:"iDNo"`
	AccNo 			string `json:"accNo"`
	OperType 		string `json:"operType"`
	BankCardNo 		string `json:"bankCardNo"`
	BankCardName 	string `json:"bankCardName"`
	CardBank 		string `json:"cardBank"`
	List1 			string `json:"list1"`
	OutAccountType 	string `json:"outAccountType"`
	Amount 			string `json:"amount"`
	List2 			string `json:"list2"`
	Fee 			string `json:"fee"`
	ReturnURL 		string `json:"returnURL"`
	NotifyURL 		string `json:"notifyURL"`
	OrderInfo 		string `json:"orderInfo"`
}

//WithdrawResponse -- Withdraw Response
type QueryMultiFundResponse struct {
	OutOrderID 				string `json:"outOrderId"`
	Status 					string `json:"Status"`
	Amount 					string `json:"amount"`
	OutUserID 				string `json:"outUserID"`
	Memo 					string `json:"memo"`
	Remark 					string `json:"remark"`
}

//RechangeResponse -- Rechange Response
type RechangeResponse struct {
	OutOrderID 				string `json:"outOrderId"`
	BusinessDate 			string `json:"businessDate"`
	BankAccountNo 			string `json:"bankAccountNo"`
 	VAcctNo 				string `json:"vAcctNo"`
	Amount 					string `json:"amount"`
	FeeType 				string `json:"feeType"`
	FeeAmount 				string `json:"feeAmount"`
	Status 					string `json:"Status"`
	TransType 				string `json:"transType"`
}

//RechargeReconciliationFile -- Recharge Reconciliation File
type RechargeReconciliationFile struct {
	OutOrderID 				string `json:"outOrderId"`
	BusinessDate 			string `json:"businessDate"`
	BankAccountNo 			string `json:"bankAccountNo"`
 	VAcctNo 				string `json:"vAcctNo"`
	Amount 					string `json:"amount"`
	FeeType 				string `json:"feeType"`
	FeeAmount 				string `json:"feeAmount"`
	Status 					string `json:"Status"`
	TransType 				string `json:"transType"`
}

//WithdrawReconciliationFile -- Withdraw Reconciliation File
type WithdrawReconciliationFile struct {
	OutOrderID 				string `json:"outOrderId"`
	BusinessDate 			string `json:"businessDate"`
	BankAccountNo 			string `json:"bankAccountNo"`
 	VAcctNo 				string `json:"vAcctNo"`
	Amount 					string `json:"amount"`
	FeeType 				string `json:"feeType"`
	FeeAmount 				string `json:"feeAmount"`
	Status 					string `json:"Status"`
	TransType 				string `json:"transType"`
}

//TransferReconciliationFile -- Transfer Reconciliation File
type TransferReconciliationFile struct {
	OutOrderID 				string `json:"outOrderID"`
	BusinessDate 			string `json:"businessDate"`
	PayuserAccountNo 		string `json:"payuserAccountNo"`
 	ReceiveuserAccountNo 	string `json:"receiveuserAccountNo"`
	Amount 					string `json:"amount"`
	FeeType 				string `json:"feeType"`
	FeeAmount 				string `json:"feeAmount"`
	Status 					string `json:"Status"`
}

//TransRecord -- Trans Record
type TransRecord struct {
	OrigOutOrderID 			string `json:"origOutOrderID"`
	OutUserID 				string `json:"outUserID"`
	TransType 		        string `json:"transType"`
}

//ChannelRefund -- Channel Refund
type ChannelRefund struct {
	OutOrderID 			string `json:"origOutOrderID"`
	OrigConfirmOrderID	string `json:"outUserID"`
	OutPayID 		    string `json:"outPayID"`
	Amount              string `json:"amount"`
	Reason 				string `json:"reason"`
	RefundType 			string `json:"refundType"`
	PayChannelCode 		string `json:"payChannelCode"`
	CallbackURL 		string `json:"callbackURL"`
	Remark 				string `json:"remark"`
}

//ChannelRefundResponse -- Channel Refund Response
type ChannelRefundResponse struct {
	OutOrderID 			string `json:"origOutOrderID"`
	Status 				string `json:"status"`
	OrigConfirmOrderID	string `json:"origConfirmOrderID"`
	Amount              string `json:"amount"`
	RefundType          string `json:"refundType"`
	Remark 				string `json:"remark"`
}

//ConsumptionRevoke -- Consumption Revoke
type ConsumptionRevoke struct {
	OutOrderID 			string `json:"outOrderID"`
	OrigOutOrderID 		string `json:"origOutOrderID"`
	OutPayID 			string `json:"outPayID"`
	Amount 				string `json:"amount"`
	Reason 				string `json:"reason"`
	RefundType 			string `json:"refundType"`
	OutUserID 			string `json:"outUserID"`
	PayChannelCode 		string `json:"payChannelCode"`
	CallbackURL 		string `json:"callbackURL"`
	Remark 				string `json:"remark"`
}

//RevokeResponse -- Revoke Response
type RevokeResponse struct {
	OutOrderID 			string `json:"outOrderID"`
	Status 				string `json:"status"`
	OrigOutOrderID		string `json:"origOutOrderID"`
	Amount              string `json:"amount"`
	RefundType          string `json:"refundType"`
	Remark 				string `json:"remark"`
}

//QueryBatchOrder -- Query Batch Order
type QueryBatchOrder struct {
	OutBatchID 			string `json:"outBatchID"`
	TransType 			string `json:"transType"`
}

//TransBatchResponse -- Trans Batch Response
type TransBatchResponse struct {
	OutBatchID 			string `json:"outBatchID"`
	Status 				string `json:"status"`
	TransList 			string `json:"transList"`
	OutOrderID 			string `json:"outOrderId"`
	FailReason  		string `json:"failReason"`
	FailList 			string `json:"failList"`
	Remark  			string `json:"remark"`
}

//ShareDisResponse -- Share Dis Response
type ShareDisResponse struct {
	OutBatchID 			string `json:"outBatchID"`
	Status 				string `json:"status"`
	Remark  			string `json:"remark"`
}

//BatchTransfer -- Batch Transfer
type BatchTransfer struct {
	OutBatchID 				string `json:"outBatchID"`
	TransferListStart  		string `json:"transferListStart"`
	OutOrderID 				string `json:"outOrderID"`
	PayUserID 				string `json:"payUserID"`
	DebitAccountNo 			string `json:"debitAccountNo"`
	ReceiveUserID  			string `json:"receiveUserID"`
	ReceiveUserAccountNo 	string `json:"receiveUserAccountNo"`
	Amount 					string `json:"amount"`
	FeeType 				string `json:"feeType"`
	FeeAmt 					string `json:"feeAmt"`
	TransferListEnd 		string `json:"transferListEnd"`
	Count 					string `json:"count"`
	Remark  				string `json:"Remark"`
	CallbackURL 			string `json:"callbackURL"`
}

//ShareDis -- Share Dis
type ShareDis struct {
	OutBatchID 				string `json:"outBatchID"`
	PayUserID 				string `json:"payUserID"`
	DebitAccountNo 			string `json:"debitAccountNo"`
	TransferListStart 		string `json:"transferListStart"`
	OutOrderID 				string `json:"outOrderID"`
	ReceiveUserID 			string `json:"receiveUserID"`
	Amount 					string `json:"amount"`
	FeeType 				string `json:"feeType"`
	FeeAmt 					string `json:"feeAmt"`
	TransferListEnd 		string `json:"transferListEnd"`
	Count 					string `json:"count"`
	Remark  				string `json:"remark"`
	CallbackURL 			string `json:"callbackURL"`
}

//QueryOrder -- Query Order
type QueryOrder struct {
	OutOrderID 			string `json:"outOrderID"`
	TransType 			string `json:"transType"`
}

//ConsumeBalance -- Consume Balance
type ConsumeBalance struct {
	OutOrderID 			string `json:"outOrderID"`
	OutUserID 			string `json:"outUserID"`
	Amount 				string `json:"amount"`
	ProductListStart 	string `json:"productListStart"`
	ProductName 		string `json:"productName"`
	ProductListEnd  	string `json:"productListEnd"`
	Remark 				string `json:"remark"`
}

//ConsumeConfirm -- Consume Confirm
type ConsumeConfirm struct {
	OutOrderID 			string `json:"outOrderID"`
	OrigOutOrderID 		string `json:"origOutOrderID"`
	Amount 				string `json:"amount"`
	ReceiveUserID 		string `json:"receiveUserID"`
	Remark 				string `json:"remark"`	
}

//QueryRevokeOrder -- Query RevokeOrder
type QueryRevokeOrder struct {
	OutOrderID 			string `json:"outOrderID"`
}

//EditAccountBalance -- Edit Account Balance
type  EditAccountBalance struct {
	OutOrderID 			string `json:"outOrderID"`
    OutUserID 			string `json:"outUserID"`
	VAcctNo 			string `json:"vAcctNo"`
	AccountName 		string `json:"accountName"`
	DefinedAccountNo 	string `json:"definedAccountNo"`
	Amount 				string `json:"amount"`
	Opt 				string `json:"opt"`
	Remark 				string `json:"remark"`
}

//TransUserDefinedBalance -- Trans User Defined Balance
type  TransUserDefinedBalance struct {
	OutOrderID 				string `json:"outOrderID"`
	TradeType 				string `json:"tradeType"`
	PayUserID 				string `json:"payUserID"`
	DebitAccountNo 			string `json:"debitAccountNo"`
	PayerAccountName 		string `json:"payerAccountName"`
	ReceiveUserID 			string `json:"receiveUserID"`
	ReceiveUserAccountNo 	string `json:"receiveUserAccountNo"`
	PayeeAccountName 		string `json:"payeeAccountName"`
	ProvisionsAccountNo 	string `json:"provisionsAccountNo"`
	Amount 					string `json:"amount"`
	Remark 					string `json:"remark"`
}

//QueryUserDefinedBalance -- Query User Defined Balance
type  QueryUserDefinedBalance struct {
	VAcctNo 				string `json:"vAcctNo"`
}

//QueryUserDefinedResponse --Query User Defined Response
type  QueryUserDefinedResponse struct {
	OrigOutOrderID 				string `json:"origOutOrderID"`
}

//RegisterMemberID --Register Member ID
type  RegisterMemberID struct {
	OutOrderID 				string `json:"outOrderID"`
	OutUserID 				string `json:"outUserID"`
	CustomerName 			string `json:"customerName"`
}

//UpgradeMemberID --Upgrade Member ID
type  UpgradeMemberID struct {
	OutOrderID 				string `json:"outOrderID"`
	OutUserID 				string `json:"outUserID"`
	VerifyCodeipID 			string `json:"verifyCodeipID"`
	CustomerName 			string `json:"customerName"`
	UserBizType 			string `json:"userBizType"`
	CustomerType  			string `json:"customerType"`
	CustomerCertType		string `json:"customerCertType"`
	CustomerCertNo 			string `json:"customerCertNo"`
	CustomerCertExpiry 		string `json:"customerCertExpiry"`
	CustomerCertExpiryend 	string `json:"customerCertExpiryend"`
	RegistCell 				string `json:"registCell"`
	CustomerEmail 			string `json:"customerEmail"`
	CustomerAddress 		string `json:"customerAddress"`
	Postcode 				string `json:"postcode"`
	CorpratePic1 			string `json:"corpratePic1"`
	CorpratePic2 			string `json:"corpratePic2"`
	ProtocolID 				string `json:"protocolID"`
	CorprateAddress 		string `json:"corprateAddress"`
	CustomerOccupation 		string `json:"customerOccupation"`
	Flage 					string `json:"flage"`
	VerifyCode 				string `json:"verifyCode"`
}

//ChannelPayment -- Channel Payment
type  ChannelPayment struct {
	OutOrderID 				string `json:"outOrderID"`
	OutUserID 				string `json:"outUserID"`
	Amount 					string `json:"amount"`
	OutPayID 				string `json:"outPayID"`
	PayType 				string `json:"payType"`
	RealName 				string `json:"realName"`
	PayChannelCode 			string `json:"payChannelCode"`
	CertNo 					string `json:"certNo"`
	BankAccountNo 			string `json:"bankAccountNo"`
	Phone 					string `json:"phone"`
	ValidDate 				string `json:"validDate"`
	Cvn2 					string `json:"cvn2"`
	ProductListStart 		string `json:"productListStart"`
 	ProductName 			string `json:"productName"`
	ProductListEnd 			string `json:"productListEnd"`
	Remark 					string `json:"remark"`
	TradeTime  				string `json:"tradeTime"`
	ReceiveUserID 			string `json:"receiveUserID"`
	ReconcileDate 			string `json:"reconcileDate"`
	SubPayChannelCode 		string `json:"subPayChannelCode"`
	PayTransactionID 		string `json:"payTransactionID"`
}

//OpenAccount -- Open Account
type OpenAccount struct {
	VAcctNo            string `json:"vAcctNo"`
	DefinedNo          string `json:"definedNo"`
}

//QueryMemberInfo -- query member info
type QueryMemberInfo struct {
	OutUserID            string `json:"outUserID"`
}

//QueryUpgrade
type QueryUpgrade struct {
	OrgOutOrderID  		string `json:"orgOutOrderID"`
	OutUserID  			string `json:"outUserID"`
	VipID 				string `json:"vipID"`
}

//FreezeFund -- Freeze customer Fund
type FreezeFund struct {
	OutOrderID  		string `json:"OutOrderID"`
	OutUserID  			string `json:"outUserID"`
	Amount 				string `json:"amount"`
	Remark				string `json:"remark"`
}
//UnFreezeFund -- unfreeze customer fund
type UnFreezeFund struct {
	OutOrderID  		string `json:"OutOrderID"`
	OutUserID  			string `json:"outUserID"`
	OrigOutOrderID  	string `json:"origOutOrderID"`
	Amount 				string `json:"amount"`
	Remark				string `json:"remark"`
}
//Recharge -- customer recharge fund
type Recharge struct {
	OutOrderID  		string `json:"OutOrderID"`
	OutUserID  			string `json:"outUserID"`
	OrigOutOrderID  	string `json:"origOutOrderID"`
	Amount 				string `json:"amount"`
	Remark				string `json:"remark"`
}

//TaxTrans -- tax trans
type TaxTrans struct {
	OutOrderID  		string `json:"OutOrderID"`
	OutUserID  			string `json:"outUserID"`
	Amount 				string `json:"amount"`
	BelongYears 		string `json:"belongYears"`
	BindCardID 			string `json:"bindCardID"`
	Note 				string `json:"note"`
}

//QueryTaxTrans -- query tax trans record
type QueryTaxTrans struct {
	OutOrderID  		string `json:"OutOrderID"`
}

//TaxFeeComputer -- tx fee calc
type TaxFeeComputer struct {
	OutUserID        string `json:"outUserID"`
	Amount           string `json:"amount"`
	BindCardID       string `json:"bindCardID"`
    Note  			 string `json:"note"`
}

//FundTrans --customer's fund transfer
type FundTrans struct {
	OutUserID 		string `json:"outUserID"`
	VAcctNo 		string `json:"vAcctNo"`
	StartTime 		string `json:"startTime"`
	EndTime 		string `json:"endTime"`
}

//VaccountFreeze --virtual account Freeze 
type VaccountFreeze struct {
	OutOrderID 		string `json:"outOrderID"`
	OutUserID       string `json:"outUserID"`
	VAcctNo 		string `json:"vAcctNo"`
	OperType		string `json:"operType"`
	Note			string `json:"note"`
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
