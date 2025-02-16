import http from 'k6/http';
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
    }
};