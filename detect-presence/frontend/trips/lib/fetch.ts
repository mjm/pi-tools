import {TripsServiceClient} from "../proto/trips_pb_service";
import {ListTripsRequest} from "../proto/trips_pb";

const client = new TripsServiceClient('http://localhost:2120')

export async function fetcher(method: string, ...args: any[]): Promise<any> {
    switch (method) {
        case 'ListTrips':
            const req = new ListTripsRequest()
            return new Promise((resolve, reject) => {
                client.listTrips(req, (err, res) => {
                    if (err) {
                        reject(err)
                    } else {
                        resolve(res.getTripsList())
                    }
                })
            })
        default:
            throw new Error(`Unimplemented method to fetch: ${method}`)
    }
}
