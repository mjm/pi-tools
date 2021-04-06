import {graphql, useFragment} from "react-relay";
import {DeploymentEvents_deploy$key} from "../../__generated__/DeploymentEvents_deploy.graphql";
import DeploymentEvent from "./DeploymentEvent";

export default function DeploymentEvents({deploy}: { deploy: DeploymentEvents_deploy$key }) {
    const data = useFragment(
        graphql`
            fragment DeploymentEvents_deploy on Deploy {
                ...DeploymentEvent_deploy
                report {
                    events {
                        ...DeploymentEvent_event
                    }
                }
            }
        `,
        deploy,
    );

    const report = data.report;
    if (!report) {
        return null;
    }

    return (
        <div className="bg-white shadow overflow-hidden sm:rounded-lg mt-8">
            <div className="bg-white px-4 py-5 border-b border-gray-200 sm:px-6">
                <div className="-ml-4 -mt-2 flex items-center justify-between flex-wrap sm:flex-nowrap">
                    <div className="ml-4 mt-2">
                        <h3 className="text-lg leading-6 font-medium text-gray-900">
                            Timeline
                        </h3>
                    </div>
                </div>
            </div>

            <div className="flow-root py-4 px-4 sm:px-6">
                <ul className="-mb-8">
                    {data.report.events.map((event, idx) => (
                        <DeploymentEvent
                            deploy={data}
                            event={event}
                            key={idx}
                            isLast={idx === data.report.events.length - 1}
                        />
                    ))}
                </ul>
            </div>
        </div>
    );
}
