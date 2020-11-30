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

//UploadPicture -- upload picture
type UploadPicture struct {
	OutOrderI			string `json:"outOrderID"`
	OurUserID			string `json:"ourUserID"`
	Pic					string `json:"pic"`
	TemplateNo          string `json:"templateNo"`
	FileType            string `json:"fileType"`
}

//RegisterPubilcUser -- Pubilc User Register
type RegisterPubilcUser struct {
	OutOrderID							string `json:"outOrderID"`
	OurUserID							string `json:"ourUserID"`
	EnterpriseName						string `json:"enterpriseName"`
	CustomerRole						string `json:"customerRole"`
	UserBizType		    				string `json:"userBizType"`
	CustomerPhone       				string `json:"customerPhone"`
	BusinessLicenseType 				string `json:"businessLicenseType"`
	BusinessLicense 					string `json:"businessLicense"`
	CustomerCertExpiry 					string `json:"customerCertExpiry"`
	CustomerCertExpiryEnd 				string `json:"customerCertExpiryEnd"`
	CustomerEmail 						string `json:"customerEmail"`
	CustorAddress 						string `json:"customerAddress "`
	BusinessScope 						string `json:"businessScope"`
	PostCode 							string `json:"postCode"`
	CustomerCertPid 					string `json:"customerCertPid"`  
	LegalName 							string `json:"legalName"`
	LegealCertType 						string `json:"legealCertType"`
	LegalCertNo 						string `json:"legalCertNo"`
	CorprateCertExpiry 					string `json:"corprateCertExpiry"`
	CorprateCertExpiryEnd 				string `json:"corprateCertExpiryEnd"`
	CorprateMobile 						string `json:"corprateMobile"`
	CorpratePic1 						string `json:"crpratePic1"`
	CorpratePic2 						string `json:"corpratePic2"`
	OperatorName 						string `json:"operatorName"`
	OperatorCertType 					string `json:"operatorCertType"`
	OperateCertNo 						string `json:"operateCertNo"`
	OperatorCertExpiry 					string `json:"operatorCertExpiry"`
	OperatorCertExpiryEnd 				string `json:"operatorCertExpiryEnd"`
	OperatorMobile 						string `json:"operatorMobile"`
	AuthLetterPicID 					string `json:"authLetterPicID"`
	ProtoColld 							string `json:"protoColld"`
	ControlShareHolderA 				string `json:"controlShareHolderA"`
	ControlShareHolderACertType 		string `json:"controlShareHolderACertType"`
	ControlShareACertNo 				string `json:"controlShareACertNo"`
	ControlShareACertExpiry 			string `json:"controlShareACertExpiry"`
	ControlShareholderACertExpiryEnd 	string `json:"controlShareholderACertExpiryEnd"`
	ControlShareholderB 				string `json:"controlShareholderB"`
	ControlShareholderBCertType 		string `json:"controlShareholderBCertType"`
	ControlShareholderBCertNo 			string `json:"controlShareholderBCertNo"`
	ControlShareholderBCertExpiry 		string `json:"controlShareholderBCertExpiry"`
	ControlShareholderBCertExpiryEnd 	string `json:"controlShareholderBCertExpiryEnd"`
	ControlShareholderC 				string `json:"cControlShareholderC"`
	ControlShareholderCCertType 		string `json:"controlShareholderCCertType"`
	ControlShareholderCCertNo 			string `json:"controlShareholderCCertNo"`
	ControlShareholderCCertExpiry 		string `json:"controlShareholderCCertExpiry"`
	ControlShareholderCCertExpiryEnd 	string `json:"controlShareholderCCertExpiryEnd"`
	Remark 								string `json:"remark"`
	CallbackURL  						string `json:"callbackURL"`
}

