'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class GetAllOfferWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        // Any initialization logic can be added here
    }

    async submitTransaction() {
        const args = {
            contractId: 'basic',
            contractFunction: 'GetAllDataHashes',
            readOnly: true
        };

        await this.sutAdapter.sendRequests(args);
    }
}

function createWorkloadModule() {
    return new GetAllOfferWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
