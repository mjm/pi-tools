import React from "react";
import {graphql, useFragment} from "react-relay/hooks";
import {MostRecentTripCard_viewer$key} from "../../__generated__/MostRecentTripCard_viewer.graphql";
import {formatDuration, intervalToDuration, parseISO} from "date-fns";
import HomePageCard from "./HomePageCard";

export default function MostRecentTripCard({viewer}: { viewer: MostRecentTripCard_viewer$key }) {
    const data = useFragment(
        graphql`
            fragment MostRecentTripCard_viewer on Viewer {
                trips(first: 1) {
                    edges {
                        node {
                            leftAt
                            returnedAt
                        }
                    }
                }
            }
        `,
        viewer,
    );

    const trip = data.trips.edges[0].node;

    return (
        <HomePageCard
            title={trip.returnedAt ? <>Most recent trip</> : <>Current trip</>}
            icon={
                <svg className="h-6 w-6 text-gray-400" xmlns="http://www.w3.org/2000/svg"
                     fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                          d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7"/>
                </svg>
            }
            footerHref="/trips"
            footer="View recent trips">
            {formatDuration(intervalToDuration({
                start: parseISO(trip.leftAt),
                end: trip.returnedAt ? parseISO(trip.returnedAt) : new Date(),
            }))}
        </HomePageCard>
    );
}
