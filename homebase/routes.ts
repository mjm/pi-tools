import {HomePage} from "com_github_mjm_pi_tools/homebase/homepage/components/HomePage";
import {tripRoutes} from "com_github_mjm_pi_tools/homebase/trips/components/TripRoutes";
import {goLinksRoutes} from "com_github_mjm_pi_tools/homebase/go-links/components/GoLinkRoutes";
import {backupRoutes} from "com_github_mjm_pi_tools/homebase/backups/components/BackupRoutes";
import {loadQuery} from "react-relay/hooks";
import RelayEnvironment from "com_github_mjm_pi_tools/homebase/lib/environment";
import HomePageQuery from "com_github_mjm_pi_tools/homebase/api/__generated__/HomePageQuery.graphql";

export default [
    {
        path: "/",
        exact: true,
        component: HomePage,
        prepare() {
            return {
                homeQuery: loadQuery(RelayEnvironment, HomePageQuery, {}),
            };
        },
    },
    ...tripRoutes("/trips"),
    ...goLinksRoutes("/go"),
    ...backupRoutes("/backups"),
];
