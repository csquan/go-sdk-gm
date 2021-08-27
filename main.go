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
	"encoding/base64"
	"reflect"
	"strconv"
	"strings"
	"time"
	//"io/ioutil"

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
	//contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
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

	//"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	lutil "github.com/hyperledger/fabric/common/ledger/util"
	//putil "github.com/hyperledger/fabric/protos/utils"
	"github.com/golang/protobuf/proto"
	putil "github.com/hyperledger/fabric-sdk-go/internal2/github.com/hyperledger/fabric/protoutil"

	"database/sql"

	"github.com/golang/protobuf/ptypes"
	_ "github.com/lib/pq"
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

var channel_genesis_hash = "573f3ff5686831e322cb1c02769ebd5519ec7b3618cabcc1dca1705a6a7e1808"

var sdk *fabsdk.FabricSDK

var db *sql.DB

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func sqlOpen() {
	var err error

	db, err = sql.Open("postgres", "port=5432 user=hppoc password=password dbname=fabricexplorer sslmode=disable")
	checkErr(err)
}

func peerSqlInsert(org uint64,channel_genesis_hash string, mspid string, requests string, server_hostname string, createdt string,peer_type string) {
	 //插入数据
	stmt, err := db.Prepare("INSERT INTO peer (org,channel_genesis_hash,mspid,requests,server_hostname,createdt,peer_type) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id")
	checkErr(err)
	res, err := stmt.Exec(org,channel_genesis_hash,mspid,requests,server_hostname,createdt,peer_type)
	checkErr(err)

	_, err = res.RowsAffected()
	checkErr(err)
}

func sqlDelete(name string) {
	//删除数据
	stmt, err := db.Prepare("delete from chaincodes where name=$1")
	checkErr(err)

	res, err := stmt.Exec(name)
	checkErr(err)

	_, err = res.RowsAffected()
	checkErr(err)
}

func chaincodeSqlDelandInsert(name string ,version string,path string,txcount int ,channel_genesis_hash string){

       sqlDelete("pingan")

       stmt, err := db.Prepare("INSERT INTO chaincodes(name,version,path,txcount,channel_genesis_hash) VALUES($1,$2,$3,$4,$5) RETURNING id")
       checkErr(err)
       res, err := stmt.Exec(name,version,path,txcount,channel_genesis_hash)
       checkErr(err)

       _, err = res.RowsAffected()
       checkErr(err)
 }

func peerSqlSelect() int{
	//查询数据
        rows, err := db.Query("SELECT *  FROM peer;")
	checkErr(err)

	count:=0
	for rows.Next(){
		count++
	}
	return count
}

func channelSqlSelect() int{
	//查询数据
	rows, err := db.Query("SELECT *  FROM channel;")
	checkErr(err)

	count:=0
	for rows.Next(){
		count++
	}
	return count
}

func chaincodeSqlGetTxCount(name string) int {
	//查询数据
	sql := "SELECT txcount  FROM chaincodes where name = '" + name + "';";
	rows, err := db.Query(sql)
	checkErr(err)


	count := 0
	for rows.Next() {
		err = rows.Scan(&count)
		checkErr(err)
	}

	return count
}

func txSqlInsert(blockid uint64,txhash string, createdt string, creator_msp_id string, chaincodename string, channel_genesis_hash string) {
	//插入数据
	stmt, err := db.Prepare("INSERT INTO transactions(blockid,txhash,createdt,creator_msp_id,chaincodename,channel_genesis_hash) VALUES($1,$2,$3,$4,$5,$6) RETURNING id")
	checkErr(err)
	res, err := stmt.Exec(blockid,txhash, createdt, creator_msp_id,chaincodename, channel_genesis_hash)
	checkErr(err)

	_, err = res.RowsAffected()
	checkErr(err)
}

func blocksSqlInsert(blocknum int, txcount int, createdt string, prev_blockhash string, blockhash string, channel_genesis_hash string) {
	stmt, err := db.Prepare("INSERT INTO blocks(blocknum,txcount, createdt,prev_blockhash,blockhash,channel_genesis_hash) VALUES($1,$2,$3,$4,$5,$6) RETURNING id")
	checkErr(err)

	fmt.Printf(string(txcount))
	res, err := stmt.Exec(blocknum, txcount, createdt, prev_blockhash, blockhash, channel_genesis_hash)
	checkErr(err)

	_, err = res.RowsAffected()
	checkErr(err)
}

