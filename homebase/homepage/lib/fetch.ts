import {client, PrometheusAlert} from "com_github_mjm_pi_tools/homebase/homepage/lib/prometheus";

export const LIST_FIRING_ALERTS = "ListFiringAlerts";

const fetchers = {
    [LIST_FIRING_ALERTS]: async (): Promise<PrometheusAlert[]> => {
        return await client.alerts();
    },
} as const;

type Fetchers = typeof fetchers;
type FetchKey = keyof Fetchers;

export async function fetcher<K extends FetchKey>(method: K, ...args: Parameters<Fetchers[K]>): Promise<any> {
    const fetchFn = fetchers[method] as any;
    if (!fetchFn) {
        throw new Error(`Unimplemented method to fetch: ${method}`);
    }
    return await fetchFn(...args);
}
