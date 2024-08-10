'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class GetAllOfferWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.chaincodeID = 'basic';
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
    }

    async submitTransaction() {
        const args = {
            contractId: 'basic',
            contractFunction: 'GetAllHistoricalOffer',
            readOnly: true
        };

        await this.sutAdapter.sendRequests(args);
    }
}

function createWorkloadModule() {
    return new GetAllOfferWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
