/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type EditLinkForm_link = {
    readonly id: string;
    readonly shortURL: string;
    readonly destinationURL: string;
    readonly description: string;
    readonly " $refType": "EditLinkForm_link";
};
export type EditLinkForm_link$data = EditLinkForm_link;
export type EditLinkForm_link$key = {
    readonly " $data"?: EditLinkForm_link$data;
    readonly " $fragmentRefs": FragmentRefs<"EditLinkForm_link">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "EditLinkForm_link",
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
      "name": "destinationURL",
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
(node as any).hash = 'e9074bcfd75f01c31ad02ca00993e9c5';
export default node;
