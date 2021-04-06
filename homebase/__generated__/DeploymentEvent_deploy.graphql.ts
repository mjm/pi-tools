/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type DeploymentEvent_deploy = {
    readonly startedAt: string;
    readonly " $refType": "DeploymentEvent_deploy";
};
export type DeploymentEvent_deploy$data = DeploymentEvent_deploy;
export type DeploymentEvent_deploy$key = {
    readonly " $data"?: DeploymentEvent_deploy$data;
    readonly " $fragmentRefs": FragmentRefs<"DeploymentEvent_deploy">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "DeploymentEvent_deploy",
  "selections": [
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
(node as any).hash = 'b3be188dedd2b7b3318b943ab97c8316';
export default node;
