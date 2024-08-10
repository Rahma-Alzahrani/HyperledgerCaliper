'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class GetAllJourneyWorkload extends WorkloadModuleBase {
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
            contractFunction: 'GetAllJourney',
            readOnly: true
        };

        const response = await this.sutAdapter.sendRequests(args);
        return response;
    }
}

function createWorkloadModule() {
    return new GetAllJourneyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
