'use strict';

const {WorkloadModuleBase} = require('@hyperledger/caliper-core');

class MyWorkLoad extends WorkloadModuleBase {
    constructor() {
        super();
    }


    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
    }

    async cleanupWorkloadModule() {
        await super.cleanupWorkloadModule();
    }

    async submitTransaction() {

        const myArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'ReadUserData',
            invokerIdentity: 'bgihr',
            contractArguments: [`eDUwOTo6Q049YmdpaHIsT1U9Y2xpZW50LE89SHlwZXJsZWRnZXIsU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUzo6Q049Y2Eub3JnMS5leGFtcGxlLmNvbSxPPW9yZzEuZXhhbXBsZS5jb20sTD1EdXJoYW0sU1Q9Tm9ydGggQ2Fyb2xpbmEsQz1VUw==`],
            readOnly: true
        };

        await this.sutAdapter.sendRequests(myArgs);
    }
}

function createWorkloadModule() {
    return new MyWorkLoad();
}

module.exports.createWorkloadModule = createWorkloadModule;