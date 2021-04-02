import {format, formatDuration, intervalToDuration, parseISO} from "date-fns";
import TripTag from "./TripTag";
import {graphql, useFragment} from "react-relay/hooks";
import {TripRow_trip$key} from "../../__generated__/TripRow_trip.graphql";
import Link from "next/link";

export default function TripRow({trip}: { trip: TripRow_trip$key }) {
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
                <Link href={`/trips/${data.id}`}>
                    <a
                        className="text-indigo-600 hover:text-indigo-900">
                        Details
                    </a>
                </Link>
            </td>
        </tr>
    );
}
