import React from "react";
import {GoLinksHomePage} from "com_github_mjm_pi_tools/homebase/go-links/components/GoLinksHomePage";
import {GoLinkDetailPage} from "com_github_mjm_pi_tools/homebase/go-links/components/GoLinkDetailPage";
import {loadQuery} from "react-relay/hooks";
import RelayEnvironment from "com_github_mjm_pi_tools/homebase/lib/environment";
import GoLinksHomePageQuery from "com_github_mjm_pi_tools/homebase/api/__generated__/GoLinksHomePageQuery.graphql";
import GoLinkDetailPageQuery from "com_github_mjm_pi_tools/homebase/api/__generated__/GoLinkDetailPageQuery.graphql";

export function goLinksRoutes(path: string): any {
    return [
        {
            path,
            exact: true,
            component: GoLinksHomePage,
            prepare() {
                return {
                    linksQuery: loadQuery(RelayEnvironment, GoLinksHomePageQuery, {}),
                };
            },
        },
        {
            path: `${path}/:id`,
            component: GoLinkDetailPage,
            prepare({id}) {
                return {
                    linkQuery: loadQuery(RelayEnvironment, GoLinkDetailPageQuery, {id}),
                };
            },
        },
    ];
}
