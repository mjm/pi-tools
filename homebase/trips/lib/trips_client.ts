import {TripsServiceClient} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb_service";
import {promisifyClient} from "com_github_mjm_pi_tools/homebase/lib/promisify";

const serviceHost = `${window.location.protocol}//${window.location.host}`;

export const client = promisifyClient(new TripsServiceClient(serviceHost));
