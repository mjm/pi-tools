import {TripsServiceClient} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb_service";
import {ListTripsRequest} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";

const serviceHost = process.env.NODE_ENV === "development" ? "http://localhost:2120" : window.location.href;

const client = new TripsServiceClient(serviceHost);

export async function fetcher(method: string, ...args: any[]): Promise<any> {
    switch (method) {
        case "ListTrips":
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
