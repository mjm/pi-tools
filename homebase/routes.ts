import {HomePage} from "com_github_mjm_pi_tools/homebase/homepage/components/HomePage";
import {tripRoutes} from "com_github_mjm_pi_tools/homebase/trips/components/TripRoutes";
import {goLinksRoutes} from "com_github_mjm_pi_tools/homebase/go-links/components/GoLinkRoutes";
import {backupRoutes} from "com_github_mjm_pi_tools/homebase/backups/components/BackupRoutes";

export default [
    {
        path: "/",
        exact: true,
        component: HomePage,
    },
    ...tripRoutes("/trips"),
    ...goLinksRoutes("/go"),
    ...backupRoutes("/backups"),
];
