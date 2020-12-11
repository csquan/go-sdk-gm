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
echo
echo "POST invoke chaincode on peers of Org1 and Org2"
echo
TX_INFO=$(curl -s -X POST \
  http://localhost:4000/channels/mychannel/invokechaincodes/insurance \
  -H "authorization:$ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
        \"channelID\":\"insurancechannel\",
	\"peers\": [\"peer0.org1.example.com\",\"peer0.org2.example.com\"],
	\"fcn\":\"createInsurance\",
	\"args\":[\"InsuranceNum002\",\"OrderNum001\",\"Shipper001\",\"Carrier\",\"aisio\",\"TV\",\"200\",\"20000\",\"4000\",\"1\"]
}")
echo $TX_INFO
echo "Total execution time : $(($(date +%s)-starttime)) secs ..."