func channelSqlInsert(name  string, blocks int, trans int,channel_genesis_hash string){
      stmt, err := db.Prepare("INSERT INTO channel(name, blocks, trans,channel_genesis_hash) VALUES($1,$2,$3,$4) RETURNING id")
      checkErr(err)

      res, err := stmt.Exec(name, blocks,trans,channel_genesis_hash)
      checkErr(err)

      _, err = res.RowsAffected()
      checkErr(err)
}

func blocksSqlUpdate(blocknum int, txcount int){
     //更新数据
     stmt, err := db.Prepare("update blocks set txcount=$1 where blocknum=$2")
     checkErr(err)

     res, err := stmt.Exec(txcount, blocknum)
     checkErr(err)

     _, err = res.RowsAffected()
     checkErr(err)
}

func blocksSqlSelect() (int,int){
	//查询数据
	rows, err := db.Query("SELECT blocknum ,txcount FROM blocks order by blocknum desc;")
	checkErr(err)

	for rows.Next() {
		var blocknum int
		var txcount int
		err = rows.Scan(&blocknum,&txcount)
		checkErr(err)
		return blocknum,txcount
	}
	return -1,-1
}


func main() {

	sqlOpen()

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
	r.Handle("/channels", http.HandlerFunc(getChannels)).Methods("GET")
	syncBlocks()
	http.ListenAndServe(":4000", handlers.LoggingHandler(os.Stdout, r))
}

