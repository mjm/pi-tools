import React from "react";
import {PageHeader} from "com_github_mjm_pi_tools/homebase/components/PageHeader";
import {NewLinkCard} from "com_github_mjm_pi_tools/homebase/go-links/components/NewLinkCard";
import {RecentLinksList} from "com_github_mjm_pi_tools/homebase/go-links/components/RecentLinksList";
import {graphql, useLazyLoadQuery} from "react-relay/hooks";
import {GoLinksHomePageQuery} from "com_github_mjm_pi_tools/homebase/api/__generated__/GoLinksHomePageQuery.graphql";

export function GoLinksHomePage() {
    const data = useLazyLoadQuery<GoLinksHomePageQuery>(
        graphql`
            query GoLinksHomePageQuery {
                viewer {
                    ...RecentLinksList_viewer
                }
            }
        `,
        {}
    );

    return (
        <main className="mb-8">
            <PageHeader>
                Go links
            </PageHeader>
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="grid grid-cols-1 sm:grid-cols-2 mt-6 gap-8">
                    <div><NewLinkCard/></div>
                    <div><RecentLinksList viewer={data.viewer}/></div>
                </div>
            </div>
        </main>
    );
}
