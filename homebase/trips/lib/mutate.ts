import {mutate} from "swr";
import {client} from "com_github_mjm_pi_tools/homebase/trips/lib/trips_client";
import {
    IgnoreTripRequest,
    Trip,
    UpdateTripTagsRequest,
} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {GET_TRIP, LIST_TRIPS} from "com_github_mjm_pi_tools/homebase/trips/lib/fetch";
import {promisify} from "com_github_mjm_pi_tools/homebase/lib/promisify";

export async function ignoreTrip(id: string): Promise<void> {
    const req = new IgnoreTripRequest();
    req.setId(id);
    await promisify(client, "ignoreTrip")(req);
    await mutate(LIST_TRIPS);
}

export async function updateTripTags(id: string, oldTags: string[], newTags: string[]): Promise<void> {
    const oldTagsSet = new Set(oldTags);
    const newTagsSet = new Set(newTags);

    const tagsToAdd = newTags.filter(tag => !oldTagsSet.has(tag));
    const tagsToRemove = oldTags.filter(tag => !newTagsSet.has(tag));

    const req = new UpdateTripTagsRequest();
    req.setTripId(id);
    req.setTagsToAddList(tagsToAdd);
    req.setTagsToRemoveList(tagsToRemove);

    await promisify(client, "updateTripTags")(req);

    await Promise.all([
        mutate(LIST_TRIPS),
        mutate([GET_TRIP, id], (trip: Trip) => {
            const newTrip = trip.cloneMessage();
            newTrip.setTagsList(newTags);
            return newTrip;
        }),
    ]);
}
