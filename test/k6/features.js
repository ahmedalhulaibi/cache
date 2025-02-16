import http, { options } from 'k6/http';
import grpc from 'k6/net/grpc';
import exec from 'k6/execution';
import { bdd } from './bdd.js';
import { checks } from './checks.js';
import { sleep } from 'k6';

export let features = {
    "Hello World": {
        scenarios: {
            "Basic scenario": () => {
                let helloResponse;
                bdd.given("A name ahmed", function () {
                    var name = "ahmed";
                });
                bdd.when("Ahmed", function () {
                    helloResponse = http.get(`${httpBaseURL}/v1/hello?name=ahmed`);
                });
                bdd.then("Expected outcome in english", function () {
                    checks.is200(helloResponse);
                    checks.isJSON(helloResponse);
                    checks.assert(helloResponse,"name is in greeting",(r) => r.json().message, "Hello, ahmed! Ya filthy animal.");
                });
                sleep(3);
            },
        },
        setup: (globalState) => {
        },
        teardown: (globalState) => {
        },
    },
    "Cache": {
        scenarios: {
            "Can set and get": () => {
                let setResponse;
                let getResponse;
                let bucket = `default`;
                let key = `mykey-${exec.vu.idInTest}-${exec.scenario.iterationInTest}`;
                let value = `myvalue"-${exec.vu.idInTest}`;
                bdd.when("Set key and value", function () {
                    setResponse = http.post(`${httpBaseURL}/v1/set`, JSON.stringify({
                        bucket: bucket,
                        key: key,
                        value: value,
                        options: {
                            ttlSeconds: 1,
                        }
                    }));
                });
                bdd.then("Get key and value", function () {
                    getResponse = http.get(`${httpBaseURL}/v1/get/${bucket}/${key}`);
                    checks.is200(getResponse);
                    checks.isJSON(getResponse);
                    checks.assert(getResponse,"value is in response",(r) => r.json().value === value);
                });
                bdd.then("Wait", function () {
                    sleep(2);
                });
                bdd.then("Get key and value again", function () {
                    getResponse = http.get(`${httpBaseURL}/v1/get/${bucket}/${key}`);
                    checks.is200(getResponse);
                    checks.assert(getResponse,"value is blank in response",(r) => r.json().value === "");
                });
            }
        },
    },
};