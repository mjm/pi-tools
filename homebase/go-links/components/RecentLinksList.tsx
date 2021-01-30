import React from "react";
import {LinkRow} from "com_github_mjm_pi_tools/homebase/go-links/components/LinkRow";
import {graphql, useFragment} from "react-relay/hooks";
import {RecentLinksList_links$key} from "com_github_mjm_pi_tools/homebase/api/__generated__/RecentLinksList_links.graphql";

export function RecentLinksList({links}: { links: RecentLinksList_links$key }) {
    const data = useFragment(
        graphql`
            fragment RecentLinksList_links on LinkConnection {
                edges {
                    node {
                        id
                        ...LinkRow_link
                    }
                }
            }
        `,
        links,
    );

    const linkNodes = data.edges.map(e => e.node);

    return (
        <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="px-4 py-5 sm:px-6">
                <h3 className="text-lg leading-6 font-medium text-gray-900">
                    Recently added links
                </h3>
            </div>
            <ul className="border-b border-gray-200">
                {linkNodes.map(link => (
                    <LinkRow key={link.id} link={link}/>
                ))}
            </ul>
        </div>
    );
}
