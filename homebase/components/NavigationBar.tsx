import React from "react";
import Link from "next/link";
import {useRouter} from "next/router";
import {HomeIcon, MenuIcon, XIcon} from "@heroicons/react/outline";
import {CodeIcon, DatabaseIcon, LinkIcon, MapIcon} from "@heroicons/react/solid";

export default function NavigationBar() {
    const [showMenu, setShowMenu] = React.useState(false);

    return (
        <nav className="bg-gray-800">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex items-center justify-between h-16">
                    <div className="flex items-center">
                        <div className="flex-shrink-0">
                            <Link href="/">
                                <a>
                                    <HomeIcon className="h-6 w-6 text-white"/>
                                </a>
                            </Link>
                        </div>
                        <div className="hidden md:block">
                            <div className="ml-10 flex items-baseline space-x-4">
                                <NavLink to="/trips">
                                    <MapIcon className="h-4 w-4 mr-2"/>
                                    Your Trips
                                </NavLink>
                                <NavLink to="/go">
                                    <LinkIcon className="h-4 w-4 mr-2"/>
                                    Go Links
                                </NavLink>
                                <NavLink to="/backups">
                                    <DatabaseIcon className="h-4 w-4 mr-2"/>
                                    Backups
                                </NavLink>
                                <NavLink to="/deploys">
                                    <CodeIcon className="h-4 w-4 mr-2"/>
                                    Deploys
                                </NavLink>
                            </div>
                        </div>
                    </div>
                    <div className="-mr-2 flex md:hidden">
                        <button
                            onClick={() => setShowMenu(shown => !shown)}
                            className="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-white hover:bg-gray-700 focus:outline-none focus:bg-gray-700 focus:text-white">
                            <MenuIcon className={`${showMenu ? "hidden" : "block"} h-6 w-6`}/>
                            <XIcon className={`${showMenu ? "block" : "hidden"} h-6 w-6`}/>
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
                    <MobileNavLink to="/deploys">
                        Deploys
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
    const router = useRouter();
    const match = exact ? router.asPath === to : router.asPath.startsWith(to);

    return (
        <Link href={to}>
            <a className={`inline-flex items-center px-3 py-2 rounded-md text-sm font-medium ${match ? "text-white bg-gray-900" : "text-gray-300 hover:text-white hover:bg-gray-700"} focus:outline-none focus:text-white focus:bg-gray-700`}>
                {children}
            </a>
        </Link>
    );
}

function MobileNavLink({to, exact, children}: {
    to: string;
    exact?: boolean;
    children: React.ReactNode;
}) {
    const router = useRouter();
    const match = exact ? router.asPath === to : router.asPath.startsWith(to);

    return (
        <Link href={to}>
            <a
                className={`flex flex-col px-3 py-2 rounded-md text-base font-medium ${match ? "text-white bg-gray-900" : "text-gray-300 hover:text-white hover:bg-gray-700"} focus:outline-none focus:text-white focus:bg-gray-700`}>
                {children}
            </a>
        </Link>
    );
}
