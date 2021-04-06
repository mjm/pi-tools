/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type Id_DeployQueryVariables = {
    id: string;
};
export type Id_DeployQueryResponse = {
    readonly viewer: {
        readonly deploy: {
            readonly id: string;
            readonly " $fragmentRefs": FragmentRefs<"DeploymentDetails_deploy">;
        } | null;
    } | null;
};
export type Id_DeployQuery = {
    readonly response: Id_DeployQueryResponse;
    readonly variables: Id_DeployQueryVariables;
};



/*
query Id_DeployQuery(
  $id: ID!
) {
  viewer {
    deploy(id: $id) {
      id
      ...DeploymentDetails_deploy
    }
  }
}

fragment DeploymentDetails_deploy on Deploy {
  commitSHA
  commitMessage
  startedAt
  finishedAt
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
};
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "Id_DeployQuery",
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
            "concreteType": "Deploy",
            "kind": "LinkedField",
            "name": "deploy",
            "plural": false,
            "selections": [
              (v2/*: any*/),
              {
                "args": null,
                "kind": "FragmentSpread",
                "name": "DeploymentDetails_deploy"
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
    "name": "Id_DeployQuery",
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
            "concreteType": "Deploy",
            "kind": "LinkedField",
            "name": "deploy",
            "plural": false,
            "selections": [
              (v2/*: any*/),
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
                "name": "finishedAt",
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
    "cacheID": "6dd00a3327ee8d00ea65141df0bd669f",
    "id": null,
    "metadata": {},
    "name": "Id_DeployQuery",
    "operationKind": "query",
    "text": "query Id_DeployQuery(\n  $id: ID!\n) {\n  viewer {\n    deploy(id: $id) {\n      id\n      ...DeploymentDetails_deploy\n    }\n  }\n}\n\nfragment DeploymentDetails_deploy on Deploy {\n  commitSHA\n  commitMessage\n  startedAt\n  finishedAt\n}\n"
  }
};
})();
(node as any).hash = 'd31c0ed5360d76660960db6325eb8ba3';
export default node;
