/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type Id_TripPageQueryVariables = {
    id: string;
};
export type Id_TripPageQueryResponse = {
    readonly viewer: {
        readonly trip: {
            readonly id: string;
            readonly leftAt: string;
            readonly returnedAt: string | null;
            readonly " $fragmentRefs": FragmentRefs<"TripTagField_trip">;
        } | null;
    } | null;
};
export type Id_TripPageQuery = {
    readonly response: Id_TripPageQueryResponse;
    readonly variables: Id_TripPageQueryVariables;
};



/*
query Id_TripPageQuery(
  $id: ID!
) {
  viewer {
    trip(id: $id) {
      id
      leftAt
      returnedAt
      ...TripTagField_trip
    }
  }
}

fragment TripTagField_trip on Trip {
  id
  tags
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
  "name": "leftAt",
  "storageKey": null
},
v4 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "returnedAt",
  "storageKey": null
};
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "Id_TripPageQuery",
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
            "concreteType": "Trip",
            "kind": "LinkedField",
            "name": "trip",
            "plural": false,
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              (v4/*: any*/),
              {
                "args": null,
                "kind": "FragmentSpread",
                "name": "TripTagField_trip"
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
    "name": "Id_TripPageQuery",
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
            "concreteType": "Trip",
            "kind": "LinkedField",
            "name": "trip",
            "plural": false,
            "selections": [
              (v2/*: any*/),
              (v3/*: any*/),
              (v4/*: any*/),
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "tags",
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
    "cacheID": "38161d9df5452d35e6b9ce78aba8a454",
    "id": null,
    "metadata": {},
    "name": "Id_TripPageQuery",
    "operationKind": "query",
    "text": "query Id_TripPageQuery(\n  $id: ID!\n) {\n  viewer {\n    trip(id: $id) {\n      id\n      leftAt\n      returnedAt\n      ...TripTagField_trip\n    }\n  }\n}\n\nfragment TripTagField_trip on Trip {\n  id\n  tags\n}\n"
  }
};
})();
(node as any).hash = '22b099ea2007a82defb538820df7bfc0';
export default node;
