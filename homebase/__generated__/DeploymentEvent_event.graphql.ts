/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type DeployEventLevel = "ERROR" | "INFO" | "UNKNOWN" | "WARNING" | "%future added value";
export type DeploymentEvent_event = {
    readonly timestamp: string;
    readonly level: DeployEventLevel;
    readonly summary: string;
    readonly description: string;
    readonly " $refType": "DeploymentEvent_event";
};
export type DeploymentEvent_event$data = DeploymentEvent_event;
export type DeploymentEvent_event$key = {
    readonly " $data"?: DeploymentEvent_event$data;
    readonly " $fragmentRefs": FragmentRefs<"DeploymentEvent_event">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "DeploymentEvent_event",
  "selections": [
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "timestamp",
      "storageKey": null
    },
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "level",
      "storageKey": null
    },
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "summary",
      "storageKey": null
    },
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "description",
      "storageKey": null
    }
  ],
  "type": "DeployEvent",
  "abstractKey": null
};
(node as any).hash = '0ed057974c1f305748652b9093d8930c';
export default node;
