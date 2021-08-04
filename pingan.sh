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
		CC_SRC_PATH="./src"
		;;
		"node")
		CC_SRC_PATH="./src"
		;;
		*) printf "\n ------ Language $LANGUAGE is not supported yet ------\n"$
		exit 1
	esac
}

setChaincodePath

echo "POST request Enroll on Org1  ..."
echo
ORG1_TOKEN=$(curl -s -X GET \
  "http://localhost:4000/users?username=jim1&orgName=org1&secret=123")
echo $ORG1_TOKEN
ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".Token" | sed "s/\"//g")
echo
echo "ORG1 token is $ORG1_TOKEN"
echo
echo "POST request Create channel  ..."
echo
curl -s -X POST \
  http://localhost:4000/channels \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"name":"mychannel",
	"path":"./artifacts/channel2/channel.tx",
	"org":"org1"
}'
echo 
echo
echo "POST request Join channel on Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"channelID":"mychannel",
	"org": "org1"
}'
echo
echo
echo
echo "POST request Join channel on Org2"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/peers \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"channelID":"mychannel",
	"org": "org2"
}'
echo
echo
echo

echo "POST Install chaincode on peer0.Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/installchaincodes?peer="peer0.org1.example.com" \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"name\":\"pingan\",
	\"path\":\"github.com/pingan\",
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
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"name\":\"pingan\",
	\"path\":\"github.com/pingan\",
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
        \"channelID\":\"mychannel\",
	\"name\":\"pingan\",
	\"version\":\"v0\",
	\"path\":\"github.com/go\",
	\"args\":[]
}"
echo
sleep 5
echo
echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TX_INFO=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/invokechaincodes/pingan \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
        \"channelID\":\"mychannel\",
	\"peers\": [\"peer0.org1.example.com\",\"peer0.org2.example.com\"],
	\"fcn\":\"PersonalRegister\",
	\"args\":[\"orderid001\",\"user001\",\"secoouserxxx\",\"13588888888\",\"0\",\"0\",\"CertNo001\",\"xxx@163.com\",\"remark\",\"acctno001\",\"merchaintid001\",\"2020-11-18\",\"2020-11-19\"]}")
echo $TX_INFO
echo "Total execution time : $(($(date +%s)-starttime)) secs ..."
