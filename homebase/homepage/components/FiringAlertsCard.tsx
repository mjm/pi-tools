import React from "react";
import useSWR from "swr";
import {PrometheusAlert} from "com_github_mjm_pi_tools/homebase/homepage/lib/prometheus";
import {fetcher, LIST_FIRING_ALERTS} from "com_github_mjm_pi_tools/homebase/homepage/lib/fetch";
import {Alert} from "com_github_mjm_pi_tools/homebase/components/Alert";

export function FiringAlertsCard() {
    const {data, error} = useSWR<PrometheusAlert[]>(LIST_FIRING_ALERTS, fetcher);
    if (error) {
        console.error(error);
    }

    return (
        <div className="bg-white overflow-hidden shadow rounded-lg">
            {error ? (
                <Alert title="Couldn't load firing alerts" severity="error" rounded={false}>
                    {error.toString()}
                </Alert>
            ) : (
                <>
                    <div className="p-5">
                        <div className="flex items-center">
                            <div className="flex-shrink-0">
                                <svg className="h-6 w-6 text-cool-gray-400" xmlns="http://www.w3.org/2000/svg"
                                     fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                          d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
                                </svg>
                            </div>
                            <div className="ml-5 w-0 flex-1">
                                {data !== undefined ? (
                                    <dl>
                                        <dt className="text-sm leading-5 font-medium text-cool-gray-500 truncate">
                                            Alerts firing
                                        </dt>
                                        <dd>
                                            <div className="text-lg leading-7 font-medium text-cool-gray-900">
                                                {data.length}
                                            </div>
                                        </dd>
                                    </dl>
                                ) : (
                                    <dl>
                                        <dt className="text-sm leading-5 font-medium text-cool-gray-500 truncate">
                                            Loading…
                                        </dt>
                                    </dl>
                                )}
                            </div>
                        </div>
                    </div>
                    <div className="bg-cool-gray-50 px-5 py-3">
                        <div className="text-sm leading-5">
                            <a href="https://alertmanager.homelab/"
                               className="font-medium text-indigo-600 hover:text-indigo-900 transition ease-in-out duration-150">
                                View active alerts
                            </a>
                        </div>
                    </div>
                </>
            )}
        </div>
    );
}
