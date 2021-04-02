/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type Id_GoLinkDetailPageQueryVariables = {
    id: string;
};
export type Id_GoLinkDetailPageQueryResponse = {
    readonly viewer: {
        readonly link: {
            readonly id: string;
            readonly shortURL: string;
            readonly " $fragmentRefs": FragmentRefs<"EditLinkForm_link">;
        } | null;
    } | null;
};
export type Id_GoLinkDetailPageQuery = {
    readonly response: Id_GoLinkDetailPageQueryResponse;
    readonly variables: Id_GoLinkDetailPageQueryVariables;
};



/*
query Id_GoLinkDetailPageQuery(
  $id: ID!
) {
  viewer {
    link(id: $id) {
      id
      shortURL
      ...EditLinkForm_link
    }
  }
}

fragment EditLinkForm_link on Link {
  id
  shortURL
  destinationURL
  description
}
*/

const node: ConcreteRequest = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "id"
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "id",
    "variableName": "id"
  }
],
v2 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "id",
  "storageKey": null
},
v3 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "shortURL",
  "storageKey": null
};
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "Id_GoLinkDetailPageQuery",
    "selections": [
      {
        "alias": null,
        "args": null,
        "concreteType": "Viewer",
        "kind": "LinkedField",
        "name": "viewer",
        "plural": false,
        "selections": [
          {
            "alias": null,
            "args": (v1/*: any*/),
            "concreteType": "Link",
            "kind": "LinkedField",
            "name": "link",
            "plural": false,
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              {
                "args": null,
                "kind": "FragmentSpread",
                "name": "EditLinkForm_link"
              }
            ],
            "storageKey": null
          }
        ],
        "storageKey": null
      }
    ],
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "Id_GoLinkDetailPageQuery",
    "selections": [
      {
        "alias": null,
        "args": null,
        "concreteType": "Viewer",
        "kind": "LinkedField",
        "name": "viewer",
        "plural": false,
        "selections": [
          {
            "alias": null,
            "args": (v1/*: any*/),
            "concreteType": "Link",
            "kind": "LinkedField",
            "name": "link",
            "plural": false,
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
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
            "storageKey": null
          }
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "c2da8970e82160a0bd12d5a43c8c19c4",
    "id": null,
    "metadata": {},
    "name": "Id_GoLinkDetailPageQuery",
    "operationKind": "query",
    "text": "query Id_GoLinkDetailPageQuery(\n  $id: ID!\n) {\n  viewer {\n    link(id: $id) {\n      id\n      shortURL\n      ...EditLinkForm_link\n    }\n  }\n}\n\nfragment EditLinkForm_link on Link {\n  id\n  shortURL\n  destinationURL\n  description\n}\n"
  }
};
})();
(node as any).hash = '56a51b2fdd5140e4cd6f326c3c2feef9';
export default node;
