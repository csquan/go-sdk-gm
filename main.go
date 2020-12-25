package main

import (
	//"context"
	//"crypto/tls"
	//"crypto/x509"
	//"github.com/hyperledger/fabric/common/util"
	//"github.com/hyperledger/fabric/protos/discovery"
	"encoding/json"
	"fmt"

	//"bytes"
	//"encoding/gob"
	//"encoding/binary"
	//"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	//	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	//"io/ioutil"
	//"os/exec"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hyperledger/fabric-sdk-go/internal2/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	ledger "github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	_ "github.com/hyperledger/fabric-sdk-go/pkg/fab/events/client"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	_ "github.com/stretchr/testify/assert"

	//	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"
	//client "github.com/hyperledger/fabric/discovery/client"
	//disc "github.com/hyperledger/fabric-sdk-go/internal2/github.com/hyperledger/fabric/discovery/client"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk/factory/defsvc"
	//"github.com/hyperledger/fabric-sdk-go/test/integration"
	"github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/comm"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/discovery"
)

const (
	orgName        = "org1"
	orgAdmin       = "Admin"
	ordererOrgName = "OrdererOrg"
)

const (
	org1CfgPath = "./config.yaml"
	peer0Org1   = "peer0.org1.example.com"
	peer0Org2   = "peer0.org2.example.com"
	peer0Org3   = "peer0.org3.example.com"
	peer0Org4   = "peer0.org4.example.com"
)

var sdk *fabsdk.FabricSDK

