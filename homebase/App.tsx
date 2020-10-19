import React from "react";
import {BrowserRouter as Router, Switch} from "react-router-dom";
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
                    <TripRoutes/>
                </Switch>
            </div>
        </Router>
    );
}
