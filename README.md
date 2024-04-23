# Health Fabric

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