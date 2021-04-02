import PageHeader from "../components/PageHeader";
import {RelayProps} from "relay-nextjs";
import Head from "next/head";
import {graphql, usePreloadedQuery} from "react-relay/hooks";
import withRelay from "../lib/withRelay";
import {trips_TripsPageQuery} from "../__generated__/trips_TripsPageQuery.graphql";
import TagFilters from "../components/trips/TagFilters";
import TripRow from "../components/trips/TripRow";

const TripsPageQuery = graphql`
    query trips_TripsPageQuery {
        viewer {
            ...TagFilters_tags
            trips(first: 30) @connection(key: "TripsPageQuery_trips") {
                edges {
                    node {
                        id
                        ...TripRow_trip
                    }
                }
            }
        }
    }
`;

function TripsPage({preloadedQuery}: RelayProps<{}, trips_TripsPageQuery>) {
    const query = usePreloadedQuery(TripsPageQuery, preloadedQuery);
    if (!query) {
        return null;
    }

    const tripNodes = query.viewer.trips.edges.map(e => e.node);

    return (
        <main className="mb-8">
            <Head>
                <title>Your Trips</title>
            </Head>

            <PageHeader buttons={<>
                <a href="/app/download"
                   className="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                    <svg className="-ml-1 mr-2 h-5 w-5 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none"
                         viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                              d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"/>
                    </svg>
                    Download app
                </a>
            </>}>
                Your trips
            </PageHeader>
            <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                <div className="flex flex-col">
                    <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                        <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
                            <div className="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
                                <table className="min-w-full divide-y divide-gray-200">
                                    <thead>
                                    <tr>
                                        <td colSpan={4}
                                            className="px-6 py-3 bg-gray-50 text-xs leading-4 text-gray-500 border-b border-gray-200">
                                            <TagFilters tags={query.viewer}/>
                                        </td>
                                    </tr>
                                    <tr>
                                        <th className="px-6 py-3 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider">
                                            Left at
                                        </th>
                                        <th className="px-6 py-3 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider">
                                            Duration
                                        </th>
                                        <th className="px-6 py-3 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider">
                                            Tags
                                        </th>
                                        <th className="px-6 py-3 bg-gray-50"/>
                                    </tr>
                                    </thead>
                                    <tbody className="bg-white divide-y divide-gray-200">
                                    {tripNodes.map(trip => (
                                        <TripRow key={trip.id} trip={trip}/>
                                    ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </main>
    );
}

export default withRelay(TripsPage, TripsPageQuery)
