import React from "react";
import {Route, Switch, useRouteMatch} from "react-router-dom";
import {BackupsPage} from "com_github_mjm_pi_tools/homebase/backups/components/BackupsPage";

export function BackupRoutes() {
    const {path} = useRouteMatch();

    return (
        <Switch>
            <Route exact path={path}>
                <BackupsPage/>
            </Route>
        </Switch>
    );
}