func main() {

	//here have a trap of comma,use sdk1 to solve
	sdk1, err := fabsdk.New(config.FromFile(org1CfgPath))
	if err != nil {
		log.Panicf("failed to create fabric sdk: %s", err)
	}
	sdk = sdk1

	storePath := "/tmp/examplestore"
	err = os.RemoveAll(storePath)
	if err != nil {
		log.Fatalf("Cleaning up directory '%s' failed: %v", storePath, err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/users", login).Methods("GET")
	r.Handle("/channels", authMiddleware(http.HandlerFunc(createChannel))).Methods("POST")
	r.Handle("/channels/{channelName}/peers", authMiddleware(http.HandlerFunc(joinChannel))).Methods("POST")
	r.Handle("/channels/{channelName}/installchaincodes", authMiddleware(http.HandlerFunc(InstallChainCode))).Methods("POST")
	r.Handle("/channels/{channelName}/instantiatechaincodes", authMiddleware(http.HandlerFunc(InstantiateChainCode))).Methods("POST")
	r.Handle("/channels/{channelName}/upgradechaincodes", authMiddleware(http.HandlerFunc(UpgradeChainCode))).Methods("POST")
	r.Handle("/channels/{channelName}/chaincodes/{chaincodeName}", authMiddleware(http.HandlerFunc(queryCC))).Methods("GET")
	r.Handle("/channels/{channelName}/invokechaincodes/{chaincodeName}", authMiddleware(http.HandlerFunc(invokeCC))).Methods("POST")
	r.Handle("/channels/{channelName}/blocks", authMiddleware(http.HandlerFunc(getBlockByNumber))).Methods("GET")
	r.Handle("/channels/configblock", authMiddleware(http.HandlerFunc(QueryConfigBlockFromOrderer))).Methods("GET")
	r.Handle("/channels/tx", authMiddleware(http.HandlerFunc(getTransactionByID))).Methods("GET")
	r.Handle("/channels/{channelName}", authMiddleware(http.HandlerFunc(getChainInfo))).Methods("GET")
	r.Handle("/blockbyhash", authMiddleware(http.HandlerFunc(getBlockByHash))).Methods("GET")
	r.Handle("/blockbytxid", authMiddleware(http.HandlerFunc(getBlockByTXID))).Methods("GET")
	r.Handle("/channelconfig", authMiddleware(http.HandlerFunc(QueryChannelConfig))).Methods("GET")
	r.Handle("/getpeers", authMiddleware(http.HandlerFunc(GetPeers))).Methods("GET")
	r.Handle("/ispeerinchannel", authMiddleware(http.HandlerFunc(IsPeerInChannel))).Methods("GET")
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
	user := r.Form.Get("username")
	org := r.Form.Get("orgName")
	secret := r.Form.Get("secret")
	log.Print(r)
	if user != "" && org != "" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user,
			"orgName":  org,
			"exp":      time.Now().Unix() + 360000,
		})
		tokenString, err := token.SignedString([]byte(secret))
		ctxProvider := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(org))
		if ctxProvider == nil {
			log.Fatalf("failed to create ctxProvider")
		}

		msp, err := mspclient.New(ctxProvider)
		if err != nil {
			log.Fatal("failed to call new for create msp", err.Error())
		}

		_, err = msp.GetSigningIdentity(user)
		if err != nil {
			log.Printf("Check if user %s is enrolled: %s", user, err.Error())
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
			identity, _ := msp.GetIdentity(user)
			if true {
				log.Printf("User %s does not exist, registering new user", user)

				_, err = msp.Register(&mspclient.RegistrationRequest{
					Name:        user,
					Type:        org,
					Attributes:  testAttributes,
					Affiliation: org,
					Secret:      secret,
				})
			} else {
				log.Printf("Identity: %s", identity.Secret)
			}
			log.Printf("secret: %s ", secret)

		}
		err = msp.Enroll(user, mspclient.WithSecret("123"))
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
	user := "Admin"
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	channelID := r.Form.Get("channelID")

	org := r.Header.Get("orgName")
	channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(user), fabsdk.WithOrg(org))
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
	org1AdminClientContext := sdk.Context(fabsdk.WithUser("AdminOrg1"), fabsdk.WithOrg("Org1"))
	client, err := resmgmt.New(org1AdminClientContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client: %s", err)
	}

	endpoint := r.URL.Query().Get("peer")

	channelID := r.URL.Query().Get("channelID")
	chaincodeQueryRes, err := client.QueryInstantiatedChaincodes(channelID, resmgmt.WithTargetEndpoints(endpoint))
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
	req := resmgmt.SaveChannelRequest{ChannelID: channel.Name,
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
		ChannelID string
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
	err = orgResMgmt.JoinChannel(join.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com"));

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

func QueryConfigBlockFromOrderer(w http.ResponseWriter, r *http.Request) {
	log.Print("====================QueryConfigBlockFromOrderer ==================")

	// define response
	type response struct {
		Success bool
		Message string
	}
	//username := "Admin" //r.Header.Get("username")
	//orgName := orgName
	//blockID, _ := strconv.ParseUint("1", 10, 64)
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	channelID := r.Form.Get("channelID")

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
		log.Print("====QueryConfigBlockFromOrderer success=============")
		log.Print(ret)
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

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
        channelID := r.Form.Get("channelID")
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
		TxID    string
	}

	type invokeConfig struct {
		ChannelID string
		Peers  []string
		Fcn    string
		Args   []string
	}

	decoder := json.NewDecoder(r.Body)
	invoke := invokeConfig{}
	decoder.Decode(&invoke)

	vars := mux.Vars(r)

	ccp := sdk.ChannelContext(invoke.ChannelID, fabsdk.WithUser("Admin"))
	cc, err := channel.New(ccp)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err.Error())
	}

	args := packArgs(invoke.Args)


	req := channel.Request{
		ChaincodeID: vars["chaincodeName"],
		Fcn:         invoke.Fcn,
		Args:        args,
	}
	peers := []string{peer0Org1}
	reqPeers := channel.WithTargetEndpoints(peers...)
	response, err := cc.Execute(req, reqPeers)

	res := response1{
		Success: true,
		Message: "",
		TxID:"",
	}
	if err != nil {
		res.Success = false
		res.Message = err.Error()
		log.Printf("failed to Execute chaincode: %s\n", err.Error())
	} else {
		log.Print("+++++++++++++++response++++++++++++++")
		log.Print(string(response.Payload))
		res.Message = string(response.Payload)
		res.TxID = string(response.TransactionID)
	}
	log.Printf("Execute chaincode success,txId:%s\n", response.TransactionID)
	fmt.Println(res)
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
	type installConfig struct {
		Name    string
		Path    string
		Version string
		Org     string
		User    string
	}

	decoder := json.NewDecoder(r.Body)
	install := installConfig{}
	decoder.Decode(&install)
	log.Print(install)

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	peer := r.Form.Get("peer")
	log.Print(peer)
	//define response
	type response struct {
		Success bool
		Message string
	}
	ccPkg, err := packager.NewCCPackage(install.Path, "/root/go")
	if err != nil {
		log.Fatalf("pack chaincode error %s", err.Error())
	}
	// new request of installing chaincode.
	req := resmgmt.InstallCCRequest{
		Name:    install.Name,
		Path:    install.Path,
		Version: install.Version,
		Package: ccPkg,
	}

	clientContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(install.Org))

	// Resource management client is responsible for managing resources (joining channels, install/instantiate/upgrade chaincodes).
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client")
	}

	reqPeers := resmgmt.WithTargetEndpoints(peer)
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

	type instantiateConfig struct {
		ChannelID string
		Name    string
		Version string
		Path    string
		Args    []string
	}
	decoder := json.NewDecoder(r.Body)
	instantiate := instantiateConfig{}
	decoder.Decode(&instantiate)


	ccPolicy := cauthdsl.SignedByMspMember("Org1MSP")

	args := packArgs(instantiate.Args)

	req := resmgmt.InstantiateCCRequest{
		Name:    instantiate.Name,
		Path:    instantiate.Path,
		Version: instantiate.Version,
		Policy:  ccPolicy,
		Args:    args,
	}

	clientContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("org1"))

	// Resource management client is responsible for managing resources (joining channels, install/instantiate/upgrade chaincodes)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client")
	}

	reqPeers := resmgmt.WithTargetEndpoints("peer0.org1.example.com")
	_, err = resMgmtClient.InstantiateCC(instantiate.ChannelID, req, reqPeers)

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

