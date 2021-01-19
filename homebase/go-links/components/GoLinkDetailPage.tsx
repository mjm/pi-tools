import React from "react";
import {useParams} from "react-router-dom";
import {Helmet} from "react-helmet";
import {EditLinkForm} from "com_github_mjm_pi_tools/homebase/go-links/components/EditLinkForm";
import {Alert} from "com_github_mjm_pi_tools/homebase/components/Alert";
import {graphql, useLazyLoadQuery} from "react-relay/hooks";
import {GoLinkDetailPageQuery} from "com_github_mjm_pi_tools/homebase/api/__generated__/GoLinkDetailPageQuery.graphql";

export function GoLinkDetailPage() {
    const {id} = useParams<{ id: string }>();
    const data = useLazyLoadQuery<GoLinkDetailPageQuery>(
        graphql`
            query GoLinkDetailPageQuery($id: ID!) {
                viewer {
                    link(id: $id) {
                        id
                        shortURL
                        ...EditLinkForm_link
                    }
                }
            }
        `,
        {id},
    );

    const link = data.viewer.link;

    return (
        <main className="max-w-3xl mx-auto py-6 sm:px-6 lg:px-8">
            <Helmet>
                <title>{`Link Details${link ? `: go/${link.shortURL}` : ""}`}</title>
            </Helmet>

            <div className="bg-white shadow overflow-hidden sm:rounded-lg">
                <div className="px-4 py-5 sm:px-6">
                    <div className="-ml-4 -mt-2 flex items-center justify-between flex-wrap sm:flex-nowrap">
                        <div className="ml-4 mt-2">
                            <h3 className="text-lg leading-6 font-medium text-gray-900">
                                go/{link ? link.shortURL : "…"}
                            </h3>
                        </div>
                    </div>
                </div>
                {link ? (
                    <EditLinkForm link={link}/>
                ) : (
                    <Alert title="Couldn't load link details" severity="error" rounded={false}>
                        No link was found with this ID.
                    </Alert>
                )}
            </div>
        </main>
    );
}
