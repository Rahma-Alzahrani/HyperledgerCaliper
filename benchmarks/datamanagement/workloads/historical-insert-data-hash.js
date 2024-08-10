'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const { randomBytes } = require('crypto');

class InsertDataHashWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.chaincodeID = 'basic'; // Replace with your chaincode ID
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

        // Set up your initialization logic here
    }

    async submitTransaction() {
        // Simulate arguments for InsertTestHistoricalDataHash function
        const offerID = 'vngttx1scbc51cba2aeb1477d91707411385322'; 
        const hashID = 'hash123_'+ randomBytes(8).toString('hex');
        const dataHash = 'datahash123'+ randomBytes(8).toString('hex');
        const filename = 'filename123'+ randomBytes(8).toString('hex');
        const entrydate = '2024-02-06T12:00:00Z';
        const offerDataHashID = 'offerDataHash_'+ randomBytes(8).toString('hex'); 

        this.txIndex++;
        const args = {
            contractId: this.chaincodeID,
            contractFunction: 'InsertTestHistoricalDataHash',
            contractArguments: [
                offerID, hashID, dataHash, filename, entrydate, offerDataHashID
            ],
            readOnly: false
        };

        try {
            await this.sutAdapter.sendRequests(args);
            // Add logging or success message here
        } catch (error) {
            // Handle errors or log them accordingly
            console.error('Transaction submission error:', error);
        }
    }
}

function createWorkloadModule() {
    return new InsertDataHashWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
