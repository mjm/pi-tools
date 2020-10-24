import React from "react";
import {BrowserRouter as Router, Redirect, Route, Switch, useLocation} from "react-router-dom";
import {Helmet} from "react-helmet";
import {NavigationBar} from "com_github_mjm_pi_tools/homebase/components/NavigationBar";
import {TripRoutes} from "com_github_mjm_pi_tools/homebase/trips/components/TripRoutes";
import {GoLinkRoutes} from "com_github_mjm_pi_tools/homebase/go-links/components/GoLinkRoutes";

export function App() {
    return (
        <Router>
            <Helmet>
                <title>Homebase</title>
                <meta name="viewport" content="width=device-width, initial-scale=1"/>
            </Helmet>
            <div>
                <NavigationBar/>

                <Switch>
                    <Redirect exact from="/" to="/trips"/>
                    <Route path="/trips">
                        <TripRoutes/>
                    </Route>
                    <Route path="/go">
                        <GoLinkRoutes/>
                    </Route>
                    <Route path="*">
                        <NoMatch/>
                    </Route>
                </Switch>
            </div>
        </Router>
    );
}

function NoMatch() {
    const location = useLocation();
    console.log(location);

    return (
        <div>Not Found</div>
    );
}
