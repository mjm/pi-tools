import React from "react";
import {PageHeader} from "com_github_mjm_pi_tools/homebase/components/PageHeader";
import {MostRecentTripCard} from "com_github_mjm_pi_tools/homebase/homepage/components/MostRecentTripCard";
import {FiringAlertsCard} from "com_github_mjm_pi_tools/homebase/homepage/components/FiringAlertsCard";

export function HomePage() {
    return (
        <main className="mb-8">
            <PageHeader>Homebase</PageHeader>
            <div className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
                <div className="mt-2 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-3">
                    <MostRecentTripCard />
                    <FiringAlertsCard />
                </div>
            </div>
        </main>
    );
}
