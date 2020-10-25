import React from "react";
import {Route, Switch, useRouteMatch} from "react-router-dom";
import {SWRConfig} from "swr";
import {GoLinksHomePage} from "com_github_mjm_pi_tools/homebase/go-links/components/GoLinksHomePage";
import {fetcher} from "com_github_mjm_pi_tools/homebase/go-links/lib/fetch";
import {GoLinkDetailPage} from "com_github_mjm_pi_tools/homebase/go-links/components/GoLinkDetailPage";

export function GoLinkRoutes() {
    const {path} = useRouteMatch();

    return (
        <SWRConfig value={{fetcher}}>
            <Switch>
                <Route exact path={path}>
                    <GoLinksHomePage/>
                </Route>
                <Route path={`${path}/:id`}>
                    <GoLinkDetailPage/>
                </Route>
            </Switch>
        </SWRConfig>
    );
}
