import {createBrowserHistory} from "history";
import {matchRoutes} from "react-router-config";
import React, {
    useContext,
    useEffect,
    Suspense,
    useState,
} from "react";
import {ErrorBoundary} from "com_github_mjm_pi_tools/homebase/components/ErrorBoundary";

export const RoutingContext = React.createContext(null);

/**
 * A component that accesses the current route entry from RoutingContext and renders
 * that entry.
 */
export function RouterRenderer() {
    // Access the router
    const router = useContext(RoutingContext);
    // Improve the route transition UX by delaying transitions: show the previous route entry
    // for a brief period while the next route is being prepared. See
    // https://reactjs.org/docs/concurrent-mode-patterns.html#transitions
    // @ts-ignore
    const [startTransition, isPending] = React.unstable_useTransition();

    // Store the active entry in state - this allows the renderer to use features like
    // useTransition to delay when state changes become visible to the user.
    const [routeEntry, setRouteEntry] = useState(router.get());

    // On mount subscribe for route changes
    useEffect(() => {
        // Check if the route has changed between the last render and commit:
        const currentEntry = router.get();
        if (currentEntry !== routeEntry) {
            // if there was a concurrent modification, rerender and exit
            setRouteEntry(currentEntry);
            return;
        }

        // If there *wasn't* a concurrent change to the route, then the UI
        // is current: subscribe for subsequent route updates
        const dispose = router.subscribe(nextEntry => {
            // startTransition() delays the effect of the setRouteEntry (setState) call
            // for a brief period, continuing to show the old state while the new
            // state (route) is prepared.
            startTransition(() => {
                setRouteEntry(nextEntry);
            });
        });
        return () => dispose();

        // Note: this hook updates routeEntry manually; we exclude that variable
        // from the hook deps to avoid recomputing the effect after each change
        // triggered by the effect itself.
        // eslint-disable-next-line
    }, [router, startTransition]);

    // The current route value is an array of matching entries - one entry per
    // level of routes (to allow nested routes). We have to map each one to a
    // RouteComponent to allow suspending, and also pass its children correctly.
    // Conceptually, we want this structure:
    // ```
    // <RouteComponent
    //   component={entry[0].component}
    //   prepared={entrry[0].prepared}>
    //   <RouteComponent
    //     component={entry[1].component}
    //     prepared={entry[1].prepared}>
    //       // continue for nested items...
    //   </RouteComponent>
    // </RouteComponent>
    // ```
    // To achieve this, we reverse the list so we can start at the bottom-most
    // component, and iteratively construct parent components w the previous
    // value as the child of the next one:
    const reversedItems = [].concat(routeEntry.entries).reverse(); // reverse is in place, but we want a copy so concat
    const firstItem = reversedItems[0];
    // the bottom-most component is special since it will have no children
    // (though we could probably just pass null children to it)
    let routeComponent = (
        <RouteComponent
            component={firstItem.component}
            prepared={firstItem.prepared}
            routeData={firstItem.routeData}
        />
    );
    for (let ii = 1; ii < reversedItems.length; ii++) {
        const nextItem = reversedItems[ii];
        routeComponent = (
            <RouteComponent
                component={nextItem.component}
                prepared={nextItem.prepared}
                routeData={nextItem.routeData}
            >
                {routeComponent}
            </RouteComponent>
        );
    }

    // Routes can error so wrap in an <ErrorBoundary>
    // Routes can suspend, so wrap in <Suspense>
    return (
        <ErrorBoundary>
            <Suspense fallback={"Loading fallback..."}>
                {/* Indicate to the user that a transition is pending, even while showing the previous UI */}
                {isPending ? (
                    <div className="RouteRenderer-pending">Loading pending...</div>
                ) : null}
                {routeComponent}
            </Suspense>
        </ErrorBoundary>
    );
}

/**
 * The `component` property from the route entry is a Resource, which may or may not be ready.
 * We use a helper child component to unwrap the resource with component.read(), and then
 * render it if its ready.
 *
 * NOTE: calling routeEntry.route.component.read() directly in RouteRenderer woldn't work the
 * way we'd expect. Because that method could throw - either suspending or on error - the error
 * would bubble up to the *caller* of RouteRenderer. We want the suspend/error to bubble up to
 * our ErrorBoundary/Suspense components, so we have to ensure that the suspend/error happens
 * in a child component.
 */
function RouteComponent(props) {
    const Component = props.component;
    // const Component = props.component.read();
    const {routeData, prepared} = props;
    return (
        <Component
            routeData={routeData}
            prepared={prepared}
            children={props.children}
        />
    );
}

/**
 * A custom router built from the same primitives as react-router. Each object in `routes`
 * contains both a Component and a prepare() function that can preload data for the component.
 * The router watches for changes to the current location via the `history` package, maps the
 * location to the corresponding route entry, and then preloads the code and data for the route.
 */
export function createRouter(routes, options?) {
    // Initialize history
    const history = createBrowserHistory(options);

    // Find the initial match and prepare it
    const initialMatches = matchRoute(routes, history.location);
    const initialEntries = prepareMatches(initialMatches);
    let currentEntry = {
        location: history.location,
        entries: initialEntries,
    };

    // maintain a set of subscribers to the active entry
    let nextId = 0;
    const subscribers = new Map();

    // Listen for location changes, match to the route entry, prepare the entry,
    // and notify subscribers. Note that this pattern ensures that data-loading
    // occurs *outside* of - and *before* - rendering.
    const cleanup = history.listen((location, action) => {
        if (location.pathname === currentEntry.location.pathname) {
            return;
        }
        const matches = matchRoute(routes, location);
        const entries = prepareMatches(matches);
        const nextEntry = {
            location,
            entries,
        };
        currentEntry = nextEntry;
        subscribers.forEach(cb => cb(nextEntry));
    });

    // The actual object that will be passed on the RoutingConext.
    const context = {
        history,
        get() {
            return currentEntry;
        },
        preloadCode(pathname) {
            // preload just the code for a route, without storing the result
            // const matches = matchRoutes(routes, pathname);
            // matches.forEach(({route}) => route.component.load());
        },
        preload(pathname) {
            // preload the code and data for a route, without storing the result
            const matches = matchRoutes(routes, pathname);
            prepareMatches(matches);
        },
        subscribe(cb) {
            const id = nextId++;
            const dispose = () => {
                subscribers.delete(id);
            };
            subscribers.set(id, cb);
            return dispose;
        },
    };

    // Return both the context object and a cleanup function
    return {cleanup, context};
}

/**
 * Match the current location to the corresponding route entry.
 */
function matchRoute(routes, location) {
    const matchedRoutes = matchRoutes(routes, location.pathname);
    if (!Array.isArray(matchedRoutes) || matchedRoutes.length === 0) {
        throw new Error("No route for " + location.pathname);
    }
    return matchedRoutes;
}

/**
 * Load the data for the matched route, given the params extracted from the route
 */
function prepareMatches(matches) {
    return matches.map(match => {
        const {route, match: matchData} = match;
        let prepared;
        if (route.prepare) {
            prepared = route.prepare(matchData.params);
        }
        // const Component = route.component.get();
        // if (Component == null) {
        //     route.component.load(); // eagerly load
        // }
        return {component: route.component, prepared, routeData: matchData};
    });
}
