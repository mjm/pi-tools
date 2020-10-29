import {TripsServiceClient} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb_service";
import {promisifyClient} from "com_github_mjm_pi_tools/homebase/lib/promisify";

const serviceHost = process.env.NODE_ENV === "development"
    ? `http://${window.location.hostname}:2120`
    : `${window.location.protocol}//${window.location.host}`;

export const client = promisifyClient(new TripsServiceClient(serviceHost));
