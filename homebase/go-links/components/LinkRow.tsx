import React from "react";
import {Link} from "com_github_mjm_pi_tools/go-links/proto/links/links_pb";
import {destinationURL} from "com_github_mjm_pi_tools/homebase/go-links/lib/links_client";

export function LinkRow({link}: { link: Link }) {
    return (
        <li className="border-t border-gray-200">
            <a href={destinationURL(link.getShortUrl())}
               className="block hover:bg-gray-50 focus:outline-none focus:bg-gray-50 transition duration-150 ease-in-out">
                <div className="px-4 py-4 sm:px-6">
                    <div className="flex items-center justify-between">
                        <div className="text-sm leading-5 font-medium text-indigo-600 truncate">
                            go/{link.getShortUrl()}
                        </div>
                    </div>
                    <div className="mt-2">
                        <div className="flex items-center text-sm leading-5 text-gray-500">
                            {link.getDescription()}
                        </div>
                    </div>
                </div>
            </a>
        </li>
    );
}
