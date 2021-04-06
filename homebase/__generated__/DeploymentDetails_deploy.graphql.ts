/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type DeploymentDetails_deploy = {
    readonly commitSHA: string;
    readonly commitMessage: string;
    readonly startedAt: string;
    readonly finishedAt: string | null;
    readonly " $refType": "DeploymentDetails_deploy";
};
export type DeploymentDetails_deploy$data = DeploymentDetails_deploy;
export type DeploymentDetails_deploy$key = {
    readonly " $data"?: DeploymentDetails_deploy$data;
    readonly " $fragmentRefs": FragmentRefs<"DeploymentDetails_deploy">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "DeploymentDetails_deploy",
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
  "type": "Deploy",
  "abstractKey": null
};
(node as any).hash = 'b88db183d9842e8d4aeb06449c480302';
export default node;
