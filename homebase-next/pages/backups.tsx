import Head from "next/head";
import PageHeader from "../components/PageHeader";
import {RelayProps} from "relay-nextjs";
import {graphql, usePreloadedQuery} from "react-relay/hooks";
import withRelay from "../lib/withRelay";
import {backups_BackupsPageQuery} from "../__generated__/backups_BackupsPageQuery.graphql";
import BackupsList from "../components/backups/BackupsList";

const BackupsQuery = graphql`
    query backups_BackupsPageQuery {
        viewer {
            ...BackupsList_viewer
        }
    }
`;

function BackupsPage({preloadedQuery}: RelayProps<{}, backups_BackupsPageQuery>) {
    const query = usePreloadedQuery(BackupsQuery, preloadedQuery);

    return (
        <main className="mb-8">
            <Head>
                <title>Backups</title>
            </Head>

            <PageHeader>
                Recent archives
            </PageHeader>

            <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                {query.viewer && <BackupsList viewer={query.viewer}/>}
            </div>
        </main>
    );
}

export default withRelay(BackupsPage, BackupsQuery);
