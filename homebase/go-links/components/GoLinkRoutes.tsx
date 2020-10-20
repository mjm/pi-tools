import React from "react";
import {Route, Switch, useRouteMatch} from "react-router-dom";
import {SWRConfig} from "swr";
import {GoLinksHomePage} from "com_github_mjm_pi_tools/homebase/go-links/components/GoLinksHomePage";

export function GoLinkRoutes() {
    const {path} = useRouteMatch();

    return (
        <SWRConfig value={{}}>
            <Switch>
                <Route exact path={path}>
                    <GoLinksHomePage/>
                </Route>
            </Switch>
        </SWRConfig>
    );
}
