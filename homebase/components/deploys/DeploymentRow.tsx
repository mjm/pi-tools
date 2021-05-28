import {graphql, useFragment} from "react-relay";
import {DeploymentRow_deploy$key} from "../../__generated__/DeploymentRow_deploy.graphql";
import {formatDistanceToNow, parseISO} from "date-fns";
import Link from "next/link";
import {CheckIcon, ChevronRightIcon, DotsHorizontalIcon, XIcon} from "@heroicons/react/outline";

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
            iconContent = <CheckIcon className="h-5 w-5"/>;
            break;
        case "FAILURE":
            iconStyle = "bg-red-600";
            iconContent = <XIcon className="h-5 w-5"/>;
            break;
        case "INACTIVE":
            iconStyle = "bg-gray-300";
            iconContent = <ChevronRightIcon className="h-5 w-5 ml-px"/>;
            break;
        case "IN_PROGRESS":
            iconStyle = "bg-yellow-500";
            iconContent = <DotsHorizontalIcon className="h-5 w-5"/>;
            break;
    }

    const [commitSubject] = data.commitMessage.split("\n", 1);

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
                                        {commitSubject}
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
