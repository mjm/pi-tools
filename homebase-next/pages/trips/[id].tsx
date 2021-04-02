import {RelayProps} from "relay-nextjs";
import {graphql, usePreloadedQuery} from "react-relay/hooks";
import {Id_TripPageQuery} from "../../__generated__/Id_TripPageQuery.graphql";
import Head from "next/head";
import {format, formatDuration, intervalToDuration, parseISO} from "date-fns";
import TripIgnoreButton from "../../components/trips/TripIgnoreButton";
import DescriptionField from "../../components/DescriptionField";
import TripTagField from "../../components/trips/TripTagField";
import Alert from "../../components/Alert";
import withRelay from "../../lib/withRelay";

const TripPageQuery = graphql`
    query Id_TripPageQuery($id: ID!) {
        viewer {
            trip(id: $id) {
                id
                leftAt
                returnedAt
                ...TripTagField_trip
            }
        }
    }
`;

function TripPage({preloadedQuery}: RelayProps<{}, Id_TripPageQuery>) {
    const query = usePreloadedQuery(TripPageQuery, preloadedQuery);
    const trip = query.viewer.trip;

    return (
        <main className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
            <Head>
                <title>Trip Details</title>
            </Head>

            <div className="bg-white shadow overflow-hidden sm:rounded-lg">
                <div className="bg-white px-4 py-5 border-b border-gray-200 sm:px-6">
                    <div className="-ml-4 -mt-2 flex items-center justify-between flex-wrap sm:flex-nowrap">
                        <div className="ml-4 mt-2">
                            <h3 className="text-lg leading-6 font-medium text-gray-900">
                                Trip Details
                            </h3>
                        </div>
                        <div className="ml-4 mt-2 flex-shrink-0 flex">
                            {trip && <TripIgnoreButton id={trip.id}/>}
                        </div>
                    </div>
                </div>
                {trip ? (
                    <div>
                        <dl>
                            <DescriptionField label="Left at" offset>
                                {format(parseISO(trip.leftAt), "PPpp")}
                            </DescriptionField>
                            {trip.returnedAt && (
                                <>
                                    <DescriptionField label="Returned at">
                                        {format(parseISO(trip.returnedAt), "PPpp")}
                                    </DescriptionField>
                                    <DescriptionField label="Duration" offset>
                                        {formatDuration(intervalToDuration({
                                            start: parseISO(trip.leftAt),
                                            end: parseISO(trip.returnedAt),
                                        }))}
                                    </DescriptionField>
                                </>
                            )}
                            <DescriptionField label="Tags">
                                <TripTagField trip={trip}/>
                            </DescriptionField>
                        </dl>
                    </div>
                ) : (
                    <Alert title="Couldn't load this trip" severity="error" rounded={false}>
                        No trip was found with this ID.
                    </Alert>
                )}
            </div>
        </main>
    );
}

export default withRelay(TripPage, TripPageQuery);
