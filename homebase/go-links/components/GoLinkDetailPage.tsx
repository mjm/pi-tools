import React from "react";
import {useParams} from "react-router-dom";
import useSWR from "swr";
import {Link} from "com_github_mjm_pi_tools/go-links/proto/links/links_pb";
import {GET_LINK} from "com_github_mjm_pi_tools/homebase/go-links/lib/fetch";
import {Helmet} from "react-helmet";
import {EditLinkForm} from "com_github_mjm_pi_tools/homebase/go-links/components/EditLinkForm";

export function GoLinkDetailPage() {
    const {id} = useParams<{ id: string }>();
    const {data, error} = useSWR<Link>([GET_LINK, id]);

    if (error) {
        console.error(error);
    }

    return (
        <main className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
            <Helmet>
                <title>{`Link Details${data ? `: go/${data.getShortUrl()}` : ""}`}</title>
            </Helmet>

            <div className="bg-white shadow overflow-hidden sm:rounded-lg">
                <div className="bg-white px-4 py-5 border-b border-gray-200 sm:px-6">
                    <div className="-ml-4 -mt-2 flex items-center justify-between flex-wrap sm:flex-no-wrap">
                        <div className="ml-4 mt-2">
                            <h3 className="text-lg leading-6 font-medium text-gray-900">
                                go/{data ? data.getShortUrl() : "…"}
                            </h3>
                        </div>
                        <div className="ml-4 mt-2 flex-shrink-0 flex">
                        </div>
                    </div>
                </div>
                {error && (
                    <div>{error.toString()}</div>
                )}
                {data && (
                    <EditLinkForm link={data}/>
                )}
            </div>
        </main>
    );
}
