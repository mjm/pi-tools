import React from "react";
import {Helmet} from "react-helmet";
import {NavigationBar} from "com_github_mjm_pi_tools/homebase/components/NavigationBar";
import {RelayEnvironmentProvider} from "react-relay/hooks";
import RelayEnvironment from "com_github_mjm_pi_tools/homebase/lib/environment";
import {createRouter, RouterRenderer, RoutingContext} from "com_github_mjm_pi_tools/homebase/components/Router";
import routes from "com_github_mjm_pi_tools/homebase/routes";

const router = createRouter(routes);

export function App() {
    return (
        <RelayEnvironmentProvider environment={RelayEnvironment}>
            <React.Suspense fallback={"Loading..."}>
                <RoutingContext.Provider value={router.context}>
                    <Helmet>
                        <title>Homebase</title>
                        <meta name="viewport" content="width=device-width, initial-scale=1"/>
                        <link rel="stylesheet" href="https://rsms.me/inter/inter.css"/>
                    </Helmet>
                    <div>
                        <NavigationBar/>

                        <React.Suspense fallback={"Loading..."}>
                            <RouterRenderer/>
                        </React.Suspense>
                    </div>
                </RoutingContext.Provider>
            </React.Suspense>
        </RelayEnvironmentProvider>
    );
}

// function NoMatch() {
//     const location = useLocation();
//     console.log(location);
//
//     return (
//         <div>Not Found</div>
//     );
// }