//UpgradeChainCode function
func UpgradeChainCode(w http.ResponseWriter, r *http.Request) {
	log.Print("================Upgrade cc======================")
	//define response
	type response struct {
		Success bool
		Message string
	}

	type upgradeConfig struct {
		ChannelID string
		Name    string
		Version string
		Path    string
		Args    string
	}

	decoder := json.NewDecoder(r.Body)
	upgrade := upgradeConfig{}
	decoder.Decode(&upgrade)

	log.Print("++++upgrade+++")
	log.Print(upgrade)

	ccPolicy := cauthdsl.SignedByMspMember("Org1MSP")


	req := resmgmt.UpgradeCCRequest{
		Name:    upgrade.Name,
		Path:    upgrade.Path,
		Version: upgrade.Version,
		Policy:  ccPolicy,
	}

	clientContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("org1"))

	// Resource management client is responsible for managing resources (joining channels, install/instantiate/upgrade chaincodes)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client")
	}

	reqPeers := resmgmt.WithTargetEndpoints("peer0.org1.example.com")
	_, err = resMgmtClient.UpgradeCC(upgrade.ChannelID, req, reqPeers)

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

// QueryChannels queries the names of all the channels that a peer has joined.
func getChannels(w http.ResponseWriter, r *http.Request) {
	log.Print("================ GET CHANNELS ======================")

	type response struct {
		Success bool
		Message string
	}

	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	peer := r.Form.Get("peer")
	org := r.Form.Get("org")
	clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(org))
	resMgmtClient, err := resmgmt.New(clientContext)

	if err != nil {
		log.Fatalf("Failed to create new resource management client: %s", err.Error())
	}
	ret, err := resMgmtClient.QueryChannels(resmgmt.WithTargetEndpoints(peer))

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

