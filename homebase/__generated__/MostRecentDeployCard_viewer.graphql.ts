/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type DeployState = "FAILURE" | "INACTIVE" | "IN_PROGRESS" | "PENDING" | "SUCCESS" | "UNKNOWN" | "%future added value";
export type MostRecentDeployCard_viewer = {
    readonly mostRecentDeploy: {
        readonly commitSHA: string;
        readonly commitMessage: string;
        readonly state: DeployState;
        readonly startedAt: string;
        readonly finishedAt: string | null;
    } | null;
    readonly " $refType": "MostRecentDeployCard_viewer";
};
export type MostRecentDeployCard_viewer$data = MostRecentDeployCard_viewer;
export type MostRecentDeployCard_viewer$key = {
    readonly " $data"?: MostRecentDeployCard_viewer$data;
    readonly " $fragmentRefs": FragmentRefs<"MostRecentDeployCard_viewer">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "MostRecentDeployCard_viewer",
  "selections": [
    {
      "alias": null,
      "args": null,
      "concreteType": "Deploy",
      "kind": "LinkedField",
      "name": "mostRecentDeploy",
      "plural": false,
      "selections": [
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
          "name": "state",
          "storageKey": null
        },
        {
          "alias": null,
          "args": null,
          "kind": "ScalarField",
          "name": "startedAt",
          "storageKey": null
        },
        {
          "alias": null,
          "args": null,
          "kind": "ScalarField",
          "name": "finishedAt",
          "storageKey": null
        }
      ],
      "storageKey": null
    }
  ],
  "type": "Viewer",
  "abstractKey": null
};
(node as any).hash = 'f1e12948df73cd5fe05fbbab9859bdfb';
export default node;
