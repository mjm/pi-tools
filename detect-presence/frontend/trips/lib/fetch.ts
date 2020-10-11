import {ListTripsRequest} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {client} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/lib/trips_client";

export const LIST_TRIPS = "ListTrips";

export async function fetcher(method: string, ...args: any[]): Promise<any> {
    switch (method) {
        case LIST_TRIPS:
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
        default:
            throw new Error(`Unimplemented method to fetch: ${method}`);
    }
}
