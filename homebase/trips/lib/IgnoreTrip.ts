import {graphql, useMutation} from "react-relay/hooks";
import {IgnoreTripMutation} from "com_github_mjm_pi_tools/homebase/api/__generated__/IgnoreTripMutation.graphql";

export function useIgnoreTrip() {
    const [commit, isInFlight] = useMutation<IgnoreTripMutation>(
        graphql`
            mutation IgnoreTripMutation($input: IgnoreTripInput!, $connections: [ID!]!) {
                ignoreTrip(input: $input) {
                    ignoredTripID @deleteEdge(connections: $connections)
                }
            }
        `,
    );

    async function myCommit(id: string) {
        return new Promise((resolve, reject) => {
            commit({
                variables: {
                    input: {id},
                    connections: ["TripsPageQuery_trips"],
                },
                onCompleted: resolve,
                onError: reject,
            });
        });
    }

    return [myCommit, isInFlight] as const;
}
