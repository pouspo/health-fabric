# Health Fabric

## Install Hyperledger Fabric
```shell
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
./install-fabric.sh docker samples binary
```

## How to deploy

### Build the gateway application
```shell
cd application-gateway
go build -o gateway . && mv gateway /Users/admin/go/src/github.com/health-fabric/bin/
```

### Change directory to test-network
```shell
cd test-network
```

### Bring the hyperledger network live and deploy three chaincodes
```shell
./network.sh up createChannel
./chaincode.sh policy-contract 1.0
./chaincode.sh access-contract 1.0
./chaincode.sh health-contract 1.0
```

### Import dataset from csv
This command first registers the user in the network and then import the health data.
```shell
./importCSVHealdData.sh dataset_v2.csv
```

### Insert some dummy policy
```shell
gateway insert-policy
```

### Insert diagnosis for beta user by alpha user
```shell
gateway -u alpha diagnosis insert beta age 100
```
