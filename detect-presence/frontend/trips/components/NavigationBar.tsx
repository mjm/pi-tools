import React from "react";
import {Link, useRouteMatch} from "react-router-dom";

export function NavigationBar() {
    return (
        <nav className="bg-gray-800">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    <div className="flex items-center">
                        <div className="flex-shrink-0">
                            <img className="h-8 w-8" src="https://tailwindui.com/img/logos/workflow-mark-on-dark.svg"
                                 alt="Workflow logo"
                            />
                        </div>
                        <div className="hidden md:block">
                            <div className="ml-10 flex items-baseline space-x-4">
                                <NavLink to="/">
                                    Your Trips
                                </NavLink>
                            </div>
                        </div>
                    </div>
                    <div className="-mr-2 flex md:hidden">
                        <button
                            className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 focus:outline-none focus:bg-gray-700 focus:text-white">
                            <svg className="block h-6 w-6" stroke="currentColor" fill="none" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                      d="M4 6h16M4 12h16M4 18h16"/>
                            </svg>
                            <svg className="hidden h-6 w-6" stroke="currentColor" fill="none" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2"
                                      d="M6 18L18 6M6 6l12 12"/>
                            </svg>
                        </button>
                    </div>
                </div>
            </div>

            <div className="hidden md:hidden">
                <div className="px-2 pt-2 pb-3 space-y-1 sm:px-3">
                    <MobileNavLink to="/">
                        Your Trips
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
        <Link to={to}
              className={`px-3 py-2 rounded-md text-sm font-medium ${match ? "text-white bg-gray-900" : "text-gray-300 hover:text-white hover:bg-gray-700"} focus:outline-none focus:text-white focus:bg-gray-700`}>
            {children}
        </Link>
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
        <Link to={to}
              className={`block px-3 py-2 rounded-md text-base font-medium ${match ? "text-white bg-gray-900" : "text-gray-300 hover:text-white hover:bg-gray-700"} focus:outline-none focus:text-white focus:bg-gray-700`}>
            {children}
        </Link>
    );
}
