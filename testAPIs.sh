#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

jq --version > /dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
	echo
	exit 1
fi

starttime=$(date +%s)

# Print the usage message
function printHelp () {
  echo "Usage: "
  echo "  ./testAPIs.sh -l golang|node"
  echo "    -l <language> - chaincode language (defaults to \"golang\")"
}
# Language defaults to "golang"
LANGUAGE="golang"

# Parse commandline args
while getopts "h?l:" opt; do
  case "$opt" in
    h|\?)
      printHelp
      exit 0
    ;;
    l)  LANGUAGE=$OPTARG
    ;;
  esac
done

##set chaincode path
function setChaincodePath(){
	LANGUAGE=`echo "$LANGUAGE" | tr '[:upper:]' '[:lower:]'`
	case "$LANGUAGE" in
		"golang")
		CC_SRC_PATH="github.com/example_cc"
		;;
		"node")
		CC_SRC_PATH="$PWD/artifacts/src/github.com/example_cc/node"
		;;
		*) printf "\n ------ Language $LANGUAGE is not supported yet ------\n"$
		exit 1
	esac
}

setChaincodePath

echo "POST request Enroll on Org1  ..."
echo
ORG1_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=jim&orgName=org1&secret=123')
echo $ORG1_TOKEN
ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".Token" | sed "s/\"//g")
echo
echo "ORG1 token is $ORG1_TOKEN"
echo
echo "POST request Enroll on Org2 ..."
echo
ORG2_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=barry&orgName=org2&secret=123')
echo $ORG2_TOKEN
ORG2_TOKEN=$(echo $ORG2_TOKEN | jq ".Token" | sed "s/\"//g")
echo
echo "ORG2 token is $ORG2_TOKEN"
echo
echo
echo "POST request Create channel  ..."
echo
curl -s -X POST \
  http://localhost:4000/channels \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"name":"mychannel",
	"path":"/root/gowork/src/github.com/balance-transfer-go/artifacts/channel/mychannel.tx",
	"org":"org1"
}'
echo
echo
sleep 5
echo "POST request Join channel on Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"org": "org1"
}'
echo
echo

echo "POST request Join channel on Org2"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization:$ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"org": "org2"
}'
echo
echo


echo "POST Install chaincode on peer0.Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/installchaincodes?peer="peer0.org1.example.com" \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"name\":\"face\",
	\"path\":\"github.com/face\",
	\"version\":\"v0\",
	\"org\":\"org1\",
	\"user\":\"jim\"
}"
echo
echo

echo "POST Install chaincode on peer1.Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/installchaincodes?peer="peer1.org1.example.com" \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"name\":\"face\",
	\"path\":\"github.com/face\",
	\"version\":\"v0\",
	\"org\":\"org1\",
	\"user\":\"jim\"
}"
echo
echo

echo "POST Install chaincode on peer0.Org2"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/installchaincodes?peer="peer0.org2.example.com" \
  -H "authorization:$ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d "{
    \"name\":\"face\",
    \"path\":\"github.com/face\",
    \"version\":\"v0\",
    \"org\":\"org2\",
    \"user\":\"barry\"
}"
echo
echo

echo "POST Install chaincode on peer1.Org2"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/installchaincodes?peer="peer1.org2.example.com" \
  -H "authorization:$ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d "{
    \"name\":\"face\",
    \"path\":\"github.com/face\",
    \"version\":\"v0\",
    \"org\":\"org2\",
    \"user\":\"barry\"
}"
echo
echo

echo "POST instantiate chaincode on Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/instantiatechaincodes \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"name\":\"face\",
	\"version\":\"v0\",
	\"path\":\"github.com/go\",
	\"args\":[]
}"
echo

sleep 5

echo "POST invoke chaincode on peers of Org1 and Org2"
echo
curl -s -X POST \
	  http://localhost:4000/channels/mychannel/invokechaincodes/face \
	    -H "authorization:$ORG1_TOKEN" \
	      -H "content-type: application/json" \
	        -d "{
        \"peers\": [\"peer0.org1.example.com\",\"peer0.org2.example.com\"],
	        \"fcn\":\"createFace\",
		        \"args\":[\"1111111\",\"2222222\",\"33333333\",\"xxxx\"]
		}"
echo 
echo


echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TX_INFO=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/invokechaincodes/face \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.org1.example.com\",\"peer0.org2.example.com\"],
	\"fcn\":\"createFace\",
	\"args\":[\"1111111\",\"2222222\",\"33333333\",\"xxxx\"]
}")
echo $TX_INFO
echo
TRX_ID=$(echo $TX_INFO | jq -r ".Proposal.TxnID")

echo "GET query Block by blockNumber"
echo
BLOCK_INFO=$(curl -s -X GET \
	  "http://localhost:4000/channels/mychannel/blocks/1?peer=peer0.org1.example.com" \
	   -H "authorization:$ORG1_TOKEN" \
	   -H "content-type: application/json")
echo $BLOCK_INFO
# Assign previvious block hash to HASH
HASH=$(echo $BLOCK_INFO | jq -r ".header.previous_hash")
echo

sleep 5

echo "GET query Transaction by TransactionID"
echo $TRX_ID
curl -s -X GET http://localhost:4000/channels/mychannel/transactions/$TRX_ID?peer=peer0.org1.example.com \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo



echo "GET query ChainInfo"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel?peer=peer0.org1.example.com" \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

echo "GET query Installed chaincodes"
echo
curl -s -X GET \
  "http://localhost:4000/chaincodes?peer=peer0.org1.example.com" \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

echo "GET query Instantiated chaincodes"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes?peer=peer0.org1.example.com" \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

echo "GET query Channels"
echo
curl -s -X GET \
  "http://localhost:4000/channels?peer=peer0.org1.example.com" \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo


echo "Total execution time : $(($(date +%s)-starttime)) secs ..."

