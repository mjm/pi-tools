import {graphql, usePreloadedQuery} from "react-relay";
import {RelayProps} from "relay-nextjs";
import {deploys_DeploysQuery} from "../__generated__/deploys_DeploysQuery.graphql";
import Head from "next/head";
import PageHeader from "../components/PageHeader";
import withRelay from "../lib/withRelay";
import RecentDeployments from "../components/deploys/RecentDeployments";

const DeploysQuery = graphql`
    query deploys_DeploysQuery {
        viewer {
            ...RecentDeployments_viewer
        }
    }
`;

function DeploysPage({preloadedQuery}: RelayProps<{}, deploys_DeploysQuery>) {
    const data = usePreloadedQuery(DeploysQuery, preloadedQuery);

    return (
        <main className="mb-8">
            <Head>
                <title>Recent Deployments</title>
            </Head>

            <PageHeader>
                Recent deployments
            </PageHeader>
            <div className="container mx-auto sm:px-6 lg:px-8 py-6">
                <RecentDeployments viewer={data.viewer}/>
            </div>
        </main>
    );
}

export default withRelay(DeploysPage, DeploysQuery);
