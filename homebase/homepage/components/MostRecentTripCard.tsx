import React from "react";
import {Link} from "react-router-dom";
import useSWR from "swr";
import {Trip} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {fetcher, GET_MOST_RECENT_TRIP} from "com_github_mjm_pi_tools/homebase/trips/lib/fetch";
import {Alert} from "com_github_mjm_pi_tools/homebase/components/Alert";
import {formatDuration, intervalToDuration, parseISO} from "date-fns";

export function MostRecentTripCard() {
    const {data, error} = useSWR<Trip | null>(GET_MOST_RECENT_TRIP, fetcher);
    if (error) {
        console.error(error);

        return (
            <Alert title="Couldn't load most recent trip" severity="error">
                {error.toString()}
            </Alert>
        );
    }

    return (
        <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-5">
                <div className="flex items-center">
                    <div className="flex-shrink-0">
                        <svg className="h-6 w-6 text-cool-gray-400" xmlns="http://www.w3.org/2000/svg"
                             fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                  d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7"/>
                        </svg>
                    </div>
                    <div className="ml-5 w-0 flex-1">
                        {data !== undefined ? (
                            <dl>
                                <dt className="text-sm leading-5 font-medium text-cool-gray-500 truncate">
                                    {data && data.getReturnedAt() === "" ? <>Current trip</> : <>Most recent trip</>}
                                </dt>
                                <dd>
                                    <div className="text-lg leading-7 font-medium text-cool-gray-900">
                                        {data ? (
                                            formatDuration(intervalToDuration({
                                                start: parseISO(data.getLeftAt()),
                                                end: data.getReturnedAt() ? parseISO(data.getReturnedAt()) : new Date()
                                            }))
                                        ) : (
                                            "None"
                                        )}
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
                    <Link to="/trips"
                          className="font-medium text-indigo-600 hover:text-indigo-900 transition ease-in-out duration-150">
                        View recent trips
                    </Link>
                </div>
            </div>
        </div>
    );
}
