import React from "react";
import {Helmet} from "react-helmet";
import {useHistory, useParams} from "react-router-dom";
import useSWR from "swr";
import {format, formatDuration, intervalToDuration, parseISO} from "date-fns";
import {DescriptionField} from "com_github_mjm_pi_tools/homebase/components/DescriptionField";
import {Trip} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {GET_TRIP} from "com_github_mjm_pi_tools/homebase/trips/lib/fetch";
import {ignoreTrip} from "com_github_mjm_pi_tools/homebase/trips/lib/mutate";
import {TripTagField} from "com_github_mjm_pi_tools/homebase/trips/components/TripTagField";
import {Alert} from "com_github_mjm_pi_tools/homebase/components/Alert";

export function TripPage() {
    const {id} = useParams<{ id: string }>();
    const {data, error} = useSWR<Trip>([GET_TRIP, id]);

    if (error) {
        console.error(error);
    }

    return (
        <main className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
            <Helmet>
                <title>Trip Details</title>
            </Helmet>

            <div className="bg-white shadow overflow-hidden sm:rounded-lg">
                <div className="bg-white px-4 py-5 border-b border-gray-200 sm:px-6">
                    <div className="-ml-4 -mt-2 flex items-center justify-between flex-wrap sm:flex-no-wrap">
                        <div className="ml-4 mt-2">
                            <h3 className="text-lg leading-6 font-medium text-gray-900">
                                Trip Details
                            </h3>
                        </div>
                        <div className="ml-4 mt-2 flex-shrink-0 flex">
                            {data && <TripIgnoreButton id={data.getId()}/>}
                        </div>
                    </div>
                </div>
                {error && (
                    <Alert severity="error" rounded={false} title="Couldn't load this trip">
                        {error.toString()}
                    </Alert>
                )}
                {data && (
                    <div>
                        <dl>
                            <DescriptionField label="Left at" offset>
                                {format(parseISO(data.getLeftAt()), "PPpp")}
                            </DescriptionField>
                            {data.getReturnedAt() && (
                                <>
                                    <DescriptionField label="Returned at">
                                        {format(parseISO(data.getReturnedAt()), "PPpp")}
                                    </DescriptionField>
                                    <DescriptionField label="Duration" offset>
                                        {formatDuration(intervalToDuration({
                                            start: parseISO(data.getLeftAt()),
                                            end: parseISO(data.getReturnedAt()),
                                        }))}
                                    </DescriptionField>
                                </>
                            )}
                            <DescriptionField label="Tags">
                                <TripTagField trip={data}/>
                            </DescriptionField>
                        </dl>
                    </div>
                )}
            </div>
        </main>
    );
}

function TripIgnoreButton({id}: { id: string }) {
    const history = useHistory();
    const [isIgnoring, setIgnoring] = React.useState(false);

    async function onIgnore() {
        setIgnoring(true);
        try {
            await ignoreTrip(id);

            // return to the trips page upon successful ignore
            history.push("/trips");
        } catch (e) {
            console.error(e);
        } finally {
            setIgnoring(false);
        }
    }

    return (
        <span className="inline-flex rounded-md shadow-sm">
            <button type="button"
                    disabled={isIgnoring}
                    onClick={onIgnore}
                    className="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm leading-5 font-medium rounded-md text-gray-700 bg-white hover:text-gray-500 focus:outline-none focus:shadow-outline-blue focus:border-blue-300 active:bg-gray-50 active:text-gray-800">
                <svg className="-ml-1 mr-2 h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg"
                     viewBox="0 0 20 20"
                     fill="currentColor">
                    <path fillRule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zM7 9a1 1 0 000 2h6a1 1 0 100-2H7z"
                          clipRule="evenodd"/>
                </svg>
                <span>Ignore</span>
            </button>
        </span>
    );
}
