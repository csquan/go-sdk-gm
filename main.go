package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/internal2/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	ledger "github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	channelID      = "mychannel"
	orgName        = "Org1"
	orgAdmin       = "Admin"
	ordererOrgName = "OrdererOrg"
)

const (
	org1CfgPath = "./config.yaml"
	peer0Org1   = "peer0.org1.example.com"
	peer0Org2   = "peer0.org2.example.com"
)

var sdk *fabsdk.FabricSDK

func main() {

	//here have a trap of comma,use sdk1 to solve
	sdk1, err := fabsdk.New(config.FromFile(org1CfgPath))
	if err != nil {
		log.Panicf("failed to create fabric sdk: %s", err)
	}
	sdk = sdk1

	r := mux.NewRouter()
	r.HandleFunc("/users", login).Methods("POST")
	r.Handle("/channels", authMiddleware(http.HandlerFunc(createChannel))).Methods("POST")
	r.Handle("/channels/{channelName}/peers", authMiddleware(http.HandlerFunc(joinChannel))).Methods("POST")
	r.Handle("/channels/{channelName}/installchaincodes", authMiddleware(http.HandlerFunc(InstallChainCode))).Methods("POST")
	r.Handle("/channels/{channelName}/instantiatechaincodes", authMiddleware(http.HandlerFunc(InstantiateChainCode))).Methods("POST")
	r.Handle("/channels/{channelName}/chaincodes/{chaincodeName}", authMiddleware(http.HandlerFunc(queryCC))).Methods("GET")
	r.Handle("/channels/{channelName}/invokechaincodes/{chaincodeName}", authMiddleware(http.HandlerFunc(invokeCC))).Methods("POST")
	r.Handle("/channels/{channelName}/blocks/{blockID}", authMiddleware(http.HandlerFunc(getBlockByNumber))).Methods("GET")
	r.Handle("/channels/{channelName}/transactions/{transactionID}", authMiddleware(http.HandlerFunc(getTransactionByID))).Methods("GET")
	r.Handle("/channels/{channelName}", authMiddleware(http.HandlerFunc(getChainInfo))).Methods("GET")
	r.Handle("/chaincodes", authMiddleware(http.HandlerFunc(getInstalledChaincodes))).Methods("GET")
	r.Handle("/channels/{channelName}/chaincodes", authMiddleware(http.HandlerFunc(getInstantiatedChaincodes))).Methods("GET")
	r.Handle("/channels", authMiddleware(http.HandlerFunc(getChannels))).Methods("GET")
	http.ListenAndServe(":4000", handlers.LoggingHandler(os.Stdout, r))
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("authorization")
		if tokenString != "" {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("123"), nil
			})
			if err == nil {
				if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
					r.Header.Add("username", claims["username"].(string))
					r.Header.Add("orgName", claims["orgName"].(string))
					next.ServeHTTP(w, r)
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(err.Error()))
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

// GenerateRandomID generates random ID
func GenerateRandomID() string {
	return randomString(10)
}

// Utility to create random string of strlen length
func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	seed := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(seed)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rnd.Intn(len(chars))]
	}
	return string(result)
}

