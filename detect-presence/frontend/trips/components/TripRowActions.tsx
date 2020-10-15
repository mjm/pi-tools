import React from "react";
import {Trip} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {ignoreTrip} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/lib/mutate";

export function TripRowActions({trip}: { trip: Trip }) {
    const [isIgnoring, setIgnoring] = React.useState(false);

    async function onIgnore() {
        setIgnoring(true);
        try {
            await ignoreTrip(trip.getId());
        } catch (e) {
            console.error(e);
        } finally {
            setIgnoring(false);
        }
    }

    return (
        <div>
            <button
                className="text-sm font-semibold bg-indigo-200 hover:bg-indigo-300 text-indigo-900 px-3 py-1 rounded"
                onClick={onIgnore}
                disabled={isIgnoring}
            >
                Ignore
            </button>
        </div>
    );
}
