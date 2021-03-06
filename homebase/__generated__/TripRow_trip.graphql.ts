/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type TripRow_trip = {
    readonly id: string;
    readonly leftAt: string;
    readonly returnedAt: string | null;
    readonly tags: ReadonlyArray<string>;
    readonly " $refType": "TripRow_trip";
};
export type TripRow_trip$data = TripRow_trip;
export type TripRow_trip$key = {
    readonly " $data"?: TripRow_trip$data;
    readonly " $fragmentRefs": FragmentRefs<"TripRow_trip">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "TripRow_trip",
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
      "name": "leftAt",
      "storageKey": null
    },
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "returnedAt",
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
(node as any).hash = 'aaeef382597cb2259305a02777313460';
export default node;
