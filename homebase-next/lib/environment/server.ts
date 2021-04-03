import {Environment, Network, RecordSource, Store} from "relay-runtime";
import fetch from "node-fetch";
// import * as nodeFetch from "node-fetch";
// import createFetch from "@vercel/fetch";
//
// const fetch = createFetch(nodeFetch);

const serverUrl = process.env.GRAPHQL_URL || "http://localhost:3000/graphql";

export function createServerEnvironment(cookie: string) {
    async function fetchRelay(params, variables) {
        const response = await fetch(serverUrl, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Cookie": cookie,
            },
            body: JSON.stringify({
                query: params.text,
                operationName: params.operationName,
                variables,
            }),
        });
        return await response.json();
    }

    return new Environment({
        network: Network.create(fetchRelay),
        store: new Store(new RecordSource()),
        isServer: true,
    });
}
