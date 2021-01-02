import React from "react";
import {graphql, useLazyLoadQuery} from "react-relay/hooks";
import {PageHeader} from "com_github_mjm_pi_tools/homebase/components/PageHeader";
import {HomePageQuery} from "com_github_mjm_pi_tools/homebase/api/__generated__/HomePageQuery.graphql";
import {MostRecentTripCard} from "com_github_mjm_pi_tools/homebase/homepage/components/MostRecentTripCard";
import {FiringAlertsCard} from "com_github_mjm_pi_tools/homebase/homepage/components/FiringAlertsCard";
import {MostRecentDeployCard} from "com_github_mjm_pi_tools/homebase/homepage/components/MostRecentDeployCard";

export function HomePage() {
    const data = useLazyLoadQuery<HomePageQuery>(
        graphql`
            query HomePageQuery {
                viewer {
                    ...MostRecentTripCard_viewer
                    ...FiringAlertsCard_viewer
                    ...MostRecentDeployCard_viewer
                }
            }
        `,
        {},
    );

    return (
        <main className="mb-8">
            <PageHeader>Homebase</PageHeader>
            <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                <div className="mt-2 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
                    <MostRecentTripCard viewer={data.viewer}/>
                    <FiringAlertsCard viewer={data.viewer}/>
                    <MostRecentDeployCard viewer={data.viewer}/>
                </div>
            </div>
        </main>
    );
}
