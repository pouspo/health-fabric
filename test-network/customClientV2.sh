#!/bin/bash

export FABRIC_CA_CLIENT_HOME=${PWD}/organizations/peerOrganizations/org1.example.com/
{ set +x; } 2>/dev/null

# Check if filename is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <csv_file>"
    exit 1
fi

# Check if file exists
if [ ! -f "$1" ]; then
    echo "Error: File $1 not found!"
    exit 1
fi

# Read CSV file
while IFS=',' read -r col1 col2 col3 col4 col5 col6 col7 col8 col9 col10
do
    username=$(echo -n "$col10" | tr -d '\r')
    echo "Register the user $username"
    set -x
    fabric-ca-client register --id.name "$username" --id.secret "$username" --id.type client --id.attrs 'groups=group1:ecert' --tls.certfiles "${PWD}/organizations/fabric-ca/org1/ca-cert.pem"
    { set +x; } 2>/dev/null

    echo "Enrolling the user $username"
    set -x
    fabric-ca-client enroll -u https://"$username":"$username"@localhost:7054 --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/$username@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/ca-cert.pem"
    { set +x; } 2>/dev/null

    cp "${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/org1.example.com/users/$username@org1.example.com/msp/config.yaml"

done < "$1"


#infoln "Register the alpha user"
#set -x
#fabric-ca-client register --id.name alpha --id.secret alphapw --id.type client --id.attrs 'groups=group1-group2:ecert' --tls.certfiles "${PWD}/organizations/fabric-ca/org1/ca-cert.pem"
#{ set +x; } 2>/dev/null
#
#infoln "Enrolling the alpha user"
#set -x
#fabric-ca-client enroll -u https://alpha:alphapw@localhost:7054 --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/alpha@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/ca-cert.pem"
#{ set +x; } 2>/dev/null
#
#cp "${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/org1.example.com/users/alpha@org1.example.com/msp/config.yaml"
#
#infoln "Register the beta user"
#set -x
#fabric-ca-client register --id.name beta --id.secret betapw --id.type client --id.attrs 'groups=group1-group2:ecert' --tls.certfiles "${PWD}/organizations/fabric-ca/org1/ca-cert.pem"
#{ set +x; } 2>/dev/null
#
#infoln "Enrolling the beta user"
#set -x
#fabric-ca-client enroll -u https://beta:betapw@localhost:7054 --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/beta@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/ca-cert.pem"
#{ set +x; } 2>/dev/null
#
#cp "${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/org1.example.com/users/beta@org1.example.com/msp/config.yaml"
#
#
#infoln "Register the gama user"
#set -x
#fabric-ca-client register --id.name gama --id.secret gamapw --id.type client --id.attrs 'groups=group1:ecert' --tls.certfiles "${PWD}/organizations/fabric-ca/org1/ca-cert.pem"
#{ set +x; } 2>/dev/null
#
#infoln "Enrolling the alpha user"
#set -x
#fabric-ca-client enroll -u https://gama:gamapw@localhost:7054 --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/gama@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/ca-cert.pem"
#{ set +x; } 2>/dev/null
#
#cp "${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/org1.example.com/users/gama@org1.example.com/msp/config.yaml"
#
#
