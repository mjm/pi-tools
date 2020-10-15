import {GetTripRequest, ListTripsRequest, Trip} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {client} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/lib/trips_client";

export const LIST_TRIPS = "ListTrips";
export const GET_TRIP = "GetTrip";

const fetchers = {
    [LIST_TRIPS]: async (): Promise<Trip[]> => {
        const req = new ListTripsRequest();
        return new Promise((resolve, reject) => {
            client.listTrips(req, (err, res) => {
                if (err) {
                    reject(err);
                } else {
                    resolve(res.getTripsList());
                }
            });
        });
    },
    [GET_TRIP]: async (id: string): Promise<Trip> => {
        const req = new GetTripRequest();
        req.setId(id);
        return new Promise((resolve, reject) => {
            client.getTrip(req, (err, res) => {
                if (err) {
                    reject(err);
                } else {
                    resolve(res.getTrip());
                }
            })
        })
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
