import React from "react";
import {BrowserRouter as Router, Route, Switch} from "react-router-dom";
import {Helmet} from "react-helmet";
import TripsPage from "com_github_mjm_pi_tools/detect-presence/frontend/trips/components/TripsPage";
import {TripPage} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/components/TripPage";
import {NavigationBar} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/components/NavigationBar";
import {SWRConfig} from "swr";
import {fetcher} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/lib/fetch";

export function App() {
    return (
        <Router>
            <SWRConfig value={{
                fetcher,
            }}>
                <Helmet>
                    <title>Presence Dashboard</title>
                </Helmet>
                <div>
                    <NavigationBar/>

                    <Switch>
                        <Route path="/trips/:id">
                            <TripPage/>
                        </Route>
                        <Route path="/">
                            <TripsPage/>
                        </Route>
                    </Switch>
                </div>
            </SWRConfig>
        </Router>
    );
}
