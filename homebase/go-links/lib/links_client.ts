import {LinksServiceClient} from "com_github_mjm_pi_tools/go-links/proto/links/links_pb_service";
import {promisifyClient} from "com_github_mjm_pi_tools/homebase/lib/promisify";

const serviceHost = process.env.NODE_ENV === "development"
    ? `http://${window.location.hostname}:4240`
    : `${window.location.protocol}//${window.location.host}`;

export const client = promisifyClient(new LinksServiceClient(serviceHost));

export function destinationURL(shortURL: string): string {
    if (process.env.NODE_ENV === "development") {
        return `${serviceHost}/${shortURL}`;
    } else {
        return `http://go/${shortURL}`;
    }
}