func login(w http.ResponseWriter, r *http.Request) {
	log.Print("================== LOGIN ==================")
	type response struct {
		Success bool
		Message string
		Token   string
	}

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	userName := r.Form.Get("username")
	orgName := r.Form.Get("orgName")
	secret := r.Form.Get("secret")
	if userName != "" && orgName != "" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": userName,
			"orgName":  orgName,
			"exp":      time.Now().Unix() + 360000,
		})
		tokenString, err := token.SignedString([]byte(secret))
		ctxProvider := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(orgName))
		if ctxProvider == nil {
			fmt.Println("failed to create ctxProvider")
			return
		}
		msp, err := mspclient.New(ctxProvider)
		if err != nil {
			fmt.Println(err)
			return
		}

		_, err = msp.GetSigningIdentity(userName)
		if err != nil {
			log.Printf("Check if user %s is enrolled: %s", userName, err.Error())
			testAttributes := []mspclient.Attribute{
				{
					Name:  GenerateRandomID(),
					Value: fmt.Sprintf("%s:ecert", GenerateRandomID()),
					ECert: true,
				},
				{
					Name:  GenerateRandomID(),
					Value: fmt.Sprintf("%s:ecert", GenerateRandomID()),
					ECert: true,
				},
			}
			// Register the new user
			identity, _ := msp.GetIdentity(userName)
			if true {
				log.Printf("User %s does not exist, registering new user", userName)
				_, err = msp.Register(&mspclient.RegistrationRequest{
					Name:        userName,
					Type:        orgName,
					Attributes:  testAttributes,
					Affiliation: orgName,
					Secret:      secret,
				})
			} else {
				log.Printf("Identity: %s", identity.Secret)
			}
			log.Printf("secret: %s ", secret)

		}

		err = msp.Enroll(userName, mspclient.WithSecret("123"))
		res := response{Success: true, Message: "success to enroll user"}
		if err != nil {
			log.Printf("enroll err：%s", err.Error())
			res.Success = false
			res.Message = err.Error()
		}
		res.Token = tokenString
		out, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func getChainInfo(w http.ResponseWriter, r *http.Request) {
	log.Print("================ GET CHANNEL INFORMATION ======================")
	// define response
	type response struct {
		Success bool
		Message string
	}
	vars := mux.Vars(r)
	username := "Admin"
	orgName := r.Header.Get("orgName")
	log.Print(username)
	log.Print(orgName)
	channelContext := sdk.ChannelContext(vars["channelName"], fabsdk.WithUser(username), fabsdk.WithOrg(orgName))
	client, err := ledger.New(channelContext)
	if err != nil {
		log.Fatalf("Failed to create new ledger client: %s", err)
	}

	blockchainInfo, _ := client.QueryInfo(ledger.WithTargetEndpoints(r.URL.Query().Get("peer")))
	type chainInfo struct {
		Height            uint64
		CurrentBlockHash  string
		PreviousBlockHash string
	}
	bci := chainInfo{}
	bci.Height = blockchainInfo.BCI.Height
	bci.CurrentBlockHash = fmt.Sprintf("%x", blockchainInfo.BCI.CurrentBlockHash)
	bci.PreviousBlockHash = fmt.Sprintf("%x", blockchainInfo.BCI.PreviousBlockHash)
	bcs, _ := json.Marshal(bci)

	res := response{
		Success: true,
		Message: string(bcs[:]),
	}
	ret, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func getInstalledChaincodes(w http.ResponseWriter, r *http.Request) {
	log.Print("================ GET INSTALLED CHAINCODES ======================")
	// define response
	type response struct {
		Success bool
		Message string
	}
	org1AdminClientContext := sdk.Context(fabsdk.WithUser("AdminOrg1"), fabsdk.WithOrg("Org1"))
	client, err := resmgmt.New(org1AdminClientContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client: %s", err)
	}
	endpoint := r.URL.Query().Get("peer")
	chaincodeQueryRes, err := client.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(endpoint))
	if err != nil {
		log.Fatalf("Failed to QueryInstalledChaincodes: %s", err.Error())
	}
	out1, _ := json.Marshal(chaincodeQueryRes)
	res := response{
		Success: true,
		Message: string(out1[:]),
	}

	out, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func getInstantiatedChaincodes(w http.ResponseWriter, r *http.Request) {
	log.Print("================ GET INSTANTIATED CHAINCODES ======================")
	// define response
	type response struct {
		Success bool
		Message string
	}
	vars := mux.Vars(r)
	org1AdminClientContext := sdk.Context(fabsdk.WithUser("AdminOrg1"), fabsdk.WithOrg("Org1"))
	client, err := resmgmt.New(org1AdminClientContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client: %s", err)
	}

	endpoint := r.URL.Query().Get("peer")
	chaincodeQueryRes, err := client.QueryInstantiatedChaincodes(vars["channelName"], resmgmt.WithTargetEndpoints(endpoint))
	if err != nil {
		log.Fatalf("Failed to QueryInstantiatedChaincodes: %s", err.Error())
	}
	out1, _ := json.Marshal(chaincodeQueryRes)

	res := response{
		Success: true,
		Message: string(out1[:]),
	}
	out, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
func createChannel(w http.ResponseWriter, r *http.Request) {
	log.Print("================ CREATE CHANNEL ======================")

	type response struct {
		Success bool
		Message string
	}

	type channelConfig struct {
		Name string
		Path string
		Org  string
	}

	decoder := json.NewDecoder(r.Body)
	channel := channelConfig{}
	decoder.Decode(&channel)

	clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(ordererOrgName))

	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Printf("Failed to create channel management client: %s", err.Error())
	}

	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
	if err != nil {
		log.Print(err)
	}
	adminIdentity, err := mspClient.GetSigningIdentity(orgAdmin)
	if err != nil {
		log.Print(err)
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfigPath: channel.Path,
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	_, err = resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com"))
	res := response{Success: true, Message: "success to create channel"}
	if err != nil {
		res.Message = err.Error()
		res.Success = false
	}
	out, err := json.Marshal(res)
	w.Write(out)
}
func joinChannel(w http.ResponseWriter, r *http.Request) {
	log.Print("================ JOIN CHANNEL ======================")

	type response struct {
		Success bool
		Message string
	}
	type joinConfig struct {
		Org string
	}
	decoder := json.NewDecoder(r.Body)
	join := joinConfig{}
	decoder.Decode(&join)

	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(join.Org))

	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		log.Panicf("Failed to create new resource management client: %s", err.Error())
	}

	// Org peers join channel
	if err = orgResMgmt.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com")); err != nil {
		log.Panicf("Org peers failed to JoinChannel: %s", err.Error())
	}
	res := response{
		Success: true,
		Message: "success to join channel",
	}
	if err != nil {
		res.Message = err.Error()
		res.Success = false
	}

	out, err := json.Marshal(res)
	w.Write(out)
}
func getBlockByNumber(w http.ResponseWriter, r *http.Request) {
	log.Print("==================== GET BLOCK BY NUMBER ==================")

	// define response
	type response struct {
		Success bool
		Message string
	}
	//username := "Admin" //r.Header.Get("username")
	//orgName := orgName
	//blockID, _ := strconv.ParseUint("1", 10, 64)

	clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(ordererOrgName))
	resMgmtClient, err := resmgmt.New(clientContext)

	ret, err := resMgmtClient.QueryConfigBlockFromOrderer(channelID, resmgmt.WithOrdererEndpoint("orderer.example.com"))
	res := response{
		Success: true,
		Message: "",
	}

	if err != nil {
		res.Message = err.Error()
		res.Success = false
		log.Panicf("QueryConfigBlockFromOrderer return error: %s", err.Error())
	} else {
		log.Print("====getBlockByNumber success=============")
		res.Message = ret.String()
	}
	out, err := json.Marshal(res)
	w.Write(out)
}
func queryCC(w http.ResponseWriter, r *http.Request) {
	log.Print("================query cc======================")
	type response1 struct {
		Success bool
		Message string
	}

	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser("Admin"))
	cc, err := channel.New(ccp)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err.Error())
	}
	// new channel request for query
	req := channel.Request{
		ChaincodeID: "face",
		Fcn:         "queryAllFace",
		Args:        packArgs([]string{"1111110", "1111112"}),
	}
	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer0Org1)
	response, err := cc.Query(req, reqPeers)

	res := response1{
		Success: true,
		Message: "",
	}

	if err != nil {
		res.Message = err.Error()
		res.Success = false
	} else {
		log.Printf("chaincode query success,the value is %s\n", string(response.Payload))
		res.Message = string(response.Payload)
	}

	out, err := json.Marshal(res)
	w.Write(out)
}