func syncBlocks() {
	height := getblocksHeight("mychannel", "org1")

	dbblocks,_ := blocksSqlSelect()

	fmt.Print("\n blocknetwork height:")
	fmt.Print(height)
	fmt.Print("\n db height:")
	fmt.Print(dbblocks)

	start := 0
	max   := 0
	if(int(height) > dbblocks){
            start = dbblocks+1
	    max = int(height)
	}else if(int(height) < dbblocks){
	    start = int(height)
	    max = dbblocks
	}
	fmt.Print("\n sync history blocks from %d to %d",start,max)

	txcounts := 0
	for ; start <= int(height); start++ {
		txcounts = txcounts + handleBlockByNumber("mychannel","pingan",uint64(start))
	}

	channelnumber := channelSqlSelect()
	if channelnumber == 0{
		//insert channel info into table
		channelSqlInsert("mychannel",max-1,txcounts,channel_genesis_hash)
	}

	number := peerSqlSelect()
	if number == 0{
		fmt.Print("\n+insert peer+")
		peerSqlInsert(1,"573f3ff5686831e322cb1c02769ebd5519ec7b3618cabcc1dca1705a6a7e1808","Org1MSP","grpcs://peer0.org1.example.com:7051","peer0.org1.example.com","2021-7-26","PEER")

		peerSqlInsert(1,"573f3ff5686831e322cb1c02769ebd5519ec7b3618cabcc1dca1705a6a7e1808","Org1MSP","grpcs://peer1.org1.example.com:8051","peer1.org1.example.com","2021-7-26","PEER")

		peerSqlInsert(2,"573f3ff5686831e322cb1c02769ebd5519ec7b3618cabcc1dca1705a6a7e1808","Org2MSP","grpcs://peer0.org2.example.com:9051","peer0.org2.example.com","2021-7-26","PEER")

		peerSqlInsert(2,"573f3ff5686831e322cb1c02769ebd5519ec7b3618cabcc1dca1705a6a7e1808","Org2MSP","grpcs://peer1.org2.example.com:10051","peer1.org2.example.com","2021-7-26","PEER")

		peerSqlInsert(3,"573f3ff5686831e322cb1c02769ebd5519ec7b3618cabcc1dca1705a6a7e1808","OrderMSP","grpcs://orderer.example.com:7050","orderer.example.com","2021-7-26","ORDERER")
		peerSqlInsert(3,"573f3ff5686831e322cb1c02769ebd5519ec7b3618cabcc1dca1705a6a7e1808","OrderMSP","grpcs://orderer2.example.com:8050","orderer2.example.com","2021-7-26","ORDERER")
		peerSqlInsert(3,"573f3ff5686831e322cb1c02769ebd5519ec7b3618cabcc1dca1705a6a7e1808","OrderMSP","grpcs://orderer3.example.com:9050","orderer3.example.com","2021-7-26","ORDERER")
		peerSqlInsert(3,"573f3ff5686831e322cb1c02769ebd5519ec7b3618cabcc1dca1705a6a7e1808","OrderMSP","grpcs://orderer4.example.com:10050","orderer4.example.com","2021-7-26","ORDERER")
		peerSqlInsert(3,"573f3ff5686831e322cb1c02769ebd5519ec7b3618cabcc1dca1705a6a7e1808","OrderMSP","grpcs://orderer5.example.com:11050","orderer5.example.com","2021-7-26","ORDERER")
	}
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
		ctxProvider := sdk.Context(fabsdk.WithUser("Admin"), fabsdk.WithOrg("org1"))
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
		/*if err != nil {
			log.Printf("enroll err：%s", err.Error())
			res.Success = false
			res.Message = err.Error()
		}*/
		res.Token = tokenString
		out, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func getblocksHeight(channelID string, orgName string) int64 {

	user := "Admin"

	channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(user), fabsdk.WithOrg(orgName))
	client, err := ledger.New(channelContext)
	if err != nil {
		log.Fatalf("Failed to create new ledger client: %s", err)
		return 0
	}

	blockchainInfo, _ := client.QueryInfo(ledger.WithTargetEndpoints("peer0.org1.example.com"))

	return int64(blockchainInfo.BCI.Height)
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
	log.Print(r)
	decoder := json.NewDecoder(r.Body)
	channel := channelConfig{}
	decoder.Decode(&channel)
	log.Print(channel)
	clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg("org1"))

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
		Org       string
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
	err = orgResMgmt.JoinChannel(join.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint("orderer.example.com"))

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
		Peers     []string
		Fcn       string
		Args      []string
	}
	channel_genesis_hash := "573f3ff5686831e322cb1c02769ebd5519ec7b3618cabcc1dca1705a6a7e1808"
	decoder := json.NewDecoder(r.Body)
	invoke := invokeConfig{}
	decoder.Decode(&invoke)

	log.Print(invoke)
	log.Print(invoke.ChannelID)
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
	timelocal:= time.FixedZone("CST", 8*3600) //translate to local time
	time.Local = timelocal
	t := time.Now()
	str := t.Format("2006-01-02 15:04:05")
	fmt.Printf("cur time:+++")
	fmt.Printf(str)
	response, err := cc.Execute(req, reqPeers)

	res := response1{
		Success: true,
		Message: "",
		TxID:    "",
	}
	if err != nil {
		res.Success = false
		res.Message = err.Error()
		log.Printf("failed to Execute chaincode: %s\n", err.Error())
	} else {
		log.Print("+++++++++++++++response++++++++++++++")
		res.Message = string(response.Payload)
		res.TxID = string(response.TransactionID)

		block :=queryBlockByTXID("mychannel", res.TxID)
		txSqlInsert(block.Header.Number,res.TxID, str, "Org1MSP",vars["chaincodeName"], channel_genesis_hash)

		dbheight,txCount := blocksSqlSelect()
		fmt.Print("++++++++++++height and txCountin db++++++++++++")
		fmt.Print(dbheight)
		fmt.Print(txCount)

		blockheight := int(block.Header.Number)


		if(blockheight > dbheight){
			blocksSqlInsert(blockheight, 1, str, base64.StdEncoding.EncodeToString(block.Header.PreviousHash), base64.StdEncoding.EncodeToString(block.Header.DataHash), channel_genesis_hash)
		}else if(blockheight == dbheight){
                        blocksSqlUpdate(blockheight,txCount+1) 
		}else{
		     fmt.Print("####################error blockheight and dbheight#####################\n")
		     fmt.Print("blockheight and dbheight \n")
		     fmt.Print(blockheight)
		     fmt.Print(dbheight)
		}
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
		Name      string
		Version   string
		Path      string
		Args      []string
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
		Name      string
		Version   string
		Path      string
		Args      string
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
		/*out, err1 := json.Marshal(ret)
		if err1 == nil {
			res.Message = string(out[:])
		}*/
	} else {
		res.Success = false
		res.Message = err.Error()
	}
	w.Header().Set("Content-Type", "application/json")
	//out, err := json.Marshal(res)
	out, err := json.Marshal(ret)
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
	block, _ := client.QueryBlockByHash([]byte(blockHash), ledger.WithTargetEndpoints(peer))

	ret, err := json.Marshal(block)
	log.Print("================ GET BLOCK BY HASH Start======================")
	fmt.Print(block)
	fmt.Print(blockHash)
	log.Print("================ GET BLOCK BY HASH End======================")
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

