import React from "react";
import {LinkRow} from "com_github_mjm_pi_tools/homebase/go-links/components/LinkRow";
import {graphql, useFragment} from "react-relay/hooks";
import {RecentLinksList_viewer$key} from "com_github_mjm_pi_tools/homebase/api/__generated__/RecentLinksList_viewer.graphql";

export function RecentLinksList({viewer}: { viewer: RecentLinksList_viewer$key }) {
    const data = useFragment(
        graphql`
            fragment RecentLinksList_viewer on Viewer {
                links(first: 30) @connection(key: "RecentLinksList_links") {
                    edges {
                        node {
                            id
                            ...LinkRow_link
                        }
                    }
                }
            }
        `,
        viewer,
    );

    const linkNodes = data.links.edges.map(e => e.node);

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
