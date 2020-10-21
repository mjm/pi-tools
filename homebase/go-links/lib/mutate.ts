import {client} from "com_github_mjm_pi_tools/homebase/go-links/lib/links_client";
import {CreateLinkRequest} from "com_github_mjm_pi_tools/go-links/proto/links/links_pb";

export interface CreateLinkParams {
    shortURL: string;
    destinationURL: string;
    description: string;
}

export async function createLink(params: CreateLinkParams): Promise<void> {
    const req = new CreateLinkRequest();
    req.setShortUrl(params.shortURL);
    req.setDestinationUrl(params.destinationURL);
    req.setDescription(params.description);

    return new Promise((resolve, reject) => {
        client.createLink(req, err => {
            if (err) {
                reject(err);
            } else {
                resolve();
            }
        });
    });
}
