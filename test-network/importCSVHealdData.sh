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
    fabric-ca-client register --id.name "$username" --id.secret "$username" --id.type client --id.attrs 'groups=group1-group2:ecert' --tls.certfiles "${PWD}/organizations/fabric-ca/org1/ca-cert.pem"
    { set +x; } 2>/dev/null

    echo "Enrolling the user $username"
    set -x
    fabric-ca-client enroll -u https://"$username":"$username"@localhost:7054 --caname ca-org1 -M "${PWD}/organizations/peerOrganizations/org1.example.com/users/$username@org1.example.com/msp" --tls.certfiles "${PWD}/organizations/fabric-ca/org1/ca-cert.pem"
    { set +x; } 2>/dev/null

    cp "${PWD}/organizations/peerOrganizations/org1.example.com/msp/config.yaml" "${PWD}/organizations/peerOrganizations/org1.example.com/users/$username@org1.example.com/msp/config.yaml"

    echo "Register User"
    gateway -u "$username" register

    echo "Insert Data into ledger"
    gateway -u "$username" diagnosis insert pregnancies "$col1" glucose "$col2" blood_pressure "$col3" skin_thickness "$col4" insulin "$col5" bmi "$col6" diabetes_pedigree_function "$col7" age "$col8" outcome "$col9"
done < "$1"