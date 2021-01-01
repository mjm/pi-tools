import {graphql, useMutation} from "react-relay/hooks";
import {UpdateTripTagsMutation} from "com_github_mjm_pi_tools/homebase/api/__generated__/UpdateTripTagsMutation.graphql";

export function useUpdateTripTags() {
    const [commit, isInFlight] = useMutation<UpdateTripTagsMutation>(
        graphql`
            mutation UpdateTripTagsMutation($input: UpdateTripTagsInput!) {
                updateTripTags(input: $input) {
                    trip {
                        id
                        tags
                    }
                }
            }
        `,
    );

    async function myCommit(id: string, oldTags: readonly string[], newTags: readonly string[]) {
        const oldTagsSet = new Set(oldTags);
        const newTagsSet = new Set(newTags);

        const tagsToAdd = newTags.filter(tag => !oldTagsSet.has(tag));
        const tagsToRemove = oldTags.filter(tag => !newTagsSet.has(tag));

        return new Promise((resolve, reject) => {
            commit({
                variables: {
                    input: {
                        tripID: id,
                        tagsToAdd,
                        tagsToRemove,
                    },
                },
                onCompleted: resolve,
                onError: reject,
            });
        });
    }

    return [myCommit, isInFlight] as const;
}
