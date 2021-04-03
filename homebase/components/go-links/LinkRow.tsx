import React from "react";
import {graphql, useFragment} from "react-relay/hooks";
import {LinkRow_link$key} from "../../__generated__/LinkRow_link.graphql";
import Link from "next/link";
import {destinationURL} from "../../lib/go_links";

export default function LinkRow({link}: { link: LinkRow_link$key }) {
    const data = useFragment(
        graphql`
            fragment LinkRow_link on Link {
                id
                shortURL
                description
            }
        `,
        link,
    );

    return (
        <li className="border-t border-gray-200">
            <Link href={`/go/${data.id}`}>
                <a
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
                </a>
            </Link>
        </li>
    );
}
