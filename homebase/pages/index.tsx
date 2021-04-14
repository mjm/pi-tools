import {RelayProps} from "relay-nextjs";
import {graphql, usePreloadedQuery} from "react-relay/hooks";
import {pages_HomePageQuery} from "../__generated__/pages_HomePageQuery.graphql";
import PageHeader from "../components/PageHeader";
import withRelay from "../lib/withRelay";
import MostRecentTripCard from "../components/homepage/MostRecentTripCard";
import MostRecentDeployCard from "../components/homepage/MostRecentDeployCard";
import FiringAlertsCard from "../components/homepage/FiringAlertsCard";
import Alert from "../components/Alert";

const HomePageQuery = graphql`
    query pages_HomePageQuery {
        viewer {
            ...MostRecentTripCard_viewer
            ...FiringAlertsCard_viewer
            ...MostRecentDeployCard_viewer
        }
    }
`;

function HomePage(props: RelayProps<{}, pages_HomePageQuery>) {
    const query = usePreloadedQuery(HomePageQuery, props.preloadedQuery);

    if (props.err) {
        console.error(props.err);
    }

    return (
        <main className="mb-8">
            <PageHeader>Homebase</PageHeader>
            <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
                    <MostRecentTripCard viewer={query.viewer}/>
                    <FiringAlertsCard viewer={query.viewer}/>
                    <MostRecentDeployCard viewer={query.viewer}/>
                </div>
            </div>
        </main>
    );
}

export default withRelay(HomePage, HomePageQuery);
