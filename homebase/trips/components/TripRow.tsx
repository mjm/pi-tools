import React from "react";
import {format, formatDuration, intervalToDuration, parseISO} from "date-fns";
import {Link} from "react-router-dom";
import {TripTag} from "com_github_mjm_pi_tools/homebase/trips/components/TripTag";
import {graphql, useFragment} from "react-relay/hooks";
import {TripRow_trip$key} from "com_github_mjm_pi_tools/homebase/api/__generated__/TripRow_trip.graphql";

export function TripRow({trip}: { trip: TripRow_trip$key }) {
    const data = useFragment(
        graphql`
            fragment TripRow_trip on Trip {
                id
                leftAt
                returnedAt
                tags
            }
        `,
        trip,
    );

    const leftAt = parseISO(data.leftAt);
    let duration = null;
    if (data.returnedAt) {
        duration = intervalToDuration({
            start: leftAt,
            end: parseISO(data.returnedAt),
        });
    }

    return (
        <tr>
            <td className="px-6 py-4 whitespace-nowrap text-sm leading-5 font-medium text-gray-900">
                {format(leftAt, "PPpp")}
            </td>
            <td className="px-6 py-4 whitespace-nowrap text-sm leading-5 text-gray-500">
                {duration ? formatDuration(duration) : "Ongoing"}
            </td>
            <td className="px-6 py-4 whitespace-nowrap">
                <div className="flex flex-row space-x-3">
                    {data.tags.map(tag => (
                        <TripTag key={tag}>
                            {tag}
                        </TripTag>
                    ))}
                </div>
            </td>
            <td className="px-6 py-4 whitespace-nowrap text-right text-sm leading-5 font-medium">
                <Link to={`/trips/${data.id}`}
                      className="text-indigo-600 hover:text-indigo-900">
                    Details
                </Link>
            </td>
        </tr>
    );
}