func getBlockByHash(w http.ResponseWriter, r *http.Request) {
	log.Print("================ GET BLOCK BY HASH ======================")
	// define response
	type response struct {
		Success bool
		Message string
	}
	user := "Admin"
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	channelID := r.Form.Get("channelID")
	peer := r.Form.Get("peer")
	org := r.Form.Get("orgName")
	channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(user), fabsdk.WithOrg(org))
	client, err := ledger.New(channelContext)
	if err != nil {
		log.Fatalf("Failed to create new ledger client: %s", err)
	}
	blockHash := r.Form.Get("blockhash")
	blockHash = "yzHpSOo7Yb7NHP8BTQs6dLueutDqC4VtEfcO90+Nn80"
	block, _ := client.QueryBlockByHash([]byte(blockHash),ledger.WithTargetEndpoints(peer))

	ret, err := json.Marshal(block)
	log.Print("================ GET BLOCK BY HASH Start======================")
	fmt.Print(block)
	fmt.Print(blockHash)
	log.Print("================ GET BLOCK BY HASH End======================")
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}



func getBlockByTXID(w http.ResponseWriter, r *http.Request) {
	log.Print("================ QueryBlockByTxID ======================")
	// define response
	type response struct {
		Success bool
		Message string
	}
	user := "Admin"
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	channelID := r.Form.Get("channelID")

	org := r.Header.Get("orgName")
	channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(user), fabsdk.WithOrg(org))
	client, err := ledger.New(channelContext)
	if err != nil {
		log.Fatalf("Failed to create new ledger client: %s", err)
	}
	txid := r.Form.Get("txid") 
	block, _ := client.QueryBlockByTxID(fab.TransactionID(txid),ledger.WithTargetEndpoints(r.URL.Query().Get("peer")))

	//write block to file,then use confixlator to parse and get info  here

	ret, err := json.Marshal(block)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func getTransactionByID(w http.ResponseWriter, r *http.Request) {
	log.Print("================ getTransactionByID ======================")
	// define response
	type response struct {
		Success bool
		Message string
	}
	user := "Admin"
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	channelID := r.Form.Get("channelID")

	org := r.Header.Get("orgName")
	channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(user), fabsdk.WithOrg(org))
	client, err := ledger.New(channelContext)
	if err != nil {
		log.Fatalf("Failed to create new ledger client: %s", err)
	}
	txid := r.Form.Get("txid") 
	tx, err := client.QueryTransaction(fab.TransactionID(txid),ledger.WithTargetEndpoints(r.URL.Query().Get("peer")))

	res := response{
		Success: true,
		Message: "",
	}

	if err != nil {
		res.Success = false
		res.Message = err.Error()
		log.Printf("Failed to queryTx : %s", err.Error())
	} else {
		res.Message = tx.String()
	}
	out, err := json.Marshal(res)
	w.Write(out)
}


func getBlockByNumber(w http.ResponseWriter, r *http.Request) {
	log.Print("================ getBlockByNumber ======================")
	// define response
	type response struct {
		Success bool
		Message string
	}
	user := "Admin"
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	channelID := r.Form.Get("channelID")
	blockID := r.Form.Get("blockID")
        log.Print(blockID)
	blockNumber,err := strconv.ParseUint(blockID,10,64)
	if err != nil {
		panic(err)
	}

	org := r.Header.Get("orgName")
	channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(user), fabsdk.WithOrg(org))
	client, err := ledger.New(channelContext)
	if err != nil {
		log.Fatalf("Failed to create new ledger client: %s", err)
	}

	block, _ := client.QueryBlock(blockNumber,ledger.WithTargetEndpoints(r.URL.Query().Get("peer")))

	//write block to file,use configtxlator to parse ,get useful info 

	//ret1, err := json.Marshal(block)

	/*ioutil.WriteFile("./block.txt",ret1,0777)

	cmd := exec.Command("/bin/bash", "-c", "./artifacts/channel2/configtxlator  proto_decode  --type common.Block --input ./block.txt")

	buf, err := cmd.Output()
	if err != nil{
		fmt.Println(err.Error())
	}
	fmt.Println("++++++++++++++++++++++++++++++++++++++parse block ret+++++++++++++++++++++++++++++++++++")
	fmt.Print(buf)*/
	pos := strings.IndexAny(block.Data.String(),"@")
	txid, err := json.Marshal(block.Data.String()[pos+1:pos+65])


	res := response{
		Success: true,
		Message: "txid:" + string(txid),
	}
	ret, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func QueryChannelConfig(w http.ResponseWriter, r *http.Request) {
	log.Print("================ QueryChannelConfig ======================")
	// define response
	type response struct {
		Success bool
		Message string
	}
	user := "Admin"
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
 	channelID := r.Form.Get("channelID")

	org := r.Header.Get("orgName")
	channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(user), fabsdk.WithOrg(org))
	client, err := ledger.New(channelContext)
	if err != nil {
		log.Fatalf("Failed to create new ledger client: %s", err)
	}

	block, _ := client.QueryConfig(ledger.WithTargetEndpoints(r.URL.Query().Get("peer")))
