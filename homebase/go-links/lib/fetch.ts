import {client} from "com_github_mjm_pi_tools/homebase/go-links/lib/links_client";
import {Link, ListRecentLinksRequest} from "com_github_mjm_pi_tools/go-links/proto/links/links_pb";

export const LIST_RECENT_LINKS = "ListRecentLinks";

const fetchers = {
    [LIST_RECENT_LINKS]: async (): Promise<Link[]> => {
        const req = new ListRecentLinksRequest();
        return new Promise((resolve, reject) => {
            client.listRecentLinks(req, (err, res) => {
                if (err) {
                    reject(err);
                } else {
                    resolve(res.getLinksList());
                }
            });
        });
    },
} as const;

type Fetchers = typeof fetchers;
type FetchKey = keyof Fetchers;

export async function fetcher<K extends FetchKey>(method: K, ...args: Parameters<Fetchers[K]>): Promise<any> {
    const fetchFn = fetchers[method] as any;
    if (!fetchFn) {
        throw new Error(`Unimplemented method to fetch: ${method}`);
    }
    return await fetchFn(...args);
}
