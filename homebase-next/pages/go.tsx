import {RelayProps} from "relay-nextjs";
import {graphql, usePreloadedQuery} from "react-relay/hooks";
import {go_GoLinksHomePageQuery} from "../__generated__/go_GoLinksHomePageQuery.graphql";
import RecentLinksList from "../components/go-links/RecentLinksList";
import withRelay from "../lib/withRelay";
import PageHeader from "../components/PageHeader";
import NewLinkCard from "../components/go-links/NewLinkCard";

const GoLinksQuery = graphql`
    query go_GoLinksHomePageQuery {
        viewer {
            links(first: 30) @connection(key: "RecentLinksList_links") {
                __id
                ...RecentLinksList_links

                # Not used here but it keeps the relay-compiler happy
                edges {
                    __id
                }
            }
        }
    }
`;

function GoLinksHomePage({preloadedQuery}: RelayProps<{}, go_GoLinksHomePageQuery>) {
    const query = usePreloadedQuery(GoLinksQuery, preloadedQuery);

    return (
        <main className="mb-8">
            <PageHeader>
                Go links
            </PageHeader>
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-6">
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-8">
                    <div><NewLinkCard connections={[query.viewer.links.__id]}/></div>
                    <div><RecentLinksList links={query.viewer.links}/></div>
                </div>
            </div>
        </main>
    );
}

export default withRelay(GoLinksHomePage, GoLinksQuery);
