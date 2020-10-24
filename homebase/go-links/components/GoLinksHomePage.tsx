import React from "react";
import {PageHeader} from "com_github_mjm_pi_tools/homebase/components/PageHeader";
import {NewLinkCard} from "com_github_mjm_pi_tools/homebase/go-links/components/NewLinkCard";
import {RecentLinksList} from "com_github_mjm_pi_tools/homebase/go-links/components/RecentLinksList";

export function GoLinksHomePage() {
    return (
        <main className="mb-8">
            <PageHeader>
                Go links
            </PageHeader>
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="grid grid-cols-1 sm:grid-cols-2 mt-6 gap-8">
                    <NewLinkCard/>
                    <RecentLinksList/>
                </div>
            </div>
        </main>
    );
}
