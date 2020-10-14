package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http"
	"os"
	"runtime/debug"
	"time"
	_ "time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	_ "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/internal2/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	ledger "github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
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

	sdk, _ = fabsdk.New(config.FromFile(org1CfgPath))
	//if err != nil {
	//	log.Panicf("failed to create fabric sdk: %s", err)
	//}

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
	//createChannel(sdk)
	//getChannels(sdk)
	//joinChannel(sdk)
	//installChainCode(sdk)  //debug ok
	//instantiateChainCode(sdk)
	//getBlockByNumber(sdk)
	//getTransactionByID(sdk, "45db61b8421b8a2b2c6e3afe1dc8f3beb8b16c356d92805d2657c337f86f2912")
	/*ccp := sdk.ChannelContext(channelID, fabsdk.WithUser("Admin"))
	cc, err := channel.New(ccp)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err)
	}

	query(cc)*/
	//execute(cc)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("authorization")
		if tokenString != "" {
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte("thisismysecret"), nil
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

func login(w http.ResponseWriter, r *http.Request) {
	log.Print("==================== LOGIN ==================")
	// define response
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
		ctxProvider := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))
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
		/*if err != ErrUserNotFound {
			//t.Fatal("Expected to not find user")
		}*/
		err = msp.Enroll(userName, mspclient.WithSecret("enrollmentSecret"), mspclient.WithProfile("tls"))
fmt.Println("enroll err")
fmt.Println(err)
		res := response{
			Success: true,
			Message: err.Error(),
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
}

func getInstalledChaincodes(w http.ResponseWriter, r *http.Request) {
}

func getInstantiatedChaincodes(w http.ResponseWriter, r *http.Request) {
}

func createChannel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("================ CREATE CHANNEL ======================")

	clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(ordererOrgName))

	// Resource management client is responsible for managing channels (create/update channel)
	// Supply user that has privileges to create channel (in this case orderer admin)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		//t.Fatalf("Failed to create channel management client: %s", err)
	}

	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
	if err != nil {
		//t.Fatal(err)
	}
	adminIdentity, err := mspClient.GetSigningIdentity(orgAdmin)
	if err != nil {
		//t.Fatal(err)
	}
	req := resmgmt.SaveChannelRequest{ChannelID: channelID,
		ChannelConfigPath: "./mychannel.tx",
		SigningIdentities: []msp.SigningIdentity{adminIdentity}}
	_, err = resMgmtClient.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com"))
	debug.PrintStack()
	var str string
	if err == nil {
		str = "success to create channel"
	} else {
		fmt.Println(err)
		str = err.Error()
	}
	fmt.Println(str)
	out, err := json.Marshal(str)
	w.Write(out)
}
func joinChannel(w http.ResponseWriter, r *http.Request) {
	fmt.Println("================ JOIN CHANNEL ======================")

	//prepare context
	adminContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))

	// Org resource management client
	orgResMgmt, err := resmgmt.New(adminContext)
	if err != nil {
		//t.Fatalf("Failed to create new resource management client: %s", err)
	}

	// Org peers join channel
	if err = orgResMgmt.JoinChannel(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com")); err != nil {
		fmt.Println("Org peers failed to JoinChannel: %s", err)
	}
	var str string
	if err == nil {
		str = "success to join channel"
	} else {
		str = err.Error()
	}

	fmt.Println(str)
	out, err := json.Marshal(str)
	w.Write(out)
}
func getBlockByNumber(w http.ResponseWriter, r *http.Request) {
	fmt.Println("==================== GET BLOCK BY NUMBER ==================")

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
	if err != nil {
		fmt.Printf("QueryConfigBlockFromOrderer return error: %s", err)
	}

	fmt.Printf("GET BLOCK BY NUMBER sucess")
	//	res, err := json.Marshal(ret)
	fmt.Println(ret)
	out, err := json.Marshal(ret)
	w.Write(out)
}
func queryCC(w http.ResponseWriter, r *http.Request) {
	fmt.Println("================query cc======================")
	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser("Admin"))
	cc, err := channel.New(ccp)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err)
	}
	// new channel request for query
	req := channel.Request{
		ChaincodeID: "face",
		Fcn:         "queryAllFace",
		Args:        packArgs([]string{"1111110", "1111112"}),
	}
	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer0Org1)
	fmt.Println("success to reqpeers")
	response, err := cc.Query(req, reqPeers)
	if err != nil {
		fmt.Printf("failed to query chaincode: %s\n", err)
	}
	fmt.Println("success to query cc")
	fmt.Println(string(response.Payload))
	if len(response.Payload) > 0 {
		fmt.Printf("chaincode query success,the value is %s\n", string(response.Payload))
		w.Write(response.Payload)
	}
}

