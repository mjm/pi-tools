import React from "react";
import {Route, Switch, useRouteMatch} from "react-router-dom";
import {SWRConfig} from "swr";
import {TripsPage} from "com_github_mjm_pi_tools/homebase/trips/components/TripsPage";
import {TripPage} from "com_github_mjm_pi_tools/homebase/trips/components/TripPage";
import {fetcher} from "com_github_mjm_pi_tools/homebase/trips/lib/fetch";

export function TripRoutes() {
    const {path} = useRouteMatch();

    return (
        <SWRConfig value={{
            fetcher,
        }}>
            <Switch>
                <Route exact path={path}>
                    <TripsPage/>
                </Route>
                <Route path={`${path}/:id`}>
                    <TripPage/>
                </Route>
            </Switch>
        </SWRConfig>
    );
}