//EditPubilcUser -- edit Pubilc User
type EditPubilcUser struct {
	OutOrderID         					string `json:"outOrderID"`
	OurUserID          					string `json:"ourUserID"`
	CustomerRole        				string `json:"customerRole"`
	CustomerPhone       				string `json:"customerPhone"`
	BusinessLicenseType 				string `json:"businessLicenseType"`
	BusinessLicense 					string `json:"businessLicense"`
	CustomerCertExpiry 					string `json:"customerCertExpiry"`
	CustomerCertExpiryEnd 				string `json:"customerCertExpiryEnd"`
	CustomerEmail 						string `json:"customerEmail"`
	CustorAddress 						string `json:"customerAddress "`
	PostCode 							string `json:"postCode"`
	CustomerCertPid 					string `json:"customerCertPid"`  
	LegalName 							string `json:"legalName"`
	LegealCertType 						string `json:"legealCertType"`
	LegalCertNo 						string `json:"legalCertNo"`
	CorprateCertExpiry 					string `json:"corprateCertExpiry"`
	CorprateCertExpiryEnd 				string `json:"corprateCertExpiryEnd"`
	CorprateMobile 						string `json:"corprateMobile"`
	CorpratePic1 						string `json:"crpratePic1"`
	CorpratePic2 						string `json:"corpratePic2"`
	OperatorName 						string `json:"operatorName"`
	OperatorCertType 					string `json:"operatorCertType"`
	OperateCertNo 						string `json:"operateCertNo"`
	OperatorCertExpiry 					string `json:"operatorCertExpiry"`
	OperatorCertExpiryEnd 				string `json:"operatorCertExpiryEnd"`
	OperatorMobile 						string `json:"operatorMobile"`
	AuthLetterPicID 					string `json:"authLetterPicID"`
	ProtoColld 							string `json:"protoColld"`
	Remark                              string `json:"remark"`
    EnterpriseName						string `json:"enterpriseName"`
}

//PersonalRegister -- Personal Register
type PersonalRegister struct {
	OutOrderID 				string `json:"outOrderID"`
	OutUserID 				string `json:"outUserID"`
	CustomerName 			string `json:"customerName"`
	UserBizType 			string `json:"userBizType"`
	CustomerType 			string `json:"customerType"`
	CustomerCertType 		string `json:"customerCertType"`
	CustomerCertNo 			string `json:"customerCertNo"`
	CustomerCertExpiry 		string `json:"customerCertExpiry"`
	CustomerCertExpiryend 	string `json:"customerCertExpiryend"`
	RegionCode 				string `json:"regionCode"`
	RegistCell 				string `json:"registCell"`
	CustomerEmail 			string `json:"customerEmail"`
	CustomerAddress 		string `json:"customerAddress"`
	Postcode 				string `json:"postcode"`
	CorpratePic1 			string `json:"corpratePic1"`
	CorpratePic2 			string `json:"corpratePic2"`
	ProtocolID 				string `json:"protocolID"`
	CorprateAddress 		string `json:"corprateAddress"`
	Sex 					string `json:"sex"`
	Citizenship 			string `json:"citizenship"`
	AlwaysLiveHome 			string `json:"alwaysLiveHome"`
	CustomerOccupation 		string `json:"customerOccupation"`
	Flage 					string `json:"flage"`
	VerifyCode 				string `json:"verifyCode"`
	Remark  				string `json:"remark"`
}

//EditPersonal -- edit Personal
type EditPersonal struct {
	OutOrderID 				string `json:"outOrderID"`
	OutUserID 				string `json:"outUserID"`
	UserBizType 			string `json:"userBizType"`
	CustomerCertExpiry 		string `json:"customerCertExpiry"`
	CustomerCertExpiryend 	string `json:"customerCertExpiryend"`
	CustomerPhone           string `json:"customerPhone"`
	CustomerEmail           string `json:"customerEmail"`
	CustomerAddress         string `json:"customerAddress"`
	Postcode          		string `json:"postcode"`
	CorpratePic1 			string `json:"corpratePic1"`
	CorpratePic2 			string `json:"corpratePic2"`
	CorprateAddress 		string `json:"corprateAddress"`
	Sex 					string `json:"sex"`
	Citizenship 			string `json:"citizenship"`
	AlwaysLiveHome 			string `json:"alwaysLiveHome"`
	CustomerOccupation 		string `json:"customerOccupation"`
	RegionCode 				string `json:"regionCode"`
	RegistCell 				string `json:"registCell"`
	Remark  				string `json:"remark"`
}

//CustomerAccountTie -- Customer Account Tie
type CustomerAccountTie struct {
	OutOrderID           	    string `json:"outOrderID"`
	VAcctNo          		string `json:"vAcctNo"`
	OutUserID      			string `json:"outUserID"`

	BindType 				string `json:"bindType"`
	AccountAttr				string `json:"accountAttr"`
	BankAccountNo			string `json:"bankAccountNo"`
	BankAccountName			string `json:"bankAccountName"`

	OpeningPermitPic		string `json:"openingPermitPic"`
	AuthPic					string `json:"authPic"`
	Flage					string `json:"flage"`
	VerifyCode				string `json:"verifyCode"`
	RealName				string `json:"realName"`
	CertNo					string `json:"certNo"`
	Remark					string `json:"remark"`
}

