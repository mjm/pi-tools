/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type TagFilters_tags = {
    readonly tags: {
        readonly edges: ReadonlyArray<{
            readonly node: {
                readonly name: string;
                readonly tripCount: number;
            };
        }>;
    } | null;
    readonly " $refType": "TagFilters_tags";
};
export type TagFilters_tags$data = TagFilters_tags;
export type TagFilters_tags$key = {
    readonly " $data"?: TagFilters_tags$data;
    readonly " $fragmentRefs": FragmentRefs<"TagFilters_tags">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "TagFilters_tags",
  "selections": [
    {
      "alias": null,
      "args": [
        {
          "kind": "Literal",
          "name": "first",
          "value": 5
        }
      ],
      "concreteType": "TagConnection",
      "kind": "LinkedField",
      "name": "tags",
      "plural": false,
      "selections": [
        {
          "alias": null,
          "args": null,
          "concreteType": "TagEdge",
          "kind": "LinkedField",
          "name": "edges",
          "plural": true,
          "selections": [
            {
              "alias": null,
              "args": null,
              "concreteType": "Tag",
              "kind": "LinkedField",
              "name": "node",
              "plural": false,
              "selections": [
                {
                  "alias": null,
                  "args": null,
                  "kind": "ScalarField",
                  "name": "name",
                  "storageKey": null
                },
                {
                  "alias": null,
                  "args": null,
                  "kind": "ScalarField",
                  "name": "tripCount",
                  "storageKey": null
                }
              ],
              "storageKey": null
            }
          ],
          "storageKey": null
        }
      ],
      "storageKey": "tags(first:5)"
    }
  ],
  "type": "Viewer",
  "abstractKey": null
};
(node as any).hash = 'f01dbbc89c7b03c223ad4137d7536831';
export default node;
