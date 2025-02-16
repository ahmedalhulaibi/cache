export let profiles = {
    load: {
        stages: [
            { duration: '5m', target: 100 }, // simulate ramp-up of traffic from 1 to 100 users over 5 minutes.
            { duration: '10m', target: 100 }, // stay at 100 users for 10 minutes
            { duration: '5m', target: 0 }, // ramp-down to 0 users
        ],
        thresholds: {
            http_req_duration: ['p(99)<1000'], // 99% of requests must complete below 1s
        },
    },
    smoke: {
        thresholds: {
            http_req_failed: ['rate==0'],   // http errors should be less than 1% 
            http_req_duration: ['p(95)<1000'], // 95% of requests should be below 1000ms
        },
        scenarios: {
            contacts: {
                executor: 'per-vu-iterations',
                vus: 1,
                iterations: 1,
                maxDuration: '1m',
            },
        },
    },
    stress: {
        stages: [
            { duration: '2m', target: 100 }, // below normal load
            { duration: '5m', target: 100 },
            { duration: '2m', target: 200 }, // normal load
            { duration: '5m', target: 200 },
            { duration: '2m', target: 300 }, // around the breaking point
            { duration: '5m', target: 300 },
            { duration: '2m', target: 400 }, // beyond the breaking point
            { duration: '5m', target: 400 },
            { duration: '10m', target: 0 }, // scale down. Recovery stage.
        ],
        thresholds: {
            http_req_failed: ['rate<0.01'],   // http errors should be less than 1% 
            http_req_duration: ['p(95)<1000'], // 95% of requests should be below 1000ms
        },
    }
};