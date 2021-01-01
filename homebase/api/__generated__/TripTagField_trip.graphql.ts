/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type TripTagField_trip = {
    readonly rawID: string;
    readonly tags: ReadonlyArray<string>;
    readonly " $refType": "TripTagField_trip";
};
export type TripTagField_trip$data = TripTagField_trip;
export type TripTagField_trip$key = {
    readonly " $data"?: TripTagField_trip$data;
    readonly " $fragmentRefs": FragmentRefs<"TripTagField_trip">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "TripTagField_trip",
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
      "name": "tags",
      "storageKey": null
    }
  ],
  "type": "Trip",
  "abstractKey": null
};
(node as any).hash = '97122a6c4084e4f16efde01f89c3d24a';
export default node;
