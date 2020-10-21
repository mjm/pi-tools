import {LinksServiceClient} from "com_github_mjm_pi_tools/go-links/proto/links/links_pb_service";

const serviceHost = process.env.NODE_ENV === "development"
    ? "http://localhost:4240"
    : `${window.location.protocol}//${window.location.host}`;

export const client = new LinksServiceClient(serviceHost);
