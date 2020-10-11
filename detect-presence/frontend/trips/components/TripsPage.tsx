import React from "react";
import useSWR from "swr";
import {fetcher} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/lib/fetch";
import {ListTripsResponse} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {TripRow} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/components/TripRow";

export default function TripsPage() {
    const {data, error} = useSWR<ListTripsResponse.Trip[]>("ListTrips", fetcher);

    if (error) {
        console.error(error);
    }

    return (
        <main>
            <h1 className="text-2xl font-bold mb-6">Your Trips</h1>
            {error && (
                <p>Error loading trips: {error.toString()}</p>
            )}
            {data && (
                <ul className="border-t border-gray-200">
                    {data.map(trip => (
                        <li key={trip.getId()}>
                            <TripRow trip={trip} />
                        </li>
                    ))}
                </ul>
            )}
        </main>
    );
}