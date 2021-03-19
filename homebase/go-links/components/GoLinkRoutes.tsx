import React from "react";
import {Route, Switch, useRouteMatch} from "react-router-dom";
import {GoLinksHomePage} from "com_github_mjm_pi_tools/homebase/go-links/components/GoLinksHomePage";
import {GoLinkDetailPage} from "com_github_mjm_pi_tools/homebase/go-links/components/GoLinkDetailPage";

export function GoLinkRoutes() {
    const {path} = useRouteMatch();

    return (
        <Switch>
            <Route exact path={path}>
                <GoLinksHomePage/>
            </Route>
            <Route path={`${path}/:id`}>
                <GoLinkDetailPage/>
            </Route>
        </Switch>
    );
}

export function goLinksRoutes(path: string) {
    return [
        {
            path,
            exact: true,
            component: GoLinksHomePage,
        },
        {
            path: `${path}/:id`,
            component: GoLinkDetailPage,
        },
    ];
}
