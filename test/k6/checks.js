import { check } from 'k6';

let _checks = {
    is200: function (res) {
        check(res, {
            'is status 200': (r) => r.status === 200,
        });
    },
    isArray: function (res, fieldName) {
        check(res, {
            'is array': (r) => Array.isArray(r.json()[fieldName]),
        });
    },
    isJSON: function (res) {
        check(res, {
            'is json': (r) => r.json() !== undefined && r.json() !== null,
        });
    },
    assert: function (res, name, fn) {
        check(res, {
            [name]: (r) => fn(r),
        });
    }
};

export let checks = _checks;