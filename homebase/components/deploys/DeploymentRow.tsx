import {graphql, useFragment} from "react-relay";
import {DeploymentRow_deploy$key} from "../../__generated__/DeploymentRow_deploy.graphql";
import {formatDistanceToNow, parseISO} from "date-fns";
import Link from "next/link";

export default function DeploymentRow({deploy, isLast}: {
    deploy: DeploymentRow_deploy$key;
    isLast: boolean;
}) {
    const data = useFragment(
        graphql`
            fragment DeploymentRow_deploy on Deploy {
                rawID
                state
                commitSHA
                commitMessage
                startedAt
            }
        `,
        deploy,
    );

    let iconStyle = "";
    let iconContent = null;

    switch (data.state) {
        case "SUCCESS":
            iconStyle = "bg-green-500";
            iconContent =
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd"
                          d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                          clipRule="evenodd"/>
                </svg>;
            break;
        case "FAILURE":
            iconStyle = "bg-red-600";
            iconContent =
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24"
                     stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12"/>
                </svg>;
            break;
        case "INACTIVE":
            iconStyle = "bg-gray-300";
            iconContent =
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path fillRule="evenodd"
                          d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
                          clipRule="evenodd"/>
                </svg>;
            break;
        case "IN_PROGRESS":
            iconStyle = "bg-yellow-500";
            iconContent =
                <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                    <path
                        d="M6 10a2 2 0 11-4 0 2 2 0 014 0zM12 10a2 2 0 11-4 0 2 2 0 014 0zM16 12a2 2 0 100-4 2 2 0 000 4z"/>
                </svg>;
            break;
    }

    return (
        <li>
            <div className="relative pb-8">
                {!isLast && <span className="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200"
                                  aria-hidden="true"/>}
                <div className="relative flex space-x-3">
                    <div>
            <span
                className={`h-8 w-8 rounded-full text-white flex items-center justify-center ring-8 ring-white ${iconStyle}`}>
                {iconContent}
            </span>
                    </div>
                    <div className="min-w-0 flex-1 pt-1.5 flex justify-between space-x-4">
                        <div>
                            <p className={`text-sm`}>
                                <Link href={`/deploys/${data.rawID}`}>
                                    <a className="font-medium text-gray-700 hover:text-indigo-900 transition ease-in-out duration-150">
                                        {data.commitMessage}
                                    </a>
                                </Link>
                            </p>
                        </div>
                        <div className="text-right text-sm whitespace-nowrap text-gray-500">
                            <time dateTime={data.startedAt}>
                                {formatDistanceToNow(parseISO(data.startedAt), {addSuffix: true})}
                            </time>
                        </div>
                    </div>
                </div>
            </div>
        </li>
    );
}
