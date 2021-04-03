import React from "react";
import {Form, Formik, FormikHelpers} from "formik";
import {graphql, useFragment} from "react-relay/hooks";
import {EditLinkForm_link$key} from "../../__generated__/EditLinkForm_link.graphql";
import {useUpdateLink} from "../../mutations/UpdateLink";
import {useRouter} from "next/router";
import {UpdateLinkInput} from "../../__generated__/UpdateLinkMutation.graphql";
import Alert from "../Alert";
import ShortURLField from "./ShortURLField";
import DestinationURLField from "./DestinationURLField";
import DescriptionField from "./DescriptionField";

export default function EditLinkForm({link}: { link: EditLinkForm_link$key }) {
    const data = useFragment(
        graphql`
            fragment EditLinkForm_link on Link {
                id
                shortURL
                destinationURL
                description
            }
        `,
        link,
    );
    const [commit] = useUpdateLink();
    const router = useRouter();

    async function onSubmit(values: UpdateLinkInput, actions: FormikHelpers<UpdateLinkInput>) {
        actions.setStatus(null);
        try {
            await commit(values);
            await router.push("/go");
        } catch (err) {
            actions.setStatus({error: err});
        }
    }

    return (
        <Formik
            initialValues={{
                id: data.id,
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
