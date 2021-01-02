import {client} from "com_github_mjm_pi_tools/homebase/go-links/lib/links_client";
import {CreateLinkRequest, UpdateLinkRequest} from "com_github_mjm_pi_tools/go-links/proto/links/links_pb";

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

    await client.createLink(req);
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

    await client.updateLink(req);
}
