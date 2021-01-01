import React from "react";
import {Link} from "react-router-dom";

export function HomePageCard({title, icon, children, footerHref, footer}: {
    title: React.ReactNode;
    icon: React.ReactNode;
    children: React.ReactNode;
    footerHref: string;
    footer: React.ReactNode;
}) {
    return (
        <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-5">
                <div className="flex items-center">
                    <div className="flex-shrink-0">
                        {icon}
                    </div>
                    <div className="ml-5 w-0 flex-1">
                        <dl>
                            <dt className="text-sm leading-5 font-medium text-gray-500 truncate">
                                {title}
                            </dt>
                            {children != null ? (
                                <dd>
                                    <div className="text-lg leading-7 font-medium text-gray-900 truncate">
                                        {children}
                                    </div>
                                </dd>
                            ) : null}
                        </dl>
                    </div>
                </div>
            </div>
            <div className="bg-gray-50 px-5 py-3">
                <div className="text-sm leading-5">
                    {footerHref.startsWith("http") ? (
                        <a href={footerHref}
                           target="_blank"
                           className="font-medium text-indigo-700 hover:text-indigo-900 transition ease-in-out duration-150">
                            {footer}
                        </a>
                    ) : (
                        <Link to={footerHref}
                              className="font-medium text-indigo-700 hover:text-indigo-900 transition ease-in-out duration-150">
                            {footer}
                        </Link>
                    )}
                </div>
            </div>
        </div>
    );
}
