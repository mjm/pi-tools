import React from "react";
import {BackupsPage} from "com_github_mjm_pi_tools/homebase/backups/components/BackupsPage";
import RelayEnvironment from "com_github_mjm_pi_tools/homebase/lib/environment";
import {loadQuery} from "react-relay/hooks";
import BackupsPageQuery from "com_github_mjm_pi_tools/homebase/api/__generated__/BackupsPageQuery.graphql";

export function backupRoutes(path: string): any {
    return [
        {
            path,
            exact: true,
            component: BackupsPage,
            prepare() {
                return {
                    backupsQuery: loadQuery(RelayEnvironment, BackupsPageQuery, {}),
                };
            },
        },
    ];
}
