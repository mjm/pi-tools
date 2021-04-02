import {graphql, useMutation} from "react-relay/hooks";
import {UpdateLinkInput, UpdateLinkMutation} from "../__generated__/UpdateLinkMutation.graphql";

export function useUpdateLink() {
    const [commit, isInFlight] = useMutation<UpdateLinkMutation>(
        graphql`
            mutation UpdateLinkMutation($input: UpdateLinkInput!) {
                updateLink(input: $input) {
                    link {
                        id
                        ...LinkRow_link
                    }
                }
            }
        `,
    );

    async function myCommit(input: UpdateLinkInput) {
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
