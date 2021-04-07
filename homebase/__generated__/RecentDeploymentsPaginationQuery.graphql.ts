/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type RecentDeploymentsPaginationQueryVariables = {
    count?: number | null;
    cursor?: string | null;
};
export type RecentDeploymentsPaginationQueryResponse = {
    readonly viewer: {
        readonly " $fragmentRefs": FragmentRefs<"RecentDeployments_viewer">;
    } | null;
};
export type RecentDeploymentsPaginationQuery = {
    readonly response: RecentDeploymentsPaginationQueryResponse;
    readonly variables: RecentDeploymentsPaginationQueryVariables;
};



/*
query RecentDeploymentsPaginationQuery(
  $count: Int
  $cursor: Cursor
) {
  viewer {
    ...RecentDeployments_viewer_1G22uz
  }
}

fragment DeploymentRow_deploy on Deploy {
  rawID
  state
  commitSHA
  commitMessage
  startedAt
}

fragment RecentDeployments_viewer_1G22uz on Viewer {
  recentDeploys(first: $count, after: $cursor) {
    edges {
      node {
        id
        ...DeploymentRow_deploy
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
    "defaultValue": null,
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
    "name": "RecentDeploymentsPaginationQuery",
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
            "name": "RecentDeployments_viewer"
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
    "name": "RecentDeploymentsPaginationQuery",
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
            "concreteType": "DeployConnection",
            "kind": "LinkedField",
            "name": "recentDeploys",
            "plural": false,
            "selections": [
              {
                "alias": null,
                "args": null,
                "concreteType": "DeployEdge",
                "kind": "LinkedField",
                "name": "edges",
                "plural": true,
                "selections": [
                  {
                    "alias": null,
                    "args": null,
                    "concreteType": "Deploy",
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
                        "name": "rawID",
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
                        "name": "startedAt",
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
            "key": "RecentDeployments_recentDeploys",
            "kind": "LinkedHandle",
            "name": "recentDeploys"
          }
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "9a64a2448d572ea60b285934312718ae",
    "id": null,
    "metadata": {},
    "name": "RecentDeploymentsPaginationQuery",
    "operationKind": "query",
    "text": "query RecentDeploymentsPaginationQuery(\n  $count: Int\n  $cursor: Cursor\n) {\n  viewer {\n    ...RecentDeployments_viewer_1G22uz\n  }\n}\n\nfragment DeploymentRow_deploy on Deploy {\n  rawID\n  state\n  commitSHA\n  commitMessage\n  startedAt\n}\n\nfragment RecentDeployments_viewer_1G22uz on Viewer {\n  recentDeploys(first: $count, after: $cursor) {\n    edges {\n      node {\n        id\n        ...DeploymentRow_deploy\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n"
  }
};
})();
(node as any).hash = '479b231eaed4af797148ebda73e8052a';
export default node;
