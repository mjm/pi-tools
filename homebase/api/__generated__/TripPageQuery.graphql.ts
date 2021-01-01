/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type TripPageQueryVariables = {
    id: string;
};
export type TripPageQueryResponse = {
    readonly viewer: {
        readonly trip: {
            readonly id: string;
            readonly rawID: string;
            readonly leftAt: string;
            readonly returnedAt: string | null;
            readonly " $fragmentRefs": FragmentRefs<"TripTagField_trip">;
        } | null;
    } | null;
};
export type TripPageQuery = {
    readonly response: TripPageQueryResponse;
    readonly variables: TripPageQueryVariables;
};



/*
query TripPageQuery(
  $id: ID!
) {
  viewer {
    trip(id: $id) {
      id
      rawID
      leftAt
      returnedAt
      ...TripTagField_trip
    }
  }
}

fragment TripTagField_trip on Trip {
  rawID
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
  "name": "rawID",
  "storageKey": null
},
v4 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "leftAt",
  "storageKey": null
},
v5 = {
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
    "name": "TripPageQuery",
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
              (v5/*: any*/),
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
    "name": "TripPageQuery",
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
              (v5/*: any*/),
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
    "cacheID": "b35ef26368d7d746be983021d28cfca1",
    "id": null,
    "metadata": {},
    "name": "TripPageQuery",
    "operationKind": "query",
    "text": "query TripPageQuery(\n  $id: ID!\n) {\n  viewer {\n    trip(id: $id) {\n      id\n      rawID\n      leftAt\n      returnedAt\n      ...TripTagField_trip\n    }\n  }\n}\n\nfragment TripTagField_trip on Trip {\n  rawID\n  tags\n}\n"
  }
};
})();
(node as any).hash = 'c250a0fb28d14ff5160605f34e3609d7';
export default node;
