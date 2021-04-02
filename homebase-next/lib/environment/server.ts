import {Environment, Network, RecordSource, Store} from "relay-runtime";
import fetch from "node-fetch";
// import * as nodeFetch from "node-fetch";
// import createFetch from "@vercel/fetch";
//
// const fetch = createFetch(nodeFetch);

const serverUrl = process.env.GRAPHQL_URL || "http://localhost:8080/graphql";

async function fetchRelay(params, variables) {
    const response = await fetch(serverUrl, {
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

export function createServerEnvironment() {
    return new Environment({
        network: Network.create(fetchRelay),
        store: new Store(new RecordSource()),
        isServer: true,
    });
}
