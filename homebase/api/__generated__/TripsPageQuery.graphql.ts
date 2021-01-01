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
    trips {
      edges {
        node {
          id
          ...TripRow_trip
        }
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
};
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
            "alias": null,
            "args": null,
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
                        "args": null,
                        "kind": "FragmentSpread",
                        "name": "TripRow_trip"
                      }
                    ],
                    "storageKey": null
                  }
                ],
                "storageKey": null
              }
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
            "args": null,
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
                      }
                    ],
                    "storageKey": null
                  }
                ],
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
    "cacheID": "182e8409dbfdd548f633ad837becc1ea",
    "id": null,
    "metadata": {},
    "name": "TripsPageQuery",
    "operationKind": "query",
    "text": "query TripsPageQuery {\n  viewer {\n    ...TagFilters_tags\n    trips {\n      edges {\n        node {\n          id\n          ...TripRow_trip\n        }\n      }\n    }\n  }\n}\n\nfragment TagFilters_tags on Viewer {\n  tags(first: 5) {\n    edges {\n      node {\n        name\n        tripCount\n      }\n    }\n  }\n}\n\nfragment TripRow_trip on Trip {\n  id\n  leftAt\n  returnedAt\n  tags\n}\n"
  }
};
})();
(node as any).hash = '9551198261ba8d7caaddce9137680661';
export default node;
