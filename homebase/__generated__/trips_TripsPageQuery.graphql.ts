/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type trips_TripsPageQueryVariables = {};
export type trips_TripsPageQueryResponse = {
    readonly viewer: {
        readonly " $fragmentRefs": FragmentRefs<"TagFilters_tags" | "TripsList_viewer">;
    } | null;
};
export type trips_TripsPageQuery = {
    readonly response: trips_TripsPageQueryResponse;
    readonly variables: trips_TripsPageQueryVariables;
};



/*
query trips_TripsPageQuery {
  viewer {
    ...TagFilters_tags
    ...TripsList_viewer
  }
}

fragment TagFilters_tags on Viewer {
  tags(first: 5) {
    edges {
      node {
        name
        tripCount
      }
    }
  }
}

fragment TripRow_trip on Trip {
  id
  leftAt
  returnedAt
  tags
}

fragment TripsList_viewer on Viewer {
  trips(first: 30) {
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
    "kind": "Literal",
    "name": "first",
    "value": 30
  }
];
return {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "trips_TripsPageQuery",
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
            "args": null,
            "kind": "FragmentSpread",
            "name": "TagFilters_tags"
          },
          {
            "args": null,
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
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "trips_TripsPageQuery",
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
            "args": [
              {
                "kind": "Literal",
                "name": "first",
                "value": 5
              }
            ],
            "concreteType": "TagConnection",
            "kind": "LinkedField",
            "name": "tags",
            "plural": false,
            "selections": [
              {
                "alias": null,
                "args": null,
                "concreteType": "TagEdge",
                "kind": "LinkedField",
                "name": "edges",
                "plural": true,
                "selections": [
                  {
                    "alias": null,
                    "args": null,
                    "concreteType": "Tag",
                    "kind": "LinkedField",
                    "name": "node",
                    "plural": false,
                    "selections": [
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
                        "name": "tripCount",
                        "storageKey": null
                      }
                    ],
                    "storageKey": null
                  }
                ],
                "storageKey": null
              }
            ],
            "storageKey": "tags(first:5)"
          },
          {
            "alias": null,
            "args": (v0/*: any*/),
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
            "storageKey": "trips(first:30)"
          },
          {
            "alias": null,
            "args": (v0/*: any*/),
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
    "cacheID": "0d216bb8e3b4519b847cc9d27c9a9f1f",
    "id": null,
    "metadata": {},
    "name": "trips_TripsPageQuery",
    "operationKind": "query",
    "text": "query trips_TripsPageQuery {\n  viewer {\n    ...TagFilters_tags\n    ...TripsList_viewer\n  }\n}\n\nfragment TagFilters_tags on Viewer {\n  tags(first: 5) {\n    edges {\n      node {\n        name\n        tripCount\n      }\n    }\n  }\n}\n\nfragment TripRow_trip on Trip {\n  id\n  leftAt\n  returnedAt\n  tags\n}\n\nfragment TripsList_viewer on Viewer {\n  trips(first: 30) {\n    edges {\n      node {\n        id\n        ...TripRow_trip\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n"
  }
};
})();
(node as any).hash = 'c0af8af7730e92f0da526f80c09fb844';
export default node;
