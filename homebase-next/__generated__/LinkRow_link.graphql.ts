/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type LinkRow_link = {
    readonly id: string;
    readonly shortURL: string;
    readonly description: string;
    readonly " $refType": "LinkRow_link";
};
export type LinkRow_link$data = LinkRow_link;
export type LinkRow_link$key = {
    readonly " $data"?: LinkRow_link$data;
    readonly " $fragmentRefs": FragmentRefs<"LinkRow_link">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "LinkRow_link",
  "selections": [
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "id",
      "storageKey": null
    },
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "shortURL",
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
  "type": "Link",
  "abstractKey": null
};
(node as any).hash = '7dee66937f9d2d22bff6b14a065743ca';
export default node;
