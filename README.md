# Health Fabric

## Install Hyperledger Fabric
```shell
curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
./install-fabric.sh docker samples binary
```
## ./network.sh down
## How to deploy

### Build the gateway application
```shell
cd application-gateway
go build -o gateway . && mv gateway ../bin/
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

### Read Data
```shell
gateway -u alpha read
gateway -u beta read
gateway -u gama read

gateway -u alpha read beta
gateway -u alpha read gama
```

### Insert diagnosis for beta user by alpha user
```shell
gateway -u alpha diagnosis insert beta glucose 120 bmi 10 insulin 5

gateway -u beta diagnosis insert alpha glucose 100 bmi 5 insulin 4

gateway -u gama diagnosis insert alpha glucose 100 bmi 5 insulin 4
```