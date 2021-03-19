import React from "react";
import {TripsPage} from "com_github_mjm_pi_tools/homebase/trips/components/TripsPage";
import {TripPage} from "com_github_mjm_pi_tools/homebase/trips/components/TripPage";
import {loadQuery} from "react-relay/hooks";
import RelayEnvironment from "com_github_mjm_pi_tools/homebase/lib/environment";
import TripsPageQuery from "com_github_mjm_pi_tools/homebase/api/__generated__/TripsPageQuery.graphql";
import TripPageQuery from "com_github_mjm_pi_tools/homebase/api/__generated__/TripPageQuery.graphql";

export function tripRoutes(path: string): any {
    return [
        {
            path,
            exact: true,
            component: TripsPage,
            prepare() {
                return {
                    tripsQuery: loadQuery(RelayEnvironment, TripsPageQuery, {}),
                };
            },
        },
        {
            path: `${path}/:id`,
            component: TripPage,
            prepare({id}) {
                return {
                    tripQuery: loadQuery(RelayEnvironment, TripPageQuery, {id}),
                };
            },
        },
    ];
}
