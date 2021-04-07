import {graphql, usePaginationFragment} from "react-relay";
import {RecentDeployments_viewer$key} from "../../__generated__/RecentDeployments_viewer.graphql";
import DeploymentRow from "./DeploymentRow";

export default function RecentDeployments({viewer}: { viewer: RecentDeployments_viewer$key }) {
    const {data} = usePaginationFragment(
        graphql`
            fragment RecentDeployments_viewer on Viewer
            @refetchable(queryName: "RecentDeploymentsPaginationQuery")
            @argumentDefinitions(
                count: {type: "Int"}
                cursor: {type: "Cursor"}
            ) {
                recentDeploys(first: $count, after: $cursor)
                @connection(key: "RecentDeployments_recentDeploys") {
                    edges {
                        node {
                            id
                            ...DeploymentRow_deploy
                        }
                    }
                }
            }
        `,
        viewer,
    );

    const deployEdges = data.recentDeploys.edges;

    return (
        <div className="bg-white overflow-hidden shadow sm:rounded-lg">
            {/*<div className="px-4 py-5 sm:p-6">*/}
                <div className="flow-root py-5 px-4 sm:px-6">
                    <ul className="-mb-8">
                        {deployEdges.map((edge, idx) => (
                            <DeploymentRow
                                key={edge.node.id}
                                deploy={edge.node}
                                isLast={idx === deployEdges.length - 1}
                            />
                        ))}
                    </ul>
                </div>
            {/*</div>*/}
        </div>
    );
}
