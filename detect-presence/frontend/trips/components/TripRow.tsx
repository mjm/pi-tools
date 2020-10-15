import React from "react";
import {ListTripsResponse} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {format, formatDuration, intervalToDuration, parseISO} from "date-fns";
import {TripRowActions} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/components/TripRowActions";
import {TripTag} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/components/TripTag";

export function TripRow({trip}: { trip: ListTripsResponse.Trip }) {
    const leftAt = parseISO(trip.getLeftAt());
    let duration = null;
    if (trip.getReturnedAt()) {
        duration = intervalToDuration({
            start: leftAt,
            end: parseISO(trip.getReturnedAt()),
        });
    }

    return (
        <div className="flex items-baseline bg-gray-100 border-b border-gray-200 p-3">
            <div>
                <span>{format(leftAt, "PPpp")}</span>{" "}
                {duration && (
                    <>({formatDuration(duration)})</>
                )}
            </div>
            <div className="pl-3">
                {trip.getTagsList().map(tag => (
                    <TripTag tag={tag} key={tag} />
                ))}
            </div>
            <div className="flex-none ml-auto">
                <TripRowActions trip={trip}/>
            </div>
        </div>
    );
}
