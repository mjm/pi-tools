/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ReaderFragment } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type ArchiveRow_archive = {
    readonly id: string;
    readonly name: string;
    readonly createdAt: string;
    readonly " $refType": "ArchiveRow_archive";
};
export type ArchiveRow_archive$data = ArchiveRow_archive;
export type ArchiveRow_archive$key = {
    readonly " $data"?: ArchiveRow_archive$data;
    readonly " $fragmentRefs": FragmentRefs<"ArchiveRow_archive">;
};



const node: ReaderFragment = {
  "argumentDefinitions": [],
  "kind": "Fragment",
  "metadata": null,
  "name": "ArchiveRow_archive",
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
      "name": "name",
      "storageKey": null
    },
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "createdAt",
      "storageKey": null
    }
  ],
  "type": "Archive",
  "abstractKey": null
};
(node as any).hash = '86b02033938ba9677cd28bb2bd77f390';
export default node;