//PublicAccountVerify -- Public Account Verify
type PublicAccountVerify struct {
	OrigOutOrderID          string `json:"origOutOrderID"`
	OutOrderID    			string `json:"outOrderID"`
	OutUserID     			string `json:"outUserID"`

	BindType 				string `json:"bindType"`
	BindCardID				string `json:"bindCardID"`
	Amount					string `json:"amount"`
}

//AccountUntie -- Account Untie
type AccountUntie struct {
	OutOrderID    			string `json:"outOrderID"`
	OutUserID     			string `json:"outUserID"`
	VAcctNo  				string `json:"vAcctNo"`
	OldBindType 			string `json:"oldBindType"`
	BindCardID 				string `json:"bindCardID"`
	Flage					string `json:"flage"`
	VerifyCode				string `json:"verifyCode"`
	Remark    				string `json:"remark"`
}

//AccountWithdraw -- Account Withdraw
type AccountWithdraw struct {
	OutOrderID    			string `json:"outOrderID"`
	OutUserID     			string `json:"outUserID"`
	DebitAccountNo 			string `json:"debitAccountNo"`
	Amount 					string `json:"amount"`
	FeeType 				string `json:"feeType"`
	FeeAmount 				string `json:"feeAmount"`
	BindCardID 				string `json:"bindCardID"`

	Flage					string `json:"flage"`
	VerifyCode				string `json:"verifyCode"`
	CallbackURL 			string `json:"callbackURL"`
	Remark    				string `json:"remark"`
}

//QueryWithdraw -- Query Withdraw
type QueryWithdraw struct {
	OutOrderID    			string `json:"outOrderID"`
	OutUserID 				string `json:"outUserID"`
	OrigOutOrderID 			string `json:"origOutOrderID"`
	OldTradeDate 			string `json:"oldTradeDate"`
	TransType 				string `json:"transType"`
}

//AccountTransfer -- Account Transfer
type AccountTransfer struct {
	OutOrderID    			string `json:"outOrderID"`
	PayUserID 				string `json:"payUserID"`
	DebitAccountNo 			string `json:"debitAccountNo"`
	ReceiveUserID 			string `json:"receiveUserID"`
	ReceiveUserAccountNo 	string `json:"receiveUserAccountNo"`
	Amount 					string `json:"amount"`
	FeeType 				string `json:"feeType"`
	FeeAmount 				string `json:"feeAmount"`
	Currency 				string `json:"currency"`
	Summary 				string `json:"summary"`
	Remark 					string `json:"remark"`
}

//QueryTransfer -- Quers Transfer
type QueryTransfer struct {
	OutOrderID    			string `json:"outOrderID"`
	OrigOutOrderID 			string `json:"origOutOrderID"`
	OldTradeDate 			string `json:"oldTradeDate"`
}

//QueryPublicAccount -- Query Public Account
type QueryPublicAccount struct {
	OutUserID    			string `json:"outUserID"`
}

//QueryPersonalAccount - -Query Personal Account
type QueryTransHistory struct {
	outPlatformId
	prdtNo
	priThirdTradeNo
	outUserId
	outAccSeqn
	tradeType
	addFlag
	begDate
	endDate
	currPage
	eachPageNum
	sortType




	OutUserID    			string `json:"outUserID"`
}

//QueryVAccount -- Query Virtual Account
type QueryVAccount struct {
	OutUserID    			string `json:"outUserID"`
	VAcctNo  				string `json:"VAcctNo"`
}

//QueryAccountTie -- Query AccountTie
type QueryAccountTie struct {
	OutUserID    			string `json:"outUserID"`
	VAcctNo  				string `json:"VAcctNo"`
}

//QueryTransRecord -- Query Trans Record
type QueryTransRecord struct {
	OutUserID    			string `json:"outUserID"`
	VAcctNo  				string `json:"VAcctNo"`
    StartTime               string `json:"startTime"`
	EndTime   				string `json:"endTime"`
}

//RegisterOrModifyResponse -- Register Or Modify Response
type RegisterOrModifyResponse struct {
	OutUserID    			string `json:"outUserID"`
	OrigOutOrderID 			string `json:"origOutOrderID"`
	Status 					string `json:"status"`
	CheckMsg 				string `json:"checkMsg"`
	VAcctNo  				string `json:"vAcctNo"`
}

//WithdrawResponse -- Withdraw Response
type WithdrawResponse struct {
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
