import {TripsServiceClient} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb_service";

const serviceHost = process.env.NODE_ENV === "development"
    ? "http://localhost:2120"
    : `${window.location.protocol}//${window.location.host}`;

export const client = new TripsServiceClient(serviceHost);
