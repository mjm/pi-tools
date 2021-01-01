/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type TripsPageQueryVariables = {};
export type TripsPageQueryResponse = {
    readonly viewer: {
        readonly trips: {
            readonly edges: ReadonlyArray<{
                readonly node: {
                    readonly id: string;
                    readonly " $fragmentRefs": FragmentRefs<"TripRow_trip">;
                };
            }>;
        } | null;
        readonly " $fragmentRefs": FragmentRefs<"TagFilters_tags">;
    } | null;
};
export type TripsPageQuery = {
    readonly response: TripsPageQueryResponse;
    readonly variables: TripsPageQueryVariables;
};



/*
query TripsPageQuery {
  viewer {
    ...TagFilters_tags
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
*/

const node: ConcreteRequest = (function(){
var v0 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "id",
  "storageKey": null
},
v1 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "__typename",
  "storageKey": null
},
v2 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "cursor",
  "storageKey": null
},
v3 = {
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
},
v4 = [
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
    "name": "TripsPageQuery",
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
            "alias": "trips",
            "args": null,
            "concreteType": "TripConnection",
            "kind": "LinkedField",
            "name": "__TripsPageQuery_trips_connection",
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
                      (v0/*: any*/),
                      (v1/*: any*/),
                      {
                        "args": null,
                        "kind": "FragmentSpread",
                        "name": "TripRow_trip"
                      }
                    ],
                    "storageKey": null
                  },
                  (v2/*: any*/)
                ],
                "storageKey": null
              },
              (v3/*: any*/)
            ],
            "storageKey": null
          },
          {
            "args": null,
            "kind": "FragmentSpread",
            "name": "TagFilters_tags"
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
    "name": "TripsPageQuery",
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
            "args": (v4/*: any*/),
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
                      (v0/*: any*/),
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
                      (v1/*: any*/)
                    ],
                    "storageKey": null
                  },
                  (v2/*: any*/)
                ],
                "storageKey": null
              },
              (v3/*: any*/)
            ],
            "storageKey": "trips(first:30)"
          },
          {
            "alias": null,
            "args": (v4/*: any*/),
            "filters": null,
            "handle": "connection",
            "key": "TripsPageQuery_trips",
            "kind": "LinkedHandle",
            "name": "trips"
          }
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "ac44b7cc027417e5982504386d2e25ee",
    "id": null,
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "viewer",
            "trips"
          ]
        }
      ]
    },
    "name": "TripsPageQuery",
    "operationKind": "query",
    "text": "query TripsPageQuery {\n  viewer {\n    ...TagFilters_tags\n    trips(first: 30) {\n      edges {\n        node {\n          id\n          ...TripRow_trip\n          __typename\n        }\n        cursor\n      }\n      pageInfo {\n        endCursor\n        hasNextPage\n      }\n    }\n  }\n}\n\nfragment TagFilters_tags on Viewer {\n  tags(first: 5) {\n    edges {\n      node {\n        name\n        tripCount\n      }\n    }\n  }\n}\n\nfragment TripRow_trip on Trip {\n  id\n  leftAt\n  returnedAt\n  tags\n}\n"
  }
};
})();
(node as any).hash = 'd64148bf568f09c66b4fa7de319587fb';
export default node;
