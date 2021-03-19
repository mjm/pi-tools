import React from "react";
import {RoutingContext} from "com_github_mjm_pi_tools/homebase/components/Router";

// export function TransitionLink({to, className, children}: {
//     to: string;
//     className?: string;
//     children: React.ReactNode;
// }) {
//     const history = useHistory();
//     // @ts-ignore
//     const [startTransition, isPending] = React.unstable_useTransition();
//
//     function onClick(e) {
//         e.preventDefault();
//         startTransition(() => {
//             history.push(to);
//         });
//     }
//
//     return (
//         <a className={className} href={to} onClick={onClick}>
//             {children}
//         </a>
//     );
// }

// import React from 'react';

const { useCallback, useContext } = React;

/**
 * An alternative to react-router's Link component that works with
 * our custom RoutingContext.
 */
export function TransitionLink(props) {
    const router: any = useContext(RoutingContext);

    // When the user clicks, change route
    const changeRoute = useCallback(
        event => {
            event.preventDefault();
            router.history.push(props.to);
        },
        [props.to, router],
    );

    // Callback to preload just the code for the route:
    // we pass this to onMouseEnter, which is a weaker signal
    // that the user *may* navigate to the route.
    const preloadRouteCode = useCallback(() => {
        router.preloadCode(props.to);
    }, [props.to, router]);

    // Callback to preload the code and data for the route:
    // we pass this to onMouseDown, since this is a stronger
    // signal that the user will likely complete the navigation
    const preloadRoute = useCallback(() => {
        router.preload(props.to);
    }, [props.to, router]);

    return (
        <a
            className={props.className}
            href={props.to}
            onClick={changeRoute}
            onMouseEnter={preloadRouteCode}
            onMouseDown={preloadRoute}
        >
            {props.children}
        </a>
    );
}
