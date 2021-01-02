import React from "react";
import {Form, Formik, FormikHelpers} from "formik";
import {updateLink, UpdateLinkParams} from "com_github_mjm_pi_tools/homebase/go-links/lib/mutate";
import {useHistory} from "react-router-dom";
import {Alert} from "com_github_mjm_pi_tools/homebase/components/Alert";
import {ShortURLField} from "com_github_mjm_pi_tools/homebase/go-links/components/ShortURLField";
import {DestinationURLField} from "com_github_mjm_pi_tools/homebase/go-links/components/DestinationURLField";
import {DescriptionField} from "com_github_mjm_pi_tools/homebase/go-links/components/DescriptionField";
import {graphql, useFragment} from "react-relay/hooks";
import {EditLinkForm_link$key} from "com_github_mjm_pi_tools/homebase/api/__generated__/EditLinkForm_link.graphql";

export function EditLinkForm({link}: { link: EditLinkForm_link$key }) {
    const data = useFragment(
        graphql`
            fragment EditLinkForm_link on Link {
                id
                rawID
                shortURL
                destinationURL
                description
            }
        `,
        link,
    );
    const history = useHistory();

    async function onSubmit(values: UpdateLinkParams, actions: FormikHelpers<UpdateLinkParams>) {
        actions.setStatus(null);
        try {
            await updateLink(values);
            history.push("/go");
        } catch (err) {
            actions.setStatus({error: err});
        }
    }

    return (
        <Formik
            initialValues={{
                id: data.rawID,
                shortURL: data.shortURL,
                destinationURL: data.destinationURL,
                description: data.description,
            }}
            onSubmit={onSubmit}
        >{({status, isSubmitting}) => (
            <Form>
                {status && status.error && (
                    <Alert title="Couldn't save link changes" severity="error" rounded={false}>
                        {status.error.toString()}
                    </Alert>
                )}
                <div className="bg-gray-50 px-4 py-5 sm:p-6">
                    <div className="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
                        <ShortURLField/>
                        <DestinationURLField/>
                        <DescriptionField/>
                    </div>
                </div>
                <div className="px-4 py-5 sm:px-6 text-right">
            <span className="inline-flex rounded-md shadow-sm">
              <button type="submit"
                      disabled={isSubmitting}
                      className="inline-flex justify-center py-2 px-4 border border-transparent text-sm leading-5 font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-500 focus:outline-none focus:border-indigo-700 focus:ring-indigo active:bg-indigo-700 transition duration-150 ease-in-out">
                Save
              </button>
            </span>
                </div>
            </Form>
        )}
        </Formik>
    );
}
