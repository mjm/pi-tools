/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type LinkRow_link = {
    readonly rawID: string;
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
      "name": "rawID",
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
(node as any).hash = '8fe0a85b0b9ac9dfb523783536eb3b2c';
export default node;
