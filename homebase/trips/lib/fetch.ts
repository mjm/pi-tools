import {
    GetTripRequest,
    GetTripResponse,
    ListTripsRequest,
    ListTripsResponse,
    Trip,
} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {client} from "com_github_mjm_pi_tools/homebase/trips/lib/trips_client";
import {promisify} from "com_github_mjm_pi_tools/homebase/lib/promisify";

export const LIST_TRIPS = "ListTrips";
export const GET_TRIP = "GetTrip";

const fetchers = {
    [LIST_TRIPS]: async (): Promise<Trip[]> => {
        const req = new ListTripsRequest();
        const res: ListTripsResponse = await promisify(client, "listTrips")(req);
        return res.getTripsList();
    },
    [GET_TRIP]: async (id: string): Promise<Trip> => {
        const req = new GetTripRequest();
        req.setId(id);
        const res: GetTripResponse = await promisify(client, "getTrip")(req);
        return res.getTrip();
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