// GetPayload Get Payload from Envelope message
func GetPayload(e *common.Envelope) (*common.Payload, error) {
	payload := &common.Payload{}
	_ = proto.Unmarshal(e.Payload, payload)
	return payload, nil
}

func queryBlockByTXID(channelID string, txid string) *common.Block{
	log.Print("================ queryBlockByTXID ======================")
	// define response
	type response struct {
		Success bool
		Message string
	}

	channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser("Admin"), fabsdk.WithOrg("org1"))
	client, err := ledger.New(channelContext)
	if err != nil {
		log.Fatalf("Failed to create new ledger client: %s", err)
	}
	block, _ := client.QueryBlockByTxID(fab.TransactionID(txid), ledger.WithTargetEndpoints("peer0.org1.example.com"))

        return block
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
	block, _ := client.QueryBlockByTxID(fab.TransactionID(txid), ledger.WithTargetEndpoints(r.URL.Query().Get("peer")))

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

	res := response{
		Success: true,
		Message: "",
	}

	client, err := ledger.New(channelContext)
	if err != nil {
		log.Print("Failed to create new ledger client: %s", err)
		res.Success = false
		res.Message = err.Error()

		out1, _ := json.Marshal(res)
		w.Write(out1)
		return
	}
	txid := r.Form.Get("txid")
	tx, err := client.QueryTransaction(fab.TransactionID(txid), ledger.WithTargetEndpoints(r.URL.Query().Get("peer")))

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

func extractData(buf *lutil.Buffer) (*common.BlockData, error) {
	data := &common.BlockData{}
	var numItems uint64
	var err error

	if numItems, err = buf.DecodeVarint(); err != nil {
		return nil, err
	}
	for i := uint64(0); i < numItems; i++ {
		var txEnvBytes []byte
		if txEnvBytes, err = buf.DecodeRawBytes(false); err != nil {
			return nil, err
		}
		data.Data = append(data.Data, txEnvBytes)
	}
	return data, nil
}

func extractTxID(txEnvelopBytes []byte) (string, error) {
	txEnvelope, err := putil.GetEnvelopeFromBlock(txEnvelopBytes)
	if err != nil {
		return "", err
	}
	txPayload, err := GetPayload(txEnvelope)
	if err != nil {
		return "", nil
	}
	chdr, err := putil.UnmarshalChannelHeader(txPayload.Header.ChannelHeader)
	if err != nil {
		return "", err
	}
	return chdr.TxId, nil
}

//insert history blocks and txs into db ,return txs in block
func handleBlock(block *common.Block,chaincodeName string) int{
	var str string

	if block == nil {
		return 0
	}
	txCount := 0
	if putil.IsConfigBlock(block) {
		fmt.Printf("txid=CONFIGBLOCK\n")
		return 0
	} else {
		for _, txEnvBytes := range block.GetData().GetData() {
			txEnvelopeBytes := block.Data.Data[0]
			txEnvelope, err := putil.GetEnvelopeFromBlock(txEnvelopeBytes)

			payload := &common.Payload{}
			err = proto.Unmarshal(txEnvelope.Payload, payload)
			if err != nil {
				return 0
			}

			channelHeader := &common.ChannelHeader{}
			err = proto.Unmarshal(payload.Header.ChannelHeader, channelHeader)
			if err != nil {
				return 0
			}
			t, err := ptypes.Timestamp(channelHeader.Timestamp)
			if err != nil {
				fmt.Println(err)
			}
                        m, _ := time.ParseDuration("8h")
			t = t.Add(m)

			str = t.Format("2006-01-02 15:04:05")
			if txid, err := extractTxID(txEnvBytes); err != nil {
				fmt.Printf("ERROR: Cannot extract txid, error=[%v]\n", err)
				return 0
			} else {
				//fmt.Printf("  handleBlock get txid=%s\n", txid)
				txCount++

				txSqlInsert(block.GetHeader().Number,txid, str, "Org1MSP",chaincodeName, channel_genesis_hash)
			}
		}
		//get txcount from chaincodes table
		oldtxcount := chaincodeSqlGetTxCount(chaincodeName)
		/*fmt.Print("table chaincodes already had ")
		fmt.Print(oldtxcount)
		fmt.Print(" txs,update to ")
		fmt.Print(oldtxcount + txCount)
		fmt.Print("\n")*/
		chaincodeSqlDelandInsert(chaincodeName,"v1.0","github/pingan",oldtxcount + txCount,channel_genesis_hash)
	}
	blocksSqlInsert(int(block.GetHeader().Number), txCount, str, base64.StdEncoding.EncodeToString(block.GetHeader().PreviousHash), base64.StdEncoding.EncodeToString(block.GetHeader().DataHash), channel_genesis_hash)
	return txCount
}

