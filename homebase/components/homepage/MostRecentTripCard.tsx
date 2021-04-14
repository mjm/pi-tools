import React from "react";
import {graphql, useFragment} from "react-relay/hooks";
import {MostRecentTripCard_viewer$key} from "../../__generated__/MostRecentTripCard_viewer.graphql";
import {formatDuration, intervalToDuration, parseISO} from "date-fns";
import HomePageCard from "./HomePageCard";
import {MapIcon} from "@heroicons/react/outline";

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
            icon={<MapIcon className="h-6 w-6 text-gray-400"/>}
            footerHref="/trips"
            footer="View recent trips">
            {formatDuration(intervalToDuration({
                start: parseISO(trip.leftAt),
                end: trip.returnedAt ? parseISO(trip.returnedAt) : new Date(),
            }))}
        </HomePageCard>
    );
}
