/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type HomePageQueryVariables = {};
export type HomePageQueryResponse = {
    readonly viewer: {
        readonly " $fragmentRefs": FragmentRefs<"MostRecentTripCard_viewer" | "FiringAlertsCard_viewer" | "MostRecentDeployCard_viewer">;
    } | null;
};
export type HomePageQuery = {
    readonly response: HomePageQueryResponse;
    readonly variables: HomePageQueryVariables;
};



/*
query HomePageQuery {
  viewer {
    ...MostRecentTripCard_viewer
    ...FiringAlertsCard_viewer
    ...MostRecentDeployCard_viewer
  }
}

fragment FiringAlertsCard_viewer on Viewer {
  alerts {
    activeAt
    value
  }
}

fragment MostRecentDeployCard_viewer on Viewer {
  mostRecentDeploy {
    commitSHA
    commitMessage
    state
    startedAt
    finishedAt
    id
  }
}

fragment MostRecentTripCard_viewer on Viewer {
  trips(first: 1) {
    edges {
      node {
        leftAt
        returnedAt
        id
      }
    }
  }
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
    "name": "HomePageQuery",
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
            "name": "MostRecentTripCard_viewer"
          },
          {
            "args": null,
            "kind": "FragmentSpread",
            "name": "FiringAlertsCard_viewer"
          },
          {
            "args": null,
            "kind": "FragmentSpread",
            "name": "MostRecentDeployCard_viewer"
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
    "name": "HomePageQuery",
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
                "value": 1
              }
            ],
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
                      (v0/*: any*/)
                    ],
                    "storageKey": null
                  }
                ],
                "storageKey": null
              }
            ],
            "storageKey": "trips(first:1)"
          },
          {
            "alias": null,
            "args": null,
            "concreteType": "Alert",
            "kind": "LinkedField",
            "name": "alerts",
            "plural": true,
            "selections": [
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "activeAt",
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "value",
                "storageKey": null
              }
            ],
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "concreteType": "Deploy",
            "kind": "LinkedField",
            "name": "mostRecentDeploy",
            "plural": false,
            "selections": [
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "commitSHA",
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "commitMessage",
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "state",
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "startedAt",
                "storageKey": null
              },
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "finishedAt",
                "storageKey": null
              },
              (v0/*: any*/)
            ],
            "storageKey": null
          }
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "df1cb5443fab504b7a8b7feda87426a7",
    "id": null,
    "metadata": {},
    "name": "HomePageQuery",
    "operationKind": "query",
    "text": "query HomePageQuery {\n  viewer {\n    ...MostRecentTripCard_viewer\n    ...FiringAlertsCard_viewer\n    ...MostRecentDeployCard_viewer\n  }\n}\n\nfragment FiringAlertsCard_viewer on Viewer {\n  alerts {\n    activeAt\n    value\n  }\n}\n\nfragment MostRecentDeployCard_viewer on Viewer {\n  mostRecentDeploy {\n    commitSHA\n    commitMessage\n    state\n    startedAt\n    finishedAt\n    id\n  }\n}\n\nfragment MostRecentTripCard_viewer on Viewer {\n  trips(first: 1) {\n    edges {\n      node {\n        leftAt\n        returnedAt\n        id\n      }\n    }\n  }\n}\n"
  }
};
})();
(node as any).hash = 'd2f61b101b71eff4afe78d5885c3dbf9';
export default node;
