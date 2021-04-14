import React from "react";
import {graphql, useFragment} from "react-relay/hooks";
import {MostRecentDeployCard_viewer$key} from "../../__generated__/MostRecentDeployCard_viewer.graphql";
import HomePageCard from "./HomePageCard";
import Link from "next/link";
import {CodeIcon} from "@heroicons/react/outline";

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
            icon={<CodeIcon className="h-6 w-6 text-gray-400"/>}
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
