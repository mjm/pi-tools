/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type RecentLinksList_links = {
    readonly edges: ReadonlyArray<{
        readonly node: {
            readonly id: string;
            readonly " $fragmentRefs": FragmentRefs<"LinkRow_link">;
        };
    }>;
    readonly " $refType": "RecentLinksList_links";
};
export type RecentLinksList_links$data = RecentLinksList_links;
export type RecentLinksList_links$key = {
    readonly " $data"?: RecentLinksList_links$data;
    readonly " $fragmentRefs": FragmentRefs<"RecentLinksList_links">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "RecentLinksList_links",
  "selections": [
    {
      "alias": null,
      "args": null,
      "concreteType": "LinkEdge",
      "kind": "LinkedField",
      "name": "edges",
      "plural": true,
      "selections": [
        {
          "alias": null,
          "args": null,
          "concreteType": "Link",
          "kind": "LinkedField",
          "name": "node",
          "plural": false,
          "selections": [
            {
              "alias": null,
              "args": null,
              "kind": "ScalarField",
              "name": "id",
              "storageKey": null
            },
            {
              "args": null,
              "kind": "FragmentSpread",
              "name": "LinkRow_link"
            }
          ],
          "storageKey": null
        }
      ],
      "storageKey": null
    }
  ],
  "type": "LinkConnection",
  "abstractKey": null
};
(node as any).hash = '69f47fa689e6193125a8d3b09228679b';
export default node;
