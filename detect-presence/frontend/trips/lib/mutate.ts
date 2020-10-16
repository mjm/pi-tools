import {client} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/lib/trips_client";
import {
    IgnoreTripRequest,
    Trip,
    UpdateTripTagsRequest,
} from "com_github_mjm_pi_tools/detect-presence/proto/trips/trips_pb";
import {mutate} from "swr";
import {GET_TRIP, LIST_TRIPS} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/lib/fetch";

export async function ignoreTrip(id: string): Promise<void> {
    const req = new IgnoreTripRequest();
    req.setId(id);
    return new Promise((resolve, reject) => {
        client.ignoreTrip(req, err => {
            if (err) {
                reject(err);
            } else {
                mutate(LIST_TRIPS);
                resolve();
            }
        });
    });
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

    return new Promise((resolve, reject) => {
        client.updateTripTags(req, err => {
            if (err) {
                reject(err);
            } else {
                mutate(LIST_TRIPS);
                mutate([GET_TRIP, id], (trip: Trip) => {
                    const newTrip = trip.cloneMessage();
                    newTrip.setTagsList(newTags);
                    return newTrip;
                });

                resolve();
            }
        });
    });
}
