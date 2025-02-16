import { group } from 'k6';

let _helper = (name, fn) => {
    fn();
}

export let bdd = {
    feature: (name, fn) => group(`Feature: ${name}`, fn),
    scenario: (name, fn) => group(`Scenario: ${name}`, fn),
    given: (name, fn) => _helper(`Given: ${name}`, fn),
    when: (name, fn) => _helper(`When: ${name}`, fn),
    then: (name, fn) => _helper(`Then: ${name}`, fn),
    and: (name, fn) => _helper(`And: ${name}`, fn),
    but: (name, fn) => _helper(`But: ${name}`, fn),
    rule: (name, fn) => _helper(`Rule: ${name}`, fn),
    background: (name, fn) => _helper(`Background: ${name}`, fn),
};