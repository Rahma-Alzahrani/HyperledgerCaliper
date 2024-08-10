'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

const {
    randomBytes
  } = require('crypto');
  
  const uid = Math.random().toString(36).slice(2) + randomBytes(8).toString('hex') + new Date().getTime();
/**
 * Workload module for the benchmark round.
 */
class JourneyScheduleWorkload extends WorkloadModuleBase {
    /**
     * Initializes the workload module instance.
     */
    constructor() {
        super();
        this.txIndex = 0;
        this.chaincodeID = 'basic'; // Replace with your chaincode ID
        this.journeySchedule = {}; // Replace with your JourneySchedule data
    }


    /**
     * Initialize the workload module with the given parameters.
     * @param {number} workerIndex The 0-based index of the worker instantiating the workload module.
     * @param {number} totalWorkers The total number of workers participating in the round.
     * @param {number} roundIndex The 0-based index of the currently executing round.
     * @param {Object} roundArguments The user-provided arguments for the round from the benchmark configuration file.
     * @param {BlockchainInterface} sutAdapter The adapter of the underlying SUT.
     * @param {Object} sutContext The custom context object provided by the SUT adapter.
     * @async
     */
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

        const args = this.roundArguments;
        this.chaincodeID = args.chaincodeID ? args.chaincodeID : 'basic'; // Replace with your chaincode ID


        // Replace the following data with your JourneySchedule input data
        this.journeySchedule = {
            DocType: 'Journey', // Replace with your data
            UID: uid,
            Type: 'WTT',
            Cat: '0S00',
            Journey: '1715 Birmingham New Street - London Euston',
            Valid_from: '2024-05-11',
            Valid_to: '2024-05-24',
            Days: 'SX',
            Operator: 'Chiltern Railways'
        };
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        // const uuid = 'client' + this.workerIndex + '_' + this.txIndex;
        // this.journeySchedule.UID = uuid;
        this.txIndex++;
        const args = {
            contractId: this.chaincodeID,
            contractFunction: 'JourneySchedule', // Replace with your contract function name
            contractArguments: [
                // uuid,
                JSON.stringify(this.journeySchedule)
            ],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(args);
    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new JourneyScheduleWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
