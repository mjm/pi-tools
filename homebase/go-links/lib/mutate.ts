import {client} from "com_github_mjm_pi_tools/homebase/go-links/lib/links_client";
import {CreateLinkRequest, Link} from "com_github_mjm_pi_tools/go-links/proto/links/links_pb";
import {mutate} from "swr";
import {LIST_RECENT_LINKS} from "com_github_mjm_pi_tools/homebase/go-links/lib/fetch";

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
        client.createLink(req, (err, res) => {
            if (err) {
                reject(err);
            } else {
                if (res) {
                    mutate(LIST_RECENT_LINKS, (links: Link[]) => {
                        return [res.getLink(), ...links];
                    });
                } else {
                    mutate(LIST_RECENT_LINKS);
                }
                resolve();
            }
        });
    });
}
