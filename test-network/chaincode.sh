#!/bin/bash

if [[ $# -ge 1 ]] ; then
  code_folder="$1"
  version="$2"
  # check for the createChannel subcommand
  if [ -n "$code_folder" ] && [ -n "$version" ]; then
    export FABRIC_CFG_PATH=$PWD/../config/

    peer lifecycle chaincode package "$code_folder"_"$version".tar.gz --path ../"$code_folder"/ --lang golang --label "$code_folder"_"$version"

    echo "Setting env(s) for Org1"
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=$PWD/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051

    echo Installing "$code_folder"_"$version" in Org1
    peer lifecycle chaincode install "$code_folder"_"$version".tar.gz

    echo "Setting env(s) for Org2"
    export CORE_PEER_LOCALMSPID="Org2MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PWD/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=$PWD/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    export CORE_PEER_ADDRESS=localhost:9051

    echo Installing "$code_folder"_"$version" in Org2
    peer lifecycle chaincode install "$code_folder"_"$version".tar.gz

    echo "Calculating sequence_number"
    commit=$(peer lifecycle chaincode querycommitted --channelID=mychannel --name="$code_folder")
    # Extract the sequence number using awk
    sequence_number=$(echo "$commit" | awk -F 'Sequence: ' '{print $2}' | awk '{print $1}' | awk '{sub(/,/, ""); print}')
    current_version_number=$(echo "$commit" | awk -F 'Version: ' '{print $2}' | awk '{print $1}' | awk '{sub(/,/, ""); print}')

    if [ -z "$sequence_number" ]; then
        sequence_number=0
    fi

    # Increment the sequence number by 1 to get the next sequence
    next_sequence=$((sequence_number + 1))

    # Print the results
    echo "Current Sequence Number: $sequence_number"
    echo "Next Sequence Number: $next_sequence"
    echo "Version Number: $current_version_number"

    if [ "$current_version_number" = "$version" ]; then
      echo "Version already exists"
        return 1
    fi

    echo "Approving chaincode as Org2"
    regex="$code_folder"_"$version"
    regex+=":[^,]*"
    candidate=$(peer lifecycle chaincode queryinstalled | grep -o "$regex")
    export CC_PACKAGE_ID="$candidate"
    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name "$code_folder" --version "$version" --package-id "$CC_PACKAGE_ID" --sequence "$next_sequence" --tls --cafile "$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

    echo "Setting env(s) for Org1"
    export CORE_PEER_TLS_ENABLED=true
    export CORE_PEER_LOCALMSPID="Org1MSP"
    export CORE_PEER_TLS_ROOTCERT_FILE=$PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
    export CORE_PEER_MSPCONFIGPATH=$PWD/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    export CORE_PEER_ADDRESS=localhost:7051

    echo "Approving chaincode as Org2"
    peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name "$code_folder" --version "$version" --package-id "$CC_PACKAGE_ID" --sequence "$next_sequence" --tls --cafile "$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

    echo "Commit Chaincode"
    peer lifecycle chaincode checkcommitreadiness --channelID mychannel --name "$code_folder" --version "$version" --sequence "$next_sequence" --tls --cafile "$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --output json

    peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID mychannel --name "$code_folder" --version "$version" --sequence "$next_sequence" --tls --cafile "$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" --peerAddresses localhost:7051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"

    shift
  else
    # Sample output from the command
    output="Version: 1.1, Sequence: 2, Endorsement Plugin: escc, Validation Plugin: vscc, Approvals: [Org1MSP: true, Org2MSP: true]"

    # Extract the sequence number using awk
    sequence_number=$(echo "$output" | awk -F 'Sequence: ' '{print $2}' | awk '{print $1}' | awk '{sub(/,/, ""); print}')

    if [ -z "$sequence_number" ]; then
        sequence_number=0
    fi

    echo "$sequence_number"
    # Increment the sequence number by 1 to get the next sequence
    next_sequence=$((sequence_number + 1))

    # Print the results
    echo "Current Sequence Number: $sequence_number"
    echo "Next Sequence Number: $next_sequence"
    echo "Incomplete command"
    echo "-code_folder -version"
  fi
fi

# peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n policy-contract --peerAddresses localhost:7051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"AddPolicy","Args":["admin", "user_101", "read", "field_1"]}'
# peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "$PWD/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n policy-contract --peerAddresses localhost:7051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "$PWD/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"ReadPolicy","Args":["admin"]}'
# peer chaincode query -C mychannel -n policy-contract -c '{"Args":["ReadPolicy","admin"]}'