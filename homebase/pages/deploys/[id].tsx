import {graphql, usePreloadedQuery} from "react-relay/hooks";
import {RelayProps} from "relay-nextjs";
import {Id_DeployQuery} from "../../__generated__/Id_DeployQuery.graphql";
import withRelay from "../../lib/withRelay";
import Head from "next/head";
import DescriptionField from "../../components/DescriptionField";
import {format, formatDuration, intervalToDuration, parseISO} from "date-fns";
import Alert from "../../components/Alert";

const DeployQuery = graphql`
    query Id_DeployQuery($id: ID!) {
        viewer {
            deploy(id: $id) {
                id
                commitSHA
                commitMessage
                startedAt
                finishedAt
            }
        }
    }
`;

function DeployPage({preloadedQuery}: RelayProps<{}, Id_DeployQuery>) {
    const query = usePreloadedQuery(DeployQuery, preloadedQuery);

    const deploy = query.viewer.deploy;

    return (
        <main className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
            <Head>
                <title>Deployment Report</title>
            </Head>

            <div className="bg-white shadow overflow-hidden sm:rounded-lg">
                <div className="bg-white px-4 py-5 border-b border-gray-200 sm:px-6">
                    <div className="-ml-4 -mt-2 flex items-center justify-between flex-wrap sm:flex-nowrap">
                        <div className="ml-4 mt-2">
                            <h3 className="text-lg leading-6 font-medium text-gray-900">
                                Deployment
                            </h3>
                        </div>
                    </div>
                </div>
                {deploy ? (
                    <div>
                        <dl>
                            <DescriptionField label="Commit" offset>
                                <a href={`https://github.com/mjm/pi-tools/commit/${deploy.commitSHA}`}
                                   className="font-medium text-indigo-600 hover:text-indigo-500"
                                   target="_blank">
                                    {deploy.commitMessage}
                                </a>
                            </DescriptionField>
                            <DescriptionField label="Started at">
                                {format(parseISO(deploy.startedAt), "PPpp")}
                            </DescriptionField>
                            <DescriptionField label="Finished at" offset>
                                {format(parseISO(deploy.finishedAt), "PPpp")}
                            </DescriptionField>
                            <DescriptionField label="Duration">
                                {formatDuration(intervalToDuration({
                                    start: parseISO(deploy.startedAt),
                                    end: parseISO(deploy.finishedAt),
                                }))}
                            </DescriptionField>
                        </dl>
                    </div>
                ) : (
                    <Alert title="Couldn't load this deploy" severity="error" rounded={false}>
                        No deploy was found with this ID.
                    </Alert>
                )}
            </div>
        </main>
    );
}

export default withRelay(DeployPage, DeployQuery);
