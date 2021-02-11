import React from "react";
import {graphql, useFragment} from "react-relay/hooks";
import {ArchiveRow_archive$key} from "com_github_mjm_pi_tools/homebase/api/__generated__/ArchiveRow_archive.graphql";
import {format, parseISO} from "date-fns";
import {TransitionLink} from "com_github_mjm_pi_tools/homebase/components/TransitionLink";

export function ArchiveRow({archive}: { archive: ArchiveRow_archive$key }) {
    const data = useFragment(
        graphql`
            fragment ArchiveRow_archive on Archive {
                id
                name
                createdAt
            }
        `,
        archive,
    );

    const createdAt = parseISO(data.createdAt);

    return (
        <tr>
            <td className="px-6 py-4 whitespace-nowrap text-sm leading-5 font-medium text-gray-900">
                {data.name}
            </td>
            <td className="px-6 py-4 whitespace-nowrap text-sm leading-5 text-gray-500">
                {format(createdAt, "PPpp")}
            </td>
            <td className="px-6 py-4 whitespace-nowrap text-right text-sm leading-5 font-medium">
                <TransitionLink to={`/backups/${data.id}`}
                                className="text-indigo-600 hover:text-indigo-900">
                    Details
                </TransitionLink>
            </td>
        </tr>
    )
}
