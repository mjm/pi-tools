import React from 'react'
import useSWR from 'swr'
import {fetcher} from "../lib/fetch";
import {ListTripsResponse} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";

export default function TripsPage() {
    const {data, error} = useSWR<ListTripsResponse.Trip[]>('ListTrips', fetcher)

    if (error) {
        console.error(error)
    }

    return (
        <div>
            <main>
                <h1>Your Trips</h1>
                {error && (
                    <p>Error loading trips: {error.toString()}</p>
                )}
                {data && (
                    <ul>
                        {data.map(trip => (
                            <li key={trip.getId()}>
                                {trip.getLeftAt()} - {trip.getReturnedAt() || 'now'}
                            </li>
                        ))}
                    </ul>
                )}
            </main>
        </div>
    )
}
