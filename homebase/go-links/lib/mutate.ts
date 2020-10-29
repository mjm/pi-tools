import {client} from "com_github_mjm_pi_tools/homebase/go-links/lib/links_client";
import {CreateLinkRequest, Link, UpdateLinkRequest} from "com_github_mjm_pi_tools/go-links/proto/links/links_pb";
import {mutate} from "swr";
import {GET_LINK, LIST_RECENT_LINKS} from "com_github_mjm_pi_tools/homebase/go-links/lib/fetch";

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

    const res = await client.createLink(req);

    await Promise.all([
        mutate(LIST_RECENT_LINKS, (links: Link[]) => {
            return [res.getLink(), ...links];
        }),
        mutate([GET_LINK, res.getLink().getId()], res.getLink()),
    ]);
}

export interface UpdateLinkParams {
    id: string;
    shortURL: string;
    destinationURL: string;
    description: string;
}

export async function updateLink(params: UpdateLinkParams): Promise<void> {
    const req = new UpdateLinkRequest();
    req.setId(params.id);
    req.setShortUrl(params.shortURL);
    req.setDestinationUrl(params.destinationURL);
    req.setDescription(params.description);

    const res = await client.updateLink(req);
    await Promise.all([
        mutate(LIST_RECENT_LINKS),
        mutate([GET_LINK, res.getLink().getId()], res.getLink()),
    ]);
}
