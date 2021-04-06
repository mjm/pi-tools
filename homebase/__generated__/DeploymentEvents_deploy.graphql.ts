/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type DeploymentEvents_deploy = {
    readonly report: {
        readonly events: ReadonlyArray<{
            readonly " $fragmentRefs": FragmentRefs<"DeploymentEvent_event">;
        }>;
    } | null;
    readonly " $fragmentRefs": FragmentRefs<"DeploymentEvent_deploy">;
    readonly " $refType": "DeploymentEvents_deploy";
};
export type DeploymentEvents_deploy$data = DeploymentEvents_deploy;
export type DeploymentEvents_deploy$key = {
    readonly " $data"?: DeploymentEvents_deploy$data;
    readonly " $fragmentRefs": FragmentRefs<"DeploymentEvents_deploy">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "DeploymentEvents_deploy",
  "selections": [
    {
      "alias": null,
      "args": null,
      "concreteType": "DeployReport",
      "kind": "LinkedField",
      "name": "report",
      "plural": false,
      "selections": [
        {
          "alias": null,
          "args": null,
          "concreteType": "DeployEvent",
          "kind": "LinkedField",
          "name": "events",
          "plural": true,
          "selections": [
            {
              "args": null,
              "kind": "FragmentSpread",
              "name": "DeploymentEvent_event"
            }
          ],
          "storageKey": null
        }
      ],
      "storageKey": null
    },
    {
      "args": null,
      "kind": "FragmentSpread",
      "name": "DeploymentEvent_deploy"
    }
  ],
  "type": "Deploy",
  "abstractKey": null
};
(node as any).hash = '63ed6d31aee2b846e3b4b236db5a73e0';
export default node;
