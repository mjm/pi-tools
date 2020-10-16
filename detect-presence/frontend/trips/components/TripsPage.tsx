import React from "react";
import useSWR from "swr";
import {LIST_TRIPS} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/lib/fetch";
import {Trip} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {Helmet} from "react-helmet";
import {TripRow} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/components/TripRow";

export default function TripsPage() {
    const {data, error} = useSWR<Trip[]>(LIST_TRIPS);

    if (error) {
        console.error(error);
    }

    return (
        <main className="mb-8">
            <Helmet>
                <title>Your Trips</title>
            </Helmet>

            <header className="bg-white shadow">
                <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
                    <h1 className="text-3xl font-bold leading-tight text-gray-900">
                        Your trips
                    </h1>
                </div>
            </header>
            <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                {error && (
                    <p>Error loading trips: {error.toString()}</p>
                )}
                {data && (
                    <div className="flex flex-col">
                        <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                            <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
                                <div className="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
                                    <table className="min-w-full divide-y divide-gray-200">
                                        <thead>
                                        <tr>
                                            <th className="px-6 py-3 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider">
                                                Left at
                                            </th>
                                            <th className="px-6 py-3 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider">
                                                Duration
                                            </th>
                                            <th className="px-6 py-3 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider">
                                                Tags
                                            </th>
                                            <th className="px-6 py-3 bg-gray-50"></th>
                                        </tr>
                                        </thead>
                                        <tbody className="bg-white divide-y divide-gray-200">
                                        {data.map(trip => (
                                            <TripRow key={trip.getId()} trip={trip}/>
                                        ))}
                                        </tbody>
                                    </table>
                                </div>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </main>
    );
}
