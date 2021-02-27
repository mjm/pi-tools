import React from "react";
import {Helmet} from "react-helmet";
import {PageHeader} from "com_github_mjm_pi_tools/homebase/components/PageHeader";
import {ErrorBoundary} from "com_github_mjm_pi_tools/homebase/components/ErrorBoundary";
import {graphql, useLazyLoadQuery} from "react-relay/hooks";
import {BackupsPageQuery} from "com_github_mjm_pi_tools/homebase/api/__generated__/BackupsPageQuery.graphql";
import {BackupsList} from "com_github_mjm_pi_tools/homebase/backups/components/BackupsList";

export function BackupsPage() {
    return (
        <main className="mb-8">
            <Helmet>
                <title>Backups</title>
            </Helmet>

            <PageHeader>
                Recent archives
            </PageHeader>

            <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                <React.Suspense fallback="Loading…">
                    <ErrorBoundary>
                        <BackupsPageInner/>
                    </ErrorBoundary>
                </React.Suspense>
            </div>
        </main>
    );
}

function BackupsPageInner() {
    const data = useLazyLoadQuery<BackupsPageQuery>(
        graphql`
            query BackupsPageQuery {
                viewer {
                    ...BackupsList_viewer
                }
            }
        `,
        {},
    );

    if (!data.viewer) {
        return null;
    }

    return <BackupsList viewer={data.viewer}/>;
}
