import React from "react";
import HomePageCard from "./HomePageCard";
import {graphql, useFragment} from "react-relay/hooks";
import {FiringAlertsCard_viewer$key} from "../../__generated__/FiringAlertsCard_viewer.graphql";

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
            icon={
                <svg className="h-6 w-6 text-gray-400" xmlns="http://www.w3.org/2000/svg"
                     fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                          d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
                </svg>
            }
            footerHref="https://alertmanager.homelab/"
            footer="View active alerts"
        >
            {data.alerts.length}
        </HomePageCard>
    );
}