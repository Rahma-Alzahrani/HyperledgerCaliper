
---

# Caliper Benchmarks for Hyperledger Fabric

This repository contains benchmark scripts and configurations for evaluating Hyperledger Fabric networks using the Caliper tool. Below are the step-by-step instructions to set up and run these benchmarks.

## Setup Instructions

### 1. Clone the Caliper Benchmarks Repository

```bash
git clone git@github.com:RahmaALzahrani/HyperledgerCaliper.git
cd caliper-benchmarks
```

### 2. Install Caliper Dependencies

```bash
npm install --only=prod @hyperledger/caliper-cli@0.5.0
```

### 3. Bind Caliper with Supported SDK Version for Fabric

```bash
npx caliper bind --caliper-bind-sut fabric:2.2
```

### 4. Configuration Setup

- Ensure Hyperledger Network is running if it's not running use the following command to up the network 
	- Go to fabric network directory run this command './network.sh up createChannel -ca -c mychannel -s couchdb'
	- Dont commit chaincode before running Prometheus server

- Ensure a Prometheus server is running to collect Hyperledger Fabric metrics:
  - Navigate to the `test-network` directory.
  - Bring up the Prometheus and Grafana servers:

    ```bash
    cd test-network/prometheus-grafana
    docker-compose up -d
    ```

- Deploy the modified chaincode to the Fabric network:
  - Navigate to the `test-network` directory:

    ```bash
    cd ..
    ./network.sh deployCC -ccn basic -ccp /home/rahma/Desktop/Caliper/HyperledgerCaliper/src/fabric/chaincode-go -ccl go
    ```
    
    You have to give your path of the chaincode that is present in caliper-benchmark, from the test-network directory where the chaincode file is located pass that 

- Configure the Caliper network setup:
  - Go to the `caliper-benchmarks/networks` directory.
  - Update the `test-network.yaml` file with appropriate paths and certificates as per your fabric network setup.

## Accessing CouchDB for Journey and Offer Details
Access CouchDB for journey and offer details:

Username: admin
Password: adminpw
Access: http://localhost:5984/_utils/#login
Navigate to the mychannel_basic collection to view details of all committed chaincode.

## Generating Reports with Caliper

Once the configuration is set, use the following commands to generate reports:

1. **Create a Journey:**

   ```bash
   npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/fabric/test-network.yaml --caliper-benchconfig benchmarks/datamanagement/journey.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled
   ```

2. **Create Real Data Offer on Journey UID:**

   ```bash
   npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/fabric/test-network.yaml --caliper-benchconfig benchmarks/datamanagement/insert-data-offer.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled
   ```

3. **Insert Data-Hash in the Real Data Offer:**

   ```bash
   npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/fabric/test-network.yaml --caliper-benchconfig benchmarks/datamanagement/insert-data-hash.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled


4. **Create Historical Offer on Journey UID:**

   ```bash
   npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/fabric/test-network.yaml --caliper-benchconfig benchmarks/datamanagement/historical-data-offer.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled
   ```

5. **Insert Data-Hash in the Historical Offer:**

   ```bash
   npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/fabric/test-network.yaml --caliper-benchconfig benchmarks/datamanagement/historical-data-hash.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled
   ```
   
6. **Get All Journey:**

   ```bash
   npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/fabric/test-network.yaml --caliper-benchconfig benchmarks/datamanagement/get-journey.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled
   ```
   
7. **Get All Historical Data Offers:**

	```bash
   npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/fabric/test-network.yaml --caliper-benchconfig benchmarks/datamanagement/get-historical-data-offer.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled
   ```
   
   
8. **Get All Sensor Data Query:**

	```bash
   npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/fabric/test-network.yaml --caliper-benchconfig benchmarks/datamanagement/sensor-query.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled
   ```

9. **Get All Real Data Offer:**

   ```bash
   npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/fabric/test-network.yaml --caliper-benchconfig benchmarks/datamanagement/get-data-offers.yaml --caliper-flow-only-test --caliper-fabric-gateway-enabled
   ```
   
 **After Successfull benchmark you can get the report from the root of the caliper-benchmark folder named as report.html it will keep overwriting the previous reports so don't forget to save pervious reports**  
  

Ensure to verify and validate each step according to your environment setup and configurations before running Caliper's benchmarks for generating reports.

---

Feel free to adjust these instructions according to your specific network setup and requirements.
