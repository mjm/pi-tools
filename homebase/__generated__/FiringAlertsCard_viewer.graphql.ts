/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type FiringAlertsCard_viewer = {
    readonly alerts: ReadonlyArray<{
        readonly activeAt: string;
        readonly value: string;
    }>;
    readonly " $refType": "FiringAlertsCard_viewer";
};
export type FiringAlertsCard_viewer$data = FiringAlertsCard_viewer;
export type FiringAlertsCard_viewer$key = {
    readonly " $data"?: FiringAlertsCard_viewer$data;
    readonly " $fragmentRefs": FragmentRefs<"FiringAlertsCard_viewer">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "FiringAlertsCard_viewer",
  "selections": [
    {
      "alias": null,
      "args": null,
      "concreteType": "Alert",
      "kind": "LinkedField",
      "name": "alerts",
      "plural": true,
      "selections": [
        {
          "alias": null,
          "args": null,
          "kind": "ScalarField",
          "name": "activeAt",
          "storageKey": null
        },
        {
          "alias": null,
          "args": null,
          "kind": "ScalarField",
          "name": "value",
          "storageKey": null
        }
      ],
      "storageKey": null
    }
  ],
  "type": "Viewer",
  "abstractKey": null
};
(node as any).hash = '73cc619df12475123f735fbf97a68ae8';
export default node;
