import React from 'react'
import Head from 'next/head'
import useSWR from 'swr'
import {fetcher} from "../lib/fetch";
import {ListTripsResponse} from "../proto/trips_pb";

export default function Home() {
    const {data, error} = useSWR<ListTripsResponse.Trip[]>('ListTrips', fetcher)

    return (
        <div>
            <Head>
                <title>Your Trips</title>
            </Head>

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
