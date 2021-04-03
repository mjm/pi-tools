/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type MostRecentTripCard_viewer = {
    readonly trips: {
        readonly edges: ReadonlyArray<{
            readonly node: {
                readonly leftAt: string;
                readonly returnedAt: string | null;
            };
        }>;
    } | null;
    readonly " $refType": "MostRecentTripCard_viewer";
};
export type MostRecentTripCard_viewer$data = MostRecentTripCard_viewer;
export type MostRecentTripCard_viewer$key = {
    readonly " $data"?: MostRecentTripCard_viewer$data;
    readonly " $fragmentRefs": FragmentRefs<"MostRecentTripCard_viewer">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "MostRecentTripCard_viewer",
  "selections": [
    {
      "alias": null,
      "args": [
        {
          "kind": "Literal",
          "name": "first",
          "value": 1
        }
      ],
      "concreteType": "TripConnection",
      "kind": "LinkedField",
      "name": "trips",
      "plural": false,
      "selections": [
        {
          "alias": null,
          "args": null,
          "concreteType": "TripEdge",
          "kind": "LinkedField",
          "name": "edges",
          "plural": true,
          "selections": [
            {
              "alias": null,
              "args": null,
              "concreteType": "Trip",
              "kind": "LinkedField",
              "name": "node",
              "plural": false,
              "selections": [
                {
                  "alias": null,
                  "args": null,
                  "kind": "ScalarField",
                  "name": "leftAt",
                  "storageKey": null
                },
                {
                  "alias": null,
                  "args": null,
                  "kind": "ScalarField",
                  "name": "returnedAt",
                  "storageKey": null
                }
              ],
              "storageKey": null
            }
          ],
          "storageKey": null
        }
      ],
      "storageKey": "trips(first:1)"
    }
  ],
  "type": "Viewer",
  "abstractKey": null
};
(node as any).hash = '21d790b015ab4e1fc43092bd2358d53b';
export default node;
