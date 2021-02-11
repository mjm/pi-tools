import React from "react";
import {graphql, useFragment, usePaginationFragment} from "react-relay/hooks";
import {BackupsList_viewer$key} from "com_github_mjm_pi_tools/homebase/api/__generated__/BackupsList_viewer.graphql";
import {ArchiveRow} from "com_github_mjm_pi_tools/homebase/backups/components/ArchiveRow";

export function BackupsList({viewer}: { viewer: BackupsList_viewer$key }) {
    const data = useFragment(
        graphql`
            fragment BackupsList_viewer on Viewer
            @argumentDefinitions(
                count: { type: "Int", defaultValue: 10 }
                cursor: { type: "Cursor" }
            ) {
                backupArchives(first: $count, after: $cursor)
                @connection(key: "BackupsList_backupArchives", filters: ["kind"]) {
                    edges {
                        node {
                            id
                            ...ArchiveRow_archive
                        }
                    }
                }
            }
        `,
        viewer,
    );

    if (!data) {
        return null;
    }

    const archiveNodes = data.backupArchives.edges.map(e => e.node);

    return (
        <div className="flex flex-col">
            <div className="-my-2 overflow-x-auto sm:-mx-6 lg:-mx-8">
                <div className="py-2 align-middle inline-block min-w-full sm:px-6 lg:px-8">
                    <div className="shadow overflow-hidden border-b border-gray-200 sm:rounded-lg">
                        <table className="min-w-full divide-y divide-gray-200">
                            <thead>
                            <tr>
                                <th className="px-6 py-3 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider">
                                    Name
                                </th>
                                <th className="px-6 py-3 bg-gray-50 text-left text-xs leading-4 font-medium text-gray-500 uppercase tracking-wider">
                                    Created at
                                </th>
                                <th className="px-6 py-3 bg-gray-50"/>
                            </tr>
                            </thead>
                            <tbody className="bg-white divide-y divide-gray-200">
                            {archiveNodes.map(archive => (
                                <ArchiveRow archive={archive}/>
                            ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    );
}
