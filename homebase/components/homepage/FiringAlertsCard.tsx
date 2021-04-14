import React from "react";
import HomePageCard from "./HomePageCard";
import {graphql, useFragment} from "react-relay/hooks";
import {FiringAlertsCard_viewer$key} from "../../__generated__/FiringAlertsCard_viewer.graphql";
import {ExclamationIcon} from "@heroicons/react/outline";

export default function FiringAlertsCard({viewer}: { viewer: FiringAlertsCard_viewer$key }) {
    const data = useFragment(
        graphql`
            fragment FiringAlertsCard_viewer on Viewer {
                alerts {
                    activeAt
                    value
                }
            }
        `,
        viewer,
    );

    return (
        <HomePageCard
            title="Alerts firing"
            icon={<ExclamationIcon className="h-6 w-6 text-gray-400"/>}
            footerHref="https://alertmanager.homelab/"
            footer="View active alerts"
        >
            {data.alerts.length}
        </HomePageCard>
    );
}
