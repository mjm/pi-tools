/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type deploys_DeploysQueryVariables = {};
export type deploys_DeploysQueryResponse = {
    readonly viewer: {
        readonly " $fragmentRefs": FragmentRefs<"RecentDeployments_viewer">;
    } | null;
};
export type deploys_DeploysQuery = {
    readonly response: deploys_DeploysQueryResponse;
    readonly variables: deploys_DeploysQueryVariables;
};



/*
query deploys_DeploysQuery {
  viewer {
    ...RecentDeployments_viewer
  }
}

fragment DeploymentRow_deploy on Deploy {
  rawID
  state
  commitSHA
  commitMessage
  startedAt
}

fragment RecentDeployments_viewer on Viewer {
  recentDeploys {
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

const node: ConcreteRequest = {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "deploys_DeploysQuery",
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
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "deploys_DeploysQuery",
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
            "args": null,
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
    "cacheID": "20434589afc45bc6f81bbdb0b32d096e",
    "id": null,
    "metadata": {},
    "name": "deploys_DeploysQuery",
    "operationKind": "query",
    "text": "query deploys_DeploysQuery {\n  viewer {\n    ...RecentDeployments_viewer\n  }\n}\n\nfragment DeploymentRow_deploy on Deploy {\n  rawID\n  state\n  commitSHA\n  commitMessage\n  startedAt\n}\n\nfragment RecentDeployments_viewer on Viewer {\n  recentDeploys {\n    edges {\n      node {\n        id\n        ...DeploymentRow_deploy\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n"
  }
};
(node as any).hash = '844db39a40ebf3c70b3bbe8c8cccc218';
export default node;
