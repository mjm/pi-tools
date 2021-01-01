import React from "react";
import {BrowserRouter as Router, Route, Switch, useLocation} from "react-router-dom";
import {Helmet} from "react-helmet";
import {NavigationBar} from "com_github_mjm_pi_tools/homebase/components/NavigationBar";
import {TripRoutes} from "com_github_mjm_pi_tools/homebase/trips/components/TripRoutes";
import {GoLinkRoutes} from "com_github_mjm_pi_tools/homebase/go-links/components/GoLinkRoutes";
import {HomePage} from "com_github_mjm_pi_tools/homebase/homepage/components/HomePage";
import {RelayEnvironmentProvider} from "react-relay/hooks";
import RelayEnvironment from "com_github_mjm_pi_tools/homebase/lib/environment";

export function App() {
    return (
        <RelayEnvironmentProvider environment={RelayEnvironment}>
            <React.Suspense fallback={"Loading..."}>
                <Router>
                    <Helmet>
                        <title>Homebase</title>
                        <meta name="viewport" content="width=device-width, initial-scale=1"/>
                        <link rel="stylesheet" href="https://rsms.me/inter/inter.css"/>
                    </Helmet>
                    <div>
                        <NavigationBar/>

                        <Switch>
                            <Route exact path="/">
                                <HomePage/>
                            </Route>
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
            </React.Suspense>
        </RelayEnvironmentProvider>
    );
}

function NoMatch() {
    const location = useLocation();
    console.log(location);

    return (
        <div>Not Found</div>
    );
}