func invokeCC(w http.ResponseWriter, r *http.Request) {
	//args := packArgs([]string{"a", "b", "10"})
	fmt.Println("=====================invoke chaincode====================")
	ccp := sdk.ChannelContext(channelID, fabsdk.WithUser("Admin"))
	cc, err := channel.New(ccp)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err)
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
	if err != nil {
		fmt.Printf("failed to Execute chaincode: %s\n", err)
	}
	fmt.Printf("Execute chaincode success,txId:%s\n", response.TransactionID)
	out, err := json.Marshal(response)
	w.Write(out)
}

func packArgs(paras []string) [][]byte {
	var args [][]byte
	for _, k := range paras {
		args = append(args, []byte(k))
	}
	return args
}
func InstallChainCode(w http.ResponseWriter, r *http.Request) {
	fmt.Println("================install cc======================")

	ccPkg, err := packager.NewCCPackage("sdk_test/src", "/root/gowork")
	if err != nil {
		log.Fatalf("pack chaincode error %s", err)
	}
	// new request of installing chaincode
	req := resmgmt.InstallCCRequest{
		Name:    "face",
		Path:    "sdk_test/src",
		Version: "v0",
		Package: ccPkg,
	}

	clientContext := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("org2"))

	// Resource management client is responsible for managing resources (joining channels, install/instantiate/upgrade chaincodes)
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		log.Fatalf("Failed to create new resource management client")
	}

	reqPeers := resmgmt.WithTargetEndpoints("peer0.org2.example.com")
	_, err = resMgmtClient.InstallCC(req, reqPeers)

	str := ""
	if err == nil {
		str = "success to install cc"
	} else {
		str = err.Error()
	}

	fmt.Println(str)
	out, err := json.Marshal(str)
	w.Write(out)
}

func InstantiateChainCode(w http.ResponseWriter, r *http.Request) {
	fmt.Println("================instantiate cc======================")
	// new request of Instantiate chaincode
	/*type InstantiateCCRequest struct {
		Name       string
		Path       string
		Version    string
		Args       [][]byte
		Policy     *common.SignaturePolicyEnvelope
		CollConfig []*common.CollectionConfig
	}*/
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

	str := ""
	if err == nil {
		str = "success to instantiate cc"
	} else {
		str = err.Error()
	}

	fmt.Println(str)
	out, err := json.Marshal(str)
	w.Write(out)
}

func getChannels(w http.ResponseWriter, r *http.Request) {
	fmt.Println("================ GET CHANNELS ======================")
	// define response
	type response struct {
		Success bool
		Message string
	}
	clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	resMgmtClient, err := resmgmt.New(clientContext)

	if err != nil {
		fmt.Println("Failed to create new resource management client: %s", err)
	}

	//var peers fabsdk.Peer
	//peers.append(peers,peer())
	//org1AdminClientContext := sdk.Context(fabsdk.WithUser("AdminOrg1"), fabsdk.WithOrg("Org1"))
	//org1Peers, err := DiscoverLocalPeers(org1AdminClientContext, 2)

	ret, err := resMgmtClient.QueryChannels(resmgmt.WithTargetEndpoints("peer0.org1.example.com"))

	if err == nil {
		fmt.Printf("GET BLOCK BY NUMBER sucess")
		fmt.Println(ret)
	}
	/*res := response{
		Success: true,
		Message: string(out[:]),
	}

	if err == nil {
		fmt.Printf(out)
	}*/
	//ret, err := json.Marshal(res)
	//w.Header().Set("Content-Type", "application/json")
	out, err := json.Marshal(ret)
	w.Write(out)
}

func getTransactionByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("================ GET TRANSACTION BY TRANSACTION_ID ======================")

	//clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(ordererOrgName))
	//resMgmtClient, err := resmgmt.New(clientContext)
	channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))

	client, err := ledger.New(channelContext)
	if err != nil {
		fmt.Println("Failed to create new channel client: %s", err)
	}
	ret, err := client.QueryTransaction("")
	if err != nil {
		fmt.Println("Failed to queryTx : %s", err)
	}
	fmt.Println("success to getTransactionByID:")
	fmt.Println(ret)
	out, err := json.Marshal(ret)
	w.Write(out)
}
