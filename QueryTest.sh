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
	  "http://localhost:4000/users?username=jim&orgName=org1&secret=123")
echo $ORG1_TOKEN
ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".Token" | sed "s/\"//g")
echo
echo "ORG1 token is $ORG1_TOKEN"
echo

echo "GET query Block by blockNumber"
echo
BLOCK_INFO=$(curl -s -X GET \
	  "http://localhost:4000/channels/mychannel/blocks?peer=peer0.org1.example.com&channelID=cargochannel&blockID=4" \
	   -H "authorization:$ORG1_TOKEN" \
	   -H "content-type: application/json")
echo $BLOCK_INFO

#echo "GET query Transaction by TransactionID"
#echo $TRX_ID
#curl -s -X GET http://localhost:4000/channels/mychannel/transactions/$TRX_ID?peer=peer0.org1.example.com \
#  -H "authorization:$ORG1_TOKEN" \
#  -H "content-type: application/json"
#echo
#echo



echo "GET query ChainInfo"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel?peer=peer0.org1.example.com&channelID=cargochannel" \
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
  "http://localhost:4000/channels/mychannel/chaincodes?peer=peer0.org1.example.com&channelID=cargochannel" \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo

echo "GET query Channels"
echo
curl -s -X GET \
  "http://localhost:4000/channels?peer=peer0.org1.example.com&org=org1" \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json"
echo
echo


echo "Total execution time : $(($(date +%s)-starttime)) secs ..."

