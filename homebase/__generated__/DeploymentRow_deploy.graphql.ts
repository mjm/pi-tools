/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type DeployState = "FAILURE" | "INACTIVE" | "IN_PROGRESS" | "PENDING" | "SUCCESS" | "UNKNOWN" | "%future added value";
export type DeploymentRow_deploy = {
    readonly rawID: string;
    readonly state: DeployState;
    readonly commitSHA: string;
    readonly commitMessage: string;
    readonly startedAt: string;
    readonly " $refType": "DeploymentRow_deploy";
};
export type DeploymentRow_deploy$data = DeploymentRow_deploy;
export type DeploymentRow_deploy$key = {
    readonly " $data"?: DeploymentRow_deploy$data;
    readonly " $fragmentRefs": FragmentRefs<"DeploymentRow_deploy">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "DeploymentRow_deploy",
  "selections": [
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "rawID",
      "storageKey": null
    },
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "state",
      "storageKey": null
    },
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "commitSHA",
      "storageKey": null
    },
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "commitMessage",
      "storageKey": null
    },
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "startedAt",
      "storageKey": null
    }
  ],
  "type": "Deploy",
  "abstractKey": null
};
(node as any).hash = '1167cb8cb7f711b7b649b56dfad609c2';
export default node;
