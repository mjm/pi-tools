import React from "react";
import {BrowserRouter as Router, Redirect, Route, Switch, useLocation} from "react-router-dom";
import {Helmet} from "react-helmet";
import {NavigationBar} from "com_github_mjm_pi_tools/homebase/components/NavigationBar";
import {TripRoutes} from "com_github_mjm_pi_tools/homebase/trips/components/TripRoutes";

export function App() {
    return (
        <Router>
            <Helmet>
                <title>Homebase</title>
            </Helmet>
            <div>
                <NavigationBar/>

                <Switch>
                    <Redirect exact from="/" to="/trips"/>
                    <TripRoutes/>
                    <Route path="*">
                        <NoMatch />
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