func invokeCC(w http.ResponseWriter, r *http.Request) {
	log.Print("=====================invoke chaincode====================")
	type response1 struct {
		Success bool
		Message string
	}
	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser("Admin"))
	cc, err := channel.New(ccp)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err.Error())
	}

	args := packArgs([]string{"1111111", "2222222", "33333333", "xxxx"})
	req := channel.Request{
		ChaincodeID: "face",
		Fcn:         "createFace",
		Args:        args,
	}
	peers := []string{peer0Org1}
	reqPeers := channel.WithTargetEndpoints(peers...)
	response, err := cc.Execute(req, reqPeers)

	res := response1{
		Success: true,
		Message: "",
	}
	if err != nil {
		res.Success = false
		res.Message = err.Error()
		log.Printf("failed to Execute chaincode: %s\n", err.Error())
	} else {
		res.Message = string(response.TransactionID)
	}
	log.Printf("Execute chaincode success,txId:%s\n", response.TransactionID)

	out, err := json.Marshal(res)
	w.Write(out)
}

func packArgs(paras []string) [][]byte {
	var args [][]byte
	for _, k := range paras {
		args = append(args, []byte(k))
	}
	return args
}

//InstallChainCode function
func InstallChainCode(w http.ResponseWriter, r *http.Request) {
	log.Print("================install cc======================")
	//define response
	type response struct {
		Success bool
		Message string
	}
	ccPkg, err := packager.NewCCPackage("sdk_test/src", "/root/gowork")
	if err != nil {
		log.Fatalf("pack chaincode error %s", err.Error())
	}
	// new request of installing chaincode.
	req := resmgmt.InstallCCRequest{
		Name:    "face",
		Path:    "sdk_test/src",
		Version: "v0",
		Package: ccPkg,
	}

	clientContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("org2"))

	// Resource management client is responsible for managing resources (joining channels, install/instantiate/upgrade chaincodes).
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client")
	}

	reqPeers := resmgmt.WithTargetEndpoints("peer0.org2.example.com")
	_, err = resMgmtClient.InstallCC(req, reqPeers)

	res := response{
		Success: true,
		Message: "success to install chaincode",
	}
	if err != nil {
		res.Message = err.Error()
		res.Success = false
	}

	out, err := json.Marshal(res)
	w.Write(out)
}

