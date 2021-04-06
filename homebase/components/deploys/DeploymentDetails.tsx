import DescriptionField from "../DescriptionField";
import {format, formatDuration, intervalToDuration, parseISO} from "date-fns";
import {graphql, useFragment} from "react-relay";
import {DeploymentDetails_deploy$key} from "../../__generated__/DeploymentDetails_deploy.graphql";

export default function DeploymentDetails({deploy}: { deploy: DeploymentDetails_deploy$key }) {
    const data = useFragment(
        graphql`
            fragment DeploymentDetails_deploy on Deploy {
                commitSHA
                commitMessage
                startedAt
                finishedAt
            }
        `,
        deploy,
    );

    return (
        <div>
            <dl>
                <DescriptionField label="Commit" offset>
                    <a href={`https://github.com/mjm/pi-tools/commit/${data.commitSHA}`}
                       className="font-medium text-indigo-600 hover:text-indigo-500"
                       target="_blank">
                        {data.commitMessage}
                    </a>
                </DescriptionField>
                <DescriptionField label="Started at">
                    {format(parseISO(data.startedAt), "PPpp")}
                </DescriptionField>
                <DescriptionField label="Finished at" offset>
                    {format(parseISO(data.finishedAt), "PPpp")}
                </DescriptionField>
                <DescriptionField label="Duration">
                    {formatDuration(intervalToDuration({
                        start: parseISO(data.startedAt),
                        end: parseISO(data.finishedAt),
                    }))}
                </DescriptionField>
            </dl>
        </div>
    );
}
