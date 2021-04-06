import {graphql, useFragment} from "react-relay";
import {addHours, getHours, getMinutes, getSeconds, getUnixTime, parseISO} from "date-fns";
import {DeploymentEvent_event$key} from "../../__generated__/DeploymentEvent_event.graphql";
import {DeploymentEvent_deploy$key} from "../../__generated__/DeploymentEvent_deploy.graphql";

export default function DeploymentEvent({deploy, event, isLast}: {
    deploy: DeploymentEvent_deploy$key;
    event: DeploymentEvent_event$key;
    isLast: boolean;
}) {
    const deployData = useFragment(
        graphql`
            fragment DeploymentEvent_deploy on Deploy {
                startedAt
            }
        `,
        deploy,
    );

    const data = useFragment(
        graphql`
            fragment DeploymentEvent_event on DeployEvent {
                timestamp
                level
                summary
                description
            }
        `,
        event,
    );

    const eventTime = parseISO(data.timestamp);
    const deployStartTime = parseISO(deployData.startedAt);
    const secondsSinceStart = getUnixTime(eventTime) - getUnixTime(deployStartTime);

    return (
        <li>
            <div className="relative pb-8">
                {!isLast && <span className="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200"
                                  aria-hidden="true"/>}
                <div className="relative flex space-x-3">
                    <div>
            <span className="h-8 w-8 rounded-full bg-gray-400 flex items-center justify-center ring-8 ring-white">
<svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7"/>
</svg>
            </span>
                    </div>
                    <div className="min-w-0 flex-1 pt-1.5 flex justify-between space-x-4">
                        <div>
                            <p className="text-sm text-gray-700">{data.summary}</p>
                        </div>
                        <div className="text-right text-sm whitespace-nowrap text-gray-500">
                            <time dateTime={data.timestamp}>{convertToDuration(secondsSinceStart)}</time>
                        </div>
                    </div>
                </div>
            </div>
        </li>
    );
}

function convertToDuration(secondsAmount: number) {
    const normalizeTime = (time: string): string =>
        time.length === 1 ? `0${time}` : time;

    const SECONDS_TO_MILLISECONDS_COEFF = 1000;
    const MINUTES_IN_HOUR = 60;

    const milliseconds = secondsAmount * SECONDS_TO_MILLISECONDS_COEFF;

    const date = new Date(milliseconds);
    const timezoneDiff = date.getTimezoneOffset() / MINUTES_IN_HOUR;
    const dateWithoutTimezoneDiff = addHours(date, timezoneDiff);

    const hours = normalizeTime(String(getHours(dateWithoutTimezoneDiff)));
    const minutes = normalizeTime(String(getMinutes(dateWithoutTimezoneDiff)));
    const seconds = normalizeTime(String(getSeconds(dateWithoutTimezoneDiff)));

    const hoursOutput = hours !== "00" ? `${hours}:` : "";

    return `${hoursOutput}${minutes}:${seconds}`;
};
