import {client} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/lib/trips_client";
import {IgnoreTripRequest} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {mutate} from "swr";
import {LIST_TRIPS} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/lib/fetch";

export async function ignoreTrip(id: string): Promise<void> {
    const req = new IgnoreTripRequest();
    req.setId(id);
    return new Promise((resolve, reject) => {
        client.ignoreTrip(req, err => {
            if (err) {
                reject(err);
            } else {
                mutate(LIST_TRIPS);
                resolve();
            }
        });
    });
}
