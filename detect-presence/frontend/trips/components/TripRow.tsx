import React from "react";
import {format, formatDuration, intervalToDuration, parseISO} from "date-fns";
import {Link} from "react-router-dom";
import {Trip} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {TripTag} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/components/TripTag";

export function TripRow({trip}: { trip: Trip }) {
    const leftAt = parseISO(trip.getLeftAt());
    let duration = null;
    if (trip.getReturnedAt()) {
        duration = intervalToDuration({
            start: leftAt,
            end: parseISO(trip.getReturnedAt()),
        });
    }

    return (
        <tr key={trip.getId()}>
            <td className="px-6 py-4 whitespace-no-wrap text-sm leading-5 font-medium text-gray-900">
                {format(leftAt, "PPpp")}
            </td>
            <td className="px-6 py-4 whitespace-no-wrap text-sm leading-5 text-gray-500">
                {duration ? formatDuration(duration) : "Ongoing"}
            </td>
            <td className="px-6 py-4 whitespace-no-wrap">
                <div className="flex flex-row space-x-3">
                    {trip.getTagsList().map(tag => (
                        <TripTag key={tag} tag={tag}/>
                    ))}
                </div>
            </td>
            <td className="px-6 py-4 whitespace-no-wrap text-right text-sm leading-5 font-medium">
                <Link to={`/trips/${trip.getId()}`}
                      className="text-indigo-600 hover:text-indigo-900">
                    Details
                </Link>
            </td>
        </tr>
    );
}
