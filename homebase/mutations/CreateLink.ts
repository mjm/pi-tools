import {graphql, useMutation} from "react-relay/hooks";
import {CreateLinkInput, CreateLinkMutation} from "../__generated__/CreateLinkMutation.graphql";

export function useCreateLink() {
    const [commit, isInFlight] = useMutation<CreateLinkMutation>(
        graphql`
            mutation CreateLinkMutation($input: CreateLinkInput!, $connections: [ID!]!) {
                createLink(input: $input) {
                    link @prependNode(
                        connections: $connections
                        edgeTypeName: "LinkEdge"
                    ) {
                        id
                        ...LinkRow_link
                    }
                }
            }
        `,
    );

    async function myCommit(input: CreateLinkInput, connections: string[]) {
        return new Promise((resolve, reject) => {
            commit({
                variables: {input, connections},
                onCompleted: resolve,
                onError: reject,
            });
        });
    }

    return [myCommit, isInFlight] as const;
}
