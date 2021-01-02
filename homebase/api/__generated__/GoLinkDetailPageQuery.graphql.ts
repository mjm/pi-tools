/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type GoLinkDetailPageQueryVariables = {
    id: string;
};
export type GoLinkDetailPageQueryResponse = {
    readonly viewer: {
        readonly link: {
            readonly id: string;
            readonly shortURL: string;
            readonly " $fragmentRefs": FragmentRefs<"EditLinkForm_link">;
        } | null;
    } | null;
};
export type GoLinkDetailPageQuery = {
    readonly response: GoLinkDetailPageQueryResponse;
    readonly variables: GoLinkDetailPageQueryVariables;
};



/*
query GoLinkDetailPageQuery(
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
    "name": "GoLinkDetailPageQuery",
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
    "name": "GoLinkDetailPageQuery",
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
    "cacheID": "23391993b7f547448e0ff1fb2df198f4",
    "id": null,
    "metadata": {},
    "name": "GoLinkDetailPageQuery",
    "operationKind": "query",
    "text": "query GoLinkDetailPageQuery(\n  $id: ID!\n) {\n  viewer {\n    link(id: $id) {\n      id\n      shortURL\n      ...EditLinkForm_link\n    }\n  }\n}\n\nfragment EditLinkForm_link on Link {\n  id\n  shortURL\n  destinationURL\n  description\n}\n"
  }
};
})();
(node as any).hash = '25de1d1e24ff9f9a3998947bb88217fe';
export default node;
