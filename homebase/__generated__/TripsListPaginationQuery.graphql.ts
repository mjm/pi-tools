/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type TripsListPaginationQueryVariables = {
    count?: number | null;
    cursor?: string | null;
};
export type TripsListPaginationQueryResponse = {
    readonly viewer: {
        readonly " $fragmentRefs": FragmentRefs<"TripsList_viewer">;
    } | null;
};
export type TripsListPaginationQuery = {
    readonly response: TripsListPaginationQueryResponse;
    readonly variables: TripsListPaginationQueryVariables;
};



/*
query TripsListPaginationQuery(
  $count: Int = 30
  $cursor: Cursor
) {
  viewer {
    ...TripsList_viewer_1G22uz
  }
}

fragment TripRow_trip on Trip {
  id
  leftAt
  returnedAt
  tags
}

fragment TripsList_viewer_1G22uz on Viewer {
  trips(first: $count, after: $cursor) {
    edges {
      node {
        id
        ...TripRow_trip
        __typename
      }
      cursor
    }
    pageInfo {
      endCursor
      hasNextPage
    }
  }
}
*/

const node: ConcreteRequest = (function(){
var v0 = [
  {
    "defaultValue": 30,
    "kind": "LocalArgument",
    "name": "count"
  },
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "cursor"
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "after",
    "variableName": "cursor"
  },
  {
    "kind": "Variable",
    "name": "first",
    "variableName": "count"
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "TripsListPaginationQuery",
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
            "args": [
              {
                "kind": "Variable",
                "name": "count",
                "variableName": "count"
              },
              {
                "kind": "Variable",
                "name": "cursor",
                "variableName": "cursor"
              }
            ],
            "kind": "FragmentSpread",
            "name": "TripsList_viewer"
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
    "name": "TripsListPaginationQuery",
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
            "concreteType": "TripConnection",
            "kind": "LinkedField",
            "name": "trips",
            "plural": false,
            "selections": [
              {
                "alias": null,
                "args": null,
                "concreteType": "TripEdge",
                "kind": "LinkedField",
                "name": "edges",
                "plural": true,
                "selections": [
                  {
                    "alias": null,
                    "args": null,
                    "concreteType": "Trip",
                    "kind": "LinkedField",
                    "name": "node",
                    "plural": false,
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
                      },
                      {
                        "alias": null,
                        "args": null,
                        "kind": "ScalarField",
                        "name": "__typename",
                        "storageKey": null
                      }
                    ],
                    "storageKey": null
                  },
                  {
                    "alias": null,
                    "args": null,
                    "kind": "ScalarField",
                    "name": "cursor",
                    "storageKey": null
                  }
                ],
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "concreteType": "PageInfo",
                "kind": "LinkedField",
                "name": "pageInfo",
                "plural": false,
                "selections": [
                  {
                    "alias": null,
                    "args": null,
                    "kind": "ScalarField",
                    "name": "endCursor",
                    "storageKey": null
                  },
                  {
                    "alias": null,
                    "args": null,
                    "kind": "ScalarField",
                    "name": "hasNextPage",
                    "storageKey": null
                  }
                ],
                "storageKey": null
              }
            ],
            "storageKey": null
          },
          {
            "alias": null,
            "args": (v1/*: any*/),
            "filters": null,
            "handle": "connection",
            "key": "TripsList_trips",
            "kind": "LinkedHandle",
            "name": "trips"
          }
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "3b95428cebc66292fdaf9b772147b0ad",
    "id": null,
    "metadata": {},
    "name": "TripsListPaginationQuery",
    "operationKind": "query",
    "text": "query TripsListPaginationQuery(\n  $count: Int = 30\n  $cursor: Cursor\n) {\n  viewer {\n    ...TripsList_viewer_1G22uz\n  }\n}\n\nfragment TripRow_trip on Trip {\n  id\n  leftAt\n  returnedAt\n  tags\n}\n\nfragment TripsList_viewer_1G22uz on Viewer {\n  trips(first: $count, after: $cursor) {\n    edges {\n      node {\n        id\n        ...TripRow_trip\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n"
  }
};
})();
(node as any).hash = 'c3bf7f89215f8ae44f6c2e9b556a2451';
export default node;
