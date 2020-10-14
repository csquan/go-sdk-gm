package main

import (
	_ "encoding/json"
	"fmt"
	"log"
	_ "net/http"
	"runtime/debug"
	_ "time"

	_ "github.com/gorilla/mux"
	_ "github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/internal2/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	ledger "github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	fabAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
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

func main() {
	sdk, err := fabsdk.New(config.FromFile(org1CfgPath))
	if err != nil {
		log.Panicf("failed to create fabric sdk: %s", err)
	}

	//createChannel(sdk)
	//getChannels(sdk)
	//joinChannel(sdk)
	//installChainCode(sdk)  //debug ok
	//instantiateChainCode(sdk)
	//getBlockByNumber(sdk)
	getTransactionByID(sdk, "45db61b8421b8a2b2c6e3afe1dc8f3beb8b16c356d92805d2657c337f86f2912")
	/*ccp := sdk.ChannelContext(channelID, fabsdk.WithUser("Admin"))
	cc, err := channel.New(ccp)
	if err != nil {
		log.Panicf("failed to create channel client: %s", err)
	}

	query(cc)*/
	//execute(cc)
}

func createChannel(sdk *fabsdk.FabricSDK) {
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
}
func joinChannel(sdk *fabsdk.FabricSDK) {
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
}
func getBlockByNumber(sdk *fabsdk.FabricSDK) {
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
}
func query(cc *channel.Client) {
	fmt.Println("================query cc======================")
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
	}
}

func execute(cc *channel.Client) {
	//args := packArgs([]string{"a", "b", "10"})
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
}

func packArgs(paras []string) [][]byte {
	var args [][]byte
	for _, k := range paras {
		args = append(args, []byte(k))
	}
	return args
}
func installChainCode(sdk *fabsdk.FabricSDK) {
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

}

func instantiateChainCode(sdk *fabsdk.FabricSDK) {
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

}

// DiscoverLocalPeers queries the local peers for the given MSP context and returns all of the peers. If
// the number of peers does not match the expected number then an error is returned.
func DiscoverLocalPeers(ctxProvider contextAPI.ClientProvider, expectedPeers int) ([]fabAPI.Peer, error) {
	ctx, err := contextImpl.NewLocal(ctxProvider)
	if err != nil {
		return nil, errors.Wrap(err, "error creating local context")
	}

	discoveredPeers, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			peers, serviceErr := ctx.LocalDiscoveryService().GetPeers()
			if serviceErr != nil {
				return nil, errors.Wrapf(serviceErr, "error getting peers for MSP [%s]", ctx.Identifier().MSPID)
			}
			if len(peers) < expectedPeers {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Expecting %d peers but got %d", expectedPeers, len(peers)), nil)
			}
			return peers, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return discoveredPeers.([]fabAPI.Peer), nil
}

func getChannels(sdk *fabsdk.FabricSDK) {
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
	//sssw.Write(ret)

}

func getTransactionByID(sdk *fabsdk.FabricSDK, txId fab.TransactionID) {
	fmt.Println("================ GET TRANSACTION BY TRANSACTION_ID ======================")

	//clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(ordererOrgName))
	//resMgmtClient, err := resmgmt.New(clientContext)
	channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser("Admin"), fabsdk.WithOrg("Org1"))

	client, err := ledger.New(channelContext)
	if err != nil {
		fmt.Println("Failed to create new channel client: %s", err)
	}

	ret, err := client.QueryTransaction(txId)
	if err != nil {
		fmt.Println("Failed to queryTx : %s", err)
	}
	fmt.Println("success to getTransactionByID:")
	fmt.Println(ret)
}
