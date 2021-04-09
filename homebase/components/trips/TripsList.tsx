import {graphql, usePaginationFragment} from "react-relay";
import TripRow from "./TripRow";
import {TripsList_viewer$key} from "../../__generated__/TripsList_viewer.graphql";

export default function TripsList({viewer}: { viewer: TripsList_viewer$key }) {
    const {data} = usePaginationFragment(
        graphql`
            fragment TripsList_viewer on Viewer
            @refetchable(queryName: "TripsListPaginationQuery")
            @argumentDefinitions(
                count: {type: "Int", defaultValue: 30}
                cursor: {type: "Cursor"}
            ) {
                trips(first: $count, after: $cursor)
                @connection(key: "TripsList_trips") {
                    edges {
                        node {
                            id
                            ...TripRow_trip
                        }
                    }
                }
            }
        `,
        viewer,
    );

    const tripNodes = data.trips.edges.map(e => e.node);

    return (
        <tbody className="bg-white divide-y divide-gray-200">
        {tripNodes.map(trip => (
            <TripRow key={trip.id} trip={trip}/>
        ))}
        </tbody>
    );
}