fmt.Print(block)
	ret, err := json.Marshal(block)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

//QueryConfigBlock remian

func GetPeers(w http.ResponseWriter, r *http.Request) {
	log.Print("================ GetPeers  ======================")
	type response struct {
		Success bool
		Message string
	}
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	// Create SDK setup for channel client with dynamic selection--panic
/*	sdk, _ := fabsdk.New(integration.ConfigBackend,
	fabsdk.WithServicePkg(&dynamicDiscoveryProviderFactory{}))
	defer sdk.Close()

	chProvider := sdk.ChannelContext("insurancechannel", fabsdk.WithUser("Admin"), fabsdk.WithOrg("org1"))
	chCtx, _ := chProvider()
	discoveryService, _ := chCtx.ChannelService().Discovery()
	peers,_:= discoveryService.GetPeers()
*/
	org := r.Form.Get("org")
	ctxProvider := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(org))
	locCtx, _:= contextImpl.NewLocal(ctxProvider)
	peers, _:= locCtx.LocalDiscoveryService().GetPeers()
	fmt.Print("<<<<<<<<<<<<<<<<<<<<<<<<got peers>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	/*for _,peer :=range peers{
		pos :=strings.Index(peer.URL(), "//")
		fmt.Print(pos)
		ret := peer.URL()[pos+2:len(peer.URL())-5]
		fmt.Printf(ret)
		fmt.Print("\n")
	}*/

	peer := peers[0]
	pos :=strings.Index(peer.URL(), "//")
	fmt.Print(pos)
	ret := peer.URL()[pos+2:len(peer.URL())-5]
	fmt.Printf(ret)

	res := response{
		Success: true,
		Message: ret,
		}
	
	ret1, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret1)

}

func IsPeerInChannel(w http.ResponseWriter, r *http.Request) {
	log.Print("================ IsPeerInChannel  ======================")
	type response struct {
		Success bool
		Message string
	}
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	channelID := r.Form.Get("channelID")
	peer := r.Form.Get("peer")

	ctxProvider := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg(orgName))
	ctx, _ := ctxProvider()

	var client *discovery.Client
	client, _  = discovery.New(ctx)
	reqCtx, _:= context.NewRequest(ctx, context.WithTimeout(10*time.Second))

	req := discovery.NewRequest().OfChannel(channelID).AddPeersQuery()

	peerCfg1, _:= comm.NetworkPeerConfig(ctx.EndpointConfig(), peer)

	responses, _:= client.Send(reqCtx, req, peerCfg1.PeerConfig)
	fmt.Print(responses)
	resp := responses[0]
	fmt.Print(reflect.TypeOf(resp))
	fmt.Print("\n<<<<<<<<<<<<<<<<<<<<<<<<resp>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	fmt.Print(resp)
	str := fmt.Sprintf("%v", resp)
	fmt.Println(str)
	
	res := response{
		Success: true,
		Message: "",
	}
	
	pos :=strings.Index(str, "access denied")
	fmt.Print(pos)

	if pos>0{
		res.Success = false
		res.Message = "not in channel"
	}else{
		res.Message = "in channel"
	}

	ret, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

type dynamicDiscoveryProviderFactory struct {
	defsvc.ProviderFactory
}
