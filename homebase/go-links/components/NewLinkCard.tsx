import React from "react";
import {Form, Formik, FormikHelpers} from "formik";
import {CreateLinkParams} from "com_github_mjm_pi_tools/homebase/go-links/lib/mutate";
import {Alert} from "com_github_mjm_pi_tools/homebase/components/Alert";
import {ShortURLField} from "com_github_mjm_pi_tools/homebase/go-links/components/ShortURLField";
import {DestinationURLField} from "com_github_mjm_pi_tools/homebase/go-links/components/DestinationURLField";
import {DescriptionField} from "com_github_mjm_pi_tools/homebase/go-links/components/DescriptionField";
import {useCreateLink} from "com_github_mjm_pi_tools/homebase/go-links/lib/CreateLink";

export function NewLinkCard() {
    const [commit] = useCreateLink();

    async function onSubmit(values: CreateLinkParams, actions: FormikHelpers<CreateLinkParams>) {
        actions.setStatus(null);
        try {
            await commit(values);
            actions.resetForm();
        } catch (err) {
            actions.setStatus({error: err});
        }
    }

    return (
        <div className="bg-white overflow-hidden shadow rounded-lg">
            <Formik
                initialValues={{
                    shortURL: "",
                    destinationURL: "",
                    description: "",
                }}
                onSubmit={onSubmit}
            >{({status, isSubmitting}) => (
                <Form>
                    <div className="px-4 py-5 sm:px-6">
                        <h3 className="text-lg leading-6 font-medium text-gray-900">
                            Add a new link
                        </h3>
                    </div>
                    {status && status.error && (
                        <Alert title="Couldn't create new link" severity="error" rounded={false}>
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
                Create
              </button>
            </span>
                    </div>
                </Form>
            )}
            </Formik>
        </div>
    );
}
