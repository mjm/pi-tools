import {graphql, usePreloadedQuery} from "react-relay/hooks";
import {RelayProps} from "relay-nextjs";
import {Id_GoLinkDetailPageQuery} from "../../__generated__/Id_GoLinkDetailPageQuery.graphql";
import Head from "next/head";
import EditLinkForm from "../../components/go-links/EditLinkForm";
import Alert from "../../components/Alert";
import withRelay from "../../lib/withRelay";

const GoLinkQuery = graphql`
    query Id_GoLinkDetailPageQuery($id: ID!) {
        viewer {
            link(id: $id) {
                id
                shortURL
                ...EditLinkForm_link
            }
        }
    }
`;

function GoLinkDetailPage({preloadedQuery}: RelayProps<{}, Id_GoLinkDetailPageQuery>) {
    const query = usePreloadedQuery(GoLinkQuery, preloadedQuery);
    const link = query.viewer.link;

    return (
        <main className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
            <Head>
                <title>{`Link Details${link ? `: go/${link.shortURL}` : ""}`}</title>
            </Head>

            <div className="bg-white shadow overflow-hidden sm:rounded-lg">
                <div className="px-4 py-5 sm:px-6">
                    <div className="-ml-4 -mt-2 flex items-center justify-between flex-wrap sm:flex-nowrap">
                        <div className="ml-4 mt-2">
                            <h3 className="text-lg leading-6 font-medium text-gray-900">
                                go/{link ? link.shortURL : "…"}
                            </h3>
                        </div>
                    </div>
                </div>
                {link ? (
                    <EditLinkForm link={link}/>
                ) : (
                    <Alert title="Couldn't load link details" severity="error" rounded={false}>
                        No link was found with this ID.
                    </Alert>
                )}
            </div>
        </main>
    );
}

export default withRelay(GoLinkDetailPage, GoLinkQuery);
