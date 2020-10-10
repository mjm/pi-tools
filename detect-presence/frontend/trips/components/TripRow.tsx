import React from "react";
import {ListTripsResponse} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {format, formatDuration, intervalToDuration, parseISO} from "date-fns";

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
        <div className="bg-gray-100 border-b border-gray-200 p-3">
            <span>{format(leftAt, "PPpp")}</span>{" "}
            {duration && (
                <>({formatDuration(duration)})</>
            )}
        </div>
    );
}
