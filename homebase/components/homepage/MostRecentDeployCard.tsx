import React from "react";
import {graphql, useFragment} from "react-relay/hooks";
import {MostRecentDeployCard_viewer$key} from "../../__generated__/MostRecentDeployCard_viewer.graphql";
import HomePageCard from "./HomePageCard";
import Link from "next/link";

export default function MostRecentDeployCard({viewer}: { viewer: MostRecentDeployCard_viewer$key }) {
    const data = useFragment(
        graphql`
            fragment MostRecentDeployCard_viewer on Viewer {
                mostRecentDeploy {
                    rawID
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
            footerHref="/deploys"
            footer="View deploy history">
            <Link href={`/deploys/${deploy.rawID}`}>
                <a className="text-base">
                    {deploy.commitMessage}
                </a>
            </Link>
        </HomePageCard>
    );
}