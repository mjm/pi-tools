import {Environment, Network, RecordSource, Store} from "relay-runtime";
import {getRelaySerializedState} from "relay-nextjs";

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

let clientEnv: Environment | undefined;

export function getClientEnvironment() {
    if (typeof window === "undefined") {
        return null;
    }

    if (!clientEnv) {
        clientEnv = new Environment({
            network: Network.create(fetchRelay),
            store: new Store(new RecordSource(getRelaySerializedState()?.records)),
            isServer: false,
        });
    }

    return clientEnv;
}