func getBlockByNumber(w http.ResponseWriter, r *http.Request) {
	log.Print("================ getBlockByNumber ======================")
	// define response
	type response struct {
		Success bool
		Message string
	}
	res := response{
		Success: true,
		Message: "",
	}

	user := "Admin"
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	channelID := r.Form.Get("channelID")
	blockID := r.Form.Get("blockID")
	log.Print(blockID)
	blockNumber, err := strconv.ParseUint(blockID, 10, 64)
	if err != nil {
		panic(err)
	}

	org := r.Header.Get("orgName")
	channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(user), fabsdk.WithOrg(org))
	client, err := ledger.New(channelContext)
	if err != nil {
		log.Print("Failed to create new ledger client: %s", err)
		res.Success = false
		res.Message = err.Error()

		out1, _ := json.Marshal(res)
		w.Write(out1)
		return
	}
	block, _ := client.QueryBlock(blockNumber, ledger.WithTargetEndpoints(r.URL.Query().Get("peer")))

	//handleBlock(block)
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
	pos := strings.IndexAny(block.Data.String(), "@")
	txid, err := json.Marshal(block.Data.String()[pos+1 : pos+65])

	res.Success = true
	res.Message = "txid:" + string(txid)

	ret, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

func handleBlockByNumber(channelID string,chaincodeName string,blockNumber uint64) int {

	channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser("Admin"), fabsdk.WithOrg("org1"))
	client, err := ledger.New(channelContext)
	if err != nil {
		return 0
	}

	block, _ := client.QueryBlock(blockNumber, ledger.WithTargetEndpoints("peer0.org1.example.com"))

	txcount := handleBlock(block,chaincodeName)    //insert history blocks and txs into db
	return txcount
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

	channelID := r.Form.Get("channelID")
	org := r.Form.Get("org")

	chProvider := sdk.ChannelContext(channelID, fabsdk.WithUser("Admin"), fabsdk.WithOrg(org))
	chCtx, _ := chProvider()
	discoveryService, _ := chCtx.ChannelService().Discovery()
	peers1, err := discoveryService.GetPeers()
	ret := ""
	if err == nil {
		for _, peer1 := range peers1 {
			fmt.Print(peer1)
			fmt.Print("\n")
			ret = ret + peer1.URL()
		}
	}

	/*org := r.Form.Get("org")
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
	}

	peer := peers[0]
	pos :=strings.Index(peer.URL(), "//")
	fmt.Print(pos)
	ret := peer.URL()[pos+2:len(peer.URL())-5]
	fmt.Printf(ret)
	*/
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
	client, _ = discovery.New(ctx)
	reqCtx, _ := context.NewRequest(ctx, context.WithTimeout(10*time.Second))

	req := discovery.NewRequest().OfChannel(channelID).AddPeersQuery()

	peerCfg1, _ := comm.NetworkPeerConfig(ctx.EndpointConfig(), peer)

	responses, _ := client.Send(reqCtx, req, peerCfg1.PeerConfig)
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

	pos := strings.Index(str, "access denied")
	fmt.Print(pos)

	if pos > 0 {
		res.Success = false
		res.Message = "not in channel"
	} else {
		res.Message = "in channel"
	}

	ret, err := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(ret)
}

type dynamicDiscoveryProviderFactory struct {
	defsvc.ProviderFactory
}
