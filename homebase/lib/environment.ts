import { Environment, Network, RecordSource, Store } from 'relay-runtime';

async function fetchRelay(params, variables) {
    const response = await fetch("/graphql", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            query: params.text,
            operationName: params.operationName,
            variables,
        }),
    });
    return await response.json();
}

const store = new Store(new RecordSource());

// @ts-ignore
window.RELAY_STORE = store;

export default new Environment({
    network: Network.create(fetchRelay),
    store,
})
