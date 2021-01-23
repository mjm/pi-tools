import React from "react";
import {graphql, useFragment} from "react-relay/hooks";
import {MostRecentDeployCard_viewer$key} from "com_github_mjm_pi_tools/homebase/api/__generated__/MostRecentDeployCard_viewer.graphql";
import {HomePageCard} from "com_github_mjm_pi_tools/homebase/homepage/components/HomePageCard";

export function MostRecentDeployCard({viewer}: { viewer: MostRecentDeployCard_viewer$key }) {
    const data = useFragment(
        graphql`
            fragment MostRecentDeployCard_viewer on Viewer {
                mostRecentDeploy {
                    commitSHA
                    commitMessage
                    state
                    startedAt
                    finishedAt
                }
            }
        `,
        viewer,
    );

    const deploy = data.mostRecentDeploy;
    if (!deploy) {
        return null;
    }

    return (
        <HomePageCard
            title={deploy.state === "IN_PROGRESS" ? "Currently deploying" : "Most recent deploy"}
            icon={
                <svg className="h-6 w-6 text-gray-400" xmlns="http://www.w3.org/2000/svg" fill="none"
                     viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                          d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"/>
                </svg>
            }
            footerHref="https://github.com/mjm/pi-tools/deployments"
            footer="View deploy history">
            <a href={`https://github.com/mjm/pi-tools/commit/${deploy.commitSHA}`} target="_blank"
               className="text-base">
                {deploy.commitMessage}
            </a>
        </HomePageCard>
    );
}
