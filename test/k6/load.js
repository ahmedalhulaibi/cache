import http from 'k6/http';
import grpc from 'k6/net/grpc';
import { bdd } from './bdd.js';
import { profiles } from './profiles.js';
import { checks } from './checks.js';
import { sleep } from 'k6';
import { features } from './features.js';

export let options = profiles[__ENV.K6_TEST_PROFILE];

var httpBaseURL = `${__ENV.K6_HTTP_BASE_URL}`
var grpcBaseURL = `${__ENV.K6_GRPC_BASE_URL}`
var grpcProtoDir = `${__ENV.K6_GRPC_PROTO_DIR}`
const grpcClient = new grpc.Client(grpcBaseURL);

const protoFiles = `${__ENV.K6_GRPC_PROTO_FILES}`.split(";").filter( (file) => file.length > 0 );
grpcClient.load([grpcProtoDir, '../../googleapis', '../../grpc-gateway'], ...protoFiles);

const globalState = {
    featuresToRun: null,
};

init();

function init() {
    let envK6Features = `${__ENV.K6_FEATURES}`.split(";");
    // populate the features array from features.js
    let featuresToRun = [];
    for (const feature in features) {
        featuresToRun.push({ name: feature, scenarios: features[feature].scenarios, setup: features[feature].setup, teardown: features[feature].teardown });
    }

    // populate the features array from env
    let featureNamesToRun = envK6Features.filter(f => f != "" && f !== "undefined").map(feature => {
        let f = feature.split(":");
        return { name: f[0], scenarios: f.length > 1 ? f[1].split(",") : [] };
    });

    // filter out features that are not in the env
    if (featureNamesToRun.length > 0) {
        featuresToRun = featuresToRun.filter(feature => featureNamesToRun.find(f => f.name === feature.name));
    }

    // for each featureNamesToRun, filter out scenarios that are not in the env
    // if the scenarios are empty for a given feature, then run all scenarios
    featuresToRun = featuresToRun.map(feature => {
        const featenv = featureNamesToRun.find(f => f.name === feature.name);

        for (const scenario in feature.scenarios) {
            if (featenv) {
                // if the scenario is not in the env, remove it
                if (featenv.scenarios.length > 0 && !featenv.scenarios.includes(scenario)) {
                    delete feature.scenarios[scenario];
                }
            }
        }
        return feature;
    });

    // for each featureToRun, run setup
    for (const feature of featuresToRun) {
        if (feature.setup) {
            feature.setup(globalState);
        }
    }

    globalState.featuresToRun = featuresToRun;
};

export default function () {
    grpcClient.connect(grpcBaseURL, { plaintext: true });

    globalState.featuresToRun.forEach(feature => {
        bdd.feature(feature.name, () => {
            for (const scenario in feature.scenarios) {
                bdd.scenario(scenario, () => features[feature.name].scenarios[scenario]((cbData) => { globalState.createdMemeIDs.push(cbData.id); }));
            }
        });
    });
    // console.log('Created memes after test', JSON.stringify(globalState.createdMemeIDs));
    globalState.featuresToRun.forEach(feature => { if (feature.teardown) { feature.teardown(globalState); } });

    grpcClient.close();
};

export function teardown() {
};