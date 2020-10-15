import React from "react";
import {BrowserRouter as Router, NavLink, Route, Switch} from "react-router-dom";
import {Helmet} from "react-helmet";
import TripsPage from "com_github_mjm_pi_tools/detect-presence/frontend/trips/components/TripsPage";
import {TripPage} from "com_github_mjm_pi_tools/detect-presence/frontend/trips/components/TripPage";

export function App() {
    return (
        <Router>
            <Helmet>
                <title>Presence Dashboard</title>
            </Helmet>
            <div className="container mx-auto">
                <nav className="mb-4">
                    <ul className="text-indigo-700">
                        <li className="p-3">
                            <NavLink exact to="/" activeClassName="text-black">
                                Your Trips
                            </NavLink>
                        </li>
                    </ul>
                </nav>

                <Switch>
                    <Route path="/trips/:id">
                        <TripPage/>
                    </Route>
                    <Route path="/">
                        <TripsPage/>
                    </Route>
                </Switch>
            </div>
        </Router>
    );
}
