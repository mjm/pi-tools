import React from "react";
import {Link as RouterLink} from "react-router-dom";
import {destinationURL} from "com_github_mjm_pi_tools/homebase/go-links/lib/links_client";
import {graphql, useFragment} from "react-relay/hooks";
import {LinkRow_link$key} from "com_github_mjm_pi_tools/homebase/api/__generated__/LinkRow_link.graphql";

export function LinkRow({link}: { link: LinkRow_link$key }) {
    const data = useFragment(
        graphql`
            fragment LinkRow_link on Link {
                rawID
                shortURL
                description
            }
        `,
        link,
    );

    return (
        <li className="border-t border-gray-200">
            <RouterLink to={`/go/${data.rawID}`}
                        className="block hover:bg-gray-50 focus:outline-none focus:bg-gray-50 transition duration-150 ease-in-out">
                <div className="px-4 py-4 sm:px-6">
                    <div className="flex items-center justify-between">
                        <a href={destinationURL(data.shortURL)} target="_blank"
                           className="block text-sm leading-5 font-medium text-indigo-600 truncate"
                           onClick={e => {
                               e.stopPropagation();
                           }}>
                            go/{data.shortURL}
                        </a>
                    </div>
                    <div className="mt-2">
                        <div className="flex items-center text-sm leading-5 text-gray-500">
                            {data.description}
                        </div>
                    </div>
                </div>
            </RouterLink>
        </li>
    );
}
