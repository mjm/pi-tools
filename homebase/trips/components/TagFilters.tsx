import React from "react";
import useSWR from "swr";
import {LIST_TAGS} from "com_github_mjm_pi_tools/homebase/trips/lib/fetch";
import {Tag} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {TripTag} from "com_github_mjm_pi_tools/homebase/trips/components/TripTag";

export function TagFilters() {
    const {data, error} = useSWR<Tag[]>([LIST_TAGS, 5]);

    if (error) {
        console.error(error);
    }

    return (
        <>
            <div className="flex flex-row items-baseline space-x-3 overflow-hidden whitespace-no-wrap flex-wrap">
                <span className="font-medium uppercase tracking-wider">Popular Tags:</span>
                {data ? data.map(tag => (
                    <TripTag key={tag.getName()}>
                        {tag.getName()} ({tag.getTripCount()})
                    </TripTag>
                )) : null}
            </div>
        </>
    );
}
