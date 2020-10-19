import React from "react";
import {Route} from "react-router-dom";
import {SWRConfig} from "swr";
import {TripsPage} from "com_github_mjm_pi_tools/homebase/trips/components/TripsPage";
import {TripPage} from "com_github_mjm_pi_tools/homebase/trips/components/TripPage";
import {fetcher} from "com_github_mjm_pi_tools/homebase/trips/lib/fetch";

export function TripRoutes() {
    return (
        <SWRConfig value={{
            fetcher,
        }}>
            <Route exact path="/trips">
                <TripsPage/>
            </Route>
            <Route path="/trips/:id">
                <TripPage/>
            </Route>
        </SWRConfig>
    );
}
