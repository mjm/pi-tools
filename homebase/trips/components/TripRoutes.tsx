import React from "react";
import {Route, Switch, useRouteMatch} from "react-router-dom";
import {TripsPage} from "com_github_mjm_pi_tools/homebase/trips/components/TripsPage";
import {TripPage} from "com_github_mjm_pi_tools/homebase/trips/components/TripPage";

export function TripRoutes() {
    const {path} = useRouteMatch();

    return (
        <Switch>
            <Route exact path={path}>
                <TripsPage/>
            </Route>
            <Route path={`${path}/:id`}>
                <TripPage/>
            </Route>
        </Switch>
    );
}
