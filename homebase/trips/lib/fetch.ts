import {ListTripsRequest, Trip} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {client} from "com_github_mjm_pi_tools/homebase/trips/lib/trips_client";

export const GET_MOST_RECENT_TRIP = "GetMostRecentTrip";

const fetchers = {
    [GET_MOST_RECENT_TRIP]: async (): Promise<Trip | null> => {
        const req = new ListTripsRequest();
        req.setLimit(1);
        const res = await client.listTrips(req);
        const trips = res.getTripsList();
        if (trips.length === 0) {
            return null;
        }
        return trips[0];
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
