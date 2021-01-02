import {graphql, useMutation} from "react-relay/hooks";
import {CreateLinkInput} from "com_github_mjm_pi_tools/homebase/api/__generated__/CreateLinkMutation.graphql";

export function useCreateLink() {
    const [commit, isInFlight] = useMutation(
        graphql`
            mutation CreateLinkMutation($input: CreateLinkInput!) {
                createLink(input: $input) {
                    link @prependNode(
                        connections: ["client:root:viewer:__RecentLinksList_links_connection"]
                        edgeTypeName: "LinkEdge"
                    ) {
                        id
                        ...LinkRow_link
                    }
                }
            }
        `,
    );

    async function myCommit(input: CreateLinkInput) {
        return new Promise((resolve, reject) => {
            commit({
                variables: {input},
                onCompleted: resolve,
                onError: reject,
            });
        });
    }

    return [myCommit, isInFlight] as const;
}
