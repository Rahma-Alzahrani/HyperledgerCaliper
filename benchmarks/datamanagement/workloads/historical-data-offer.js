'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

const {
    randomBytes
} = require('crypto');

const uid = Math.random().toString(36).slice(2) + randomBytes(8).toString('hex') + new Date().getTime();

class InsertDataOfferWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = 0;
        this.chaincodeID = 'basic'; // Replace with your chaincode ID
        this.dataOffer = {}; // Replace with your InsertDataOffer data
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);


        this.dataOffer = {
            id: uid,
            validity: true,
            data_owner: 'UoB',
            equipment: 'sensor777',
            moniteredAsset: 'Track-london-AMST',
            processing_level: 'A',
            price: 10.0,
            deposit: 5.0,
            creator: 'UoB',
            operator: 'Avanti West Coast',
            journey_uid: '56fl71u6i4icc21c18c7ddc71811715461065176',
            start_date: '2024-05-12T12:00:00Z',
            end_date: '2024-05-12T14:00:00Z',
        };
    }

    async submitTransaction() {
        this.txIndex++;
        const args = {
            contractId: this.chaincodeID,
            contractFunction: 'InsertTestHistoricalDataOffer',
            contractArguments: [
                JSON.stringify(this.dataOffer)
            ],
            readOnly: false
        };

        try {
            await this.sutAdapter.sendRequests(args);
        } catch (error) {
            console.error('Transaction submission error:', error);
        }
    }
}

function createWorkloadModule() {
    return new InsertDataOfferWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
