import React from "react";
import useSWR from "swr";
import {PrometheusAlert} from "com_github_mjm_pi_tools/homebase/homepage/lib/prometheus";
import {fetcher, LIST_FIRING_ALERTS} from "com_github_mjm_pi_tools/homebase/homepage/lib/fetch";
import {Alert} from "com_github_mjm_pi_tools/homebase/components/Alert";
import {HomePageCard} from "com_github_mjm_pi_tools/homebase/homepage/components/HomePageCard";

export function FiringAlertsCard() {
    const {data, error} = useSWR<PrometheusAlert[]>(LIST_FIRING_ALERTS, fetcher);
    if (error) {
        console.error(error);
        return (
            <Alert title="Couldn't load firing alerts" severity="error">
                {error.toString()}
            </Alert>
        );
    }

    return (
        <HomePageCard
            title={data ? "Alerts firing" : "Loading…"}
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
            {data ? data.length : null}
        </HomePageCard>
    );
}
