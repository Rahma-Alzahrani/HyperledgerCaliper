'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const { randomBytes } = require('crypto');


const uid = Math.random().toString(36).slice(2) + randomBytes(8).toString('hex') + new Date().getTime();

class InsertDataOfferWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.chaincodeID = 'basic'; // Replace with your chaincode ID
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);



        this.dataOffer = {
            id: uid,
            validity: true,
            data_owner: 'UoB',
            equipment: 'sensor123',
            monitered_asset: 'Axle Journal Bearing',
            processing_level: 'A',
            price: 10.0,
            deposit: 5.0,
            creator: 'UoB',
            operator: 'London Northwestern Railway',
            journey_uid: 'y7mj8wici647352b9fea001c6c1715460922540',
            depart_time: '2024-05-12T12:00:00Z',
            arrival_time: '2024-05-12T14:00:00Z',
            is_active: true
        };
    }

    async submitTransaction() {
        this.txIndex++;
        const args = {
            contractId: this.chaincodeID,
            contractFunction: 'InsertDataOffer',
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
