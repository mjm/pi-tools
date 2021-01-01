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

    return (
        <HomePageCard
            title={deploy.state === "IN_PROGRESS" ? "Currently deploying" : "Most recent deploy"}
            icon={null}
            footerHref="https://github.com/mjm/pi-tools/deployments"
            footer="View deploy history">
            <a href={`https://github.com/mjm/pi-tools/commit/${deploy.commitSHA}`} target="_blank" className="text-base">
                {deploy.commitMessage}
            </a>
        </HomePageCard>
    );
}