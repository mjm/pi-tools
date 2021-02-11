import React from "react";
import {useRouteMatch} from "react-router-dom";
import {TransitionLink} from "com_github_mjm_pi_tools/homebase/components/TransitionLink";

export function NavigationBar() {
    const [showMenu, setShowMenu] = React.useState(false);

    return (
        <nav className="bg-gray-800">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    <div className="flex items-center">
                        <div className="flex-shrink-0">
                            <TransitionLink to="/">
                                <svg className="h-6 w-6 text-white" xmlns="http://www.w3.org/2000/svg" fill="none"
                                     viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                          d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"/>
                                </svg>
                            </TransitionLink>
                        </div>
                        <div className="hidden md:block">
                            <div className="ml-10 flex items-baseline space-x-4">
                                <NavLink to="/trips">
                                    Your Trips
                                </NavLink>
                                <NavLink to="/go">
                                    Go Links
                                </NavLink>
                                <NavLink to="/backups">
                                    Backups
                                </NavLink>
                            </div>
                        </div>
                    </div>
                    <div className="-mr-2 flex md:hidden">
                        <button
                            onClick={() => setShowMenu(shown => !shown)}
                            className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 focus:outline-none focus:bg-gray-700 focus:text-white">
                            <svg className={`${showMenu ? "hidden" : "block"} h-6 w-6`} stroke="currentColor"
                                 fill="none" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                      d="M4 6h16M4 12h16M4 18h16"/>
                            </svg>
                            <svg className={`${showMenu ? "block" : "hidden"} h-6 w-6`} stroke="currentColor"
                                 fill="none" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                      d="M6 18L18 6M6 6l12 12"/>
                            </svg>
                        </button>
                    </div>
                </div>
            </div>

            <div className={`${showMenu ? "block" : "hidden"} md:hidden`}>
                <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3">
                    <MobileNavLink to="/trips">
                        Your Trips
                    </MobileNavLink>
                    <MobileNavLink to="/go">
                        Go Links
                    </MobileNavLink>
                    <MobileNavLink to="/backups">
                        Backups
                    </MobileNavLink>
                </div>
            </div>
        </nav>
    );
}

function NavLink({to, exact, children}: {
    to: string;
    exact?: boolean;
    children: React.ReactNode;
}) {
    const match = useRouteMatch({
        path: to,
        exact,
    });

    return (
        <TransitionLink to={to}
                        className={`px-3 py-2 rounded-md text-sm font-medium ${match ? "text-white bg-gray-900" : "text-gray-300 hover:text-white hover:bg-gray-700"} focus:outline-none focus:text-white focus:bg-gray-700`}>
            {children}
        </TransitionLink>
    );
}

function MobileNavLink({to, exact, children}: {
    to: string;
    exact?: boolean;
    children: React.ReactNode;
}) {
    const match = useRouteMatch({
        path: to,
        exact,
    });

    return (
        <TransitionLink to={to}
                        className={`block px-3 py-2 rounded-md text-base font-medium ${match ? "text-white bg-gray-900" : "text-gray-300 hover:text-white hover:bg-gray-700"} focus:outline-none focus:text-white focus:bg-gray-700`}>
            {children}
        </TransitionLink>
    );
}