//InstantiateChainCode function
func InstantiateChainCode(w http.ResponseWriter, r *http.Request) {
	log.Print("================instantiate cc======================")
	//define response
	type response struct {
		Success bool
		Message string
	}
	ccPolicy := cauthdsl.SignedByMspMember("Org1MSP")

	req := resmgmt.InstantiateCCRequest{
		Name:    "face",
		Path:    "sdk_test/src",
		Version: "v0",
		Policy:  ccPolicy,
	}

	clientContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("org1"))

	// Resource management client is responsible for managing resources (joining channels, install/instantiate/upgrade chaincodes)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client")
	}

	reqPeers := resmgmt.WithTargetEndpoints("peer0.org1.example.com")
	_, err = resMgmtClient.InstantiateCC("mychannel", req, reqPeers)

	res := response{
		Success: true,
		Message: "success to instantiate chaincode",
	}
	if err != nil {
		res.Message = err.Error()
		res.Success = false
	}

	out, err := json.Marshal(res)
	w.Write(out)
}

func getChannels(w http.ResponseWriter, r *http.Request) {
	log.Print("================ GET CHANNELS ======================")

	type response struct {
		Success bool
		Message string
	}
	clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	resMgmtClient, err := resmgmt.New(clientContext)

	if err != nil {
		log.Fatalf("Failed to create new resource management client: %s", err.Error())
	}
	ret, err := resMgmtClient.QueryChannels(resmgmt.WithTargetEndpoints("peer0.org1.example.com"))

	res := response{
		Success: true,
		Message: "",
	}

	if err == nil {
		log.Print("================ GET CHANNELS SUCCESS ======================")
		out, err1 := json.Marshal(ret)
		if err1 == nil {
			res.Message = string(out[:])
		}
	} else {
		res.Success = false
		res.Message = err.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	out, err := json.Marshal(res)
	w.Write(out)
}

func getTransactionByID(w http.ResponseWriter, r *http.Request) {
	log.Print("================ GET TRANSACTION BY TRANSACTION_ID ======================")

	vars := mux.Vars(r)

	log.Print(orgName)
	type response struct {
		Success bool
		Message string
	}

	channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))

	client, err := ledger.New(channelContext)
	if err != nil {
		log.Fatalf("Failed to create new channel client: %s", err.Error())
	}
	res := response{
		Success: true,
		Message: "",
	}
	txid := vars["transactionID"]
	log.Printf("txid:%s", txid)
	ret, err := client.QueryTransaction(fab.TransactionID(txid))
	if err != nil {
		res.Success = false
		res.Message = err.Error()
		log.Printf("Failed to queryTx : %s", err.Error())
	} else {
		res.Message = ret.String()
	}
	out, err := json.Marshal(res)
	w.Write(out)
}
