import React from "react";
import {graphql, PreloadedQuery, usePreloadedQuery} from "react-relay/hooks";
import {PageHeader} from "com_github_mjm_pi_tools/homebase/components/PageHeader";
import {HomePageQuery} from "com_github_mjm_pi_tools/homebase/api/__generated__/HomePageQuery.graphql";
import {MostRecentTripCard} from "com_github_mjm_pi_tools/homebase/homepage/components/MostRecentTripCard";
import {FiringAlertsCard} from "com_github_mjm_pi_tools/homebase/homepage/components/FiringAlertsCard";
import {MostRecentDeployCard} from "com_github_mjm_pi_tools/homebase/homepage/components/MostRecentDeployCard";
import {ErrorBoundary} from "com_github_mjm_pi_tools/homebase/components/ErrorBoundary";

export function HomePage({prepared}) {
    return (
        <main className="mb-8">
            <PageHeader>Homebase</PageHeader>
            <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                <React.Suspense fallback="Loading…">
                    <ErrorBoundary>
                        <HomePageInner homeQuery={prepared.homeQuery}/>
                    </ErrorBoundary>
                </React.Suspense>
            </div>
        </main>
    );
}

function HomePageInner({homeQuery}: { homeQuery: PreloadedQuery<HomePageQuery> }) {
    const data = usePreloadedQuery<HomePageQuery>(
        graphql`
            query HomePageQuery {
                viewer {
                    ...MostRecentTripCard_viewer
                    ...FiringAlertsCard_viewer
                    ...MostRecentDeployCard_viewer
                }
            }
        `,
        homeQuery,
    );

    return (
        <div className="mt-2 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
            <MostRecentTripCard viewer={data.viewer}/>
            <FiringAlertsCard viewer={data.viewer}/>
            <MostRecentDeployCard viewer={data.viewer}/>
        </div>
    );
}
