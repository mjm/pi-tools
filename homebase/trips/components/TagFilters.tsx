import React from "react";
import {TripTag} from "com_github_mjm_pi_tools/homebase/trips/components/TripTag";
import {graphql, useFragment} from "react-relay/hooks";
import {TagFilters_tags$key} from "com_github_mjm_pi_tools/homebase/api/__generated__/TagFilters_tags.graphql";

export function TagFilters({tags}: { tags: TagFilters_tags$key }) {
    const data = useFragment(
        graphql`
            fragment TagFilters_tags on Viewer {
                tags(first: 5) {
                    edges {
                        node {
                            name
                            tripCount
                        }
                    }
                }
            }
        `,
        tags,
    );

    if (!data.tags) {
        return null;
    }

    const tagNodes = data.tags.edges.map(e => e.node);
    return (
        <div className="flex flex-row items-baseline space-x-3 overflow-hidden whitespace-nowrap flex-wrap">
            <span className="font-medium uppercase tracking-wider">Popular Tags:</span>
            {tagNodes.map(tag => (
                <TripTag key={tag.name}>
                    {tag.name} ({tag.tripCount})
                </TripTag>
            ))}
        </div>
    );
}
