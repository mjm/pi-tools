/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
export type Id_DeployQueryVariables = {
    id: string;
};
export type Id_DeployQueryResponse = {
    readonly viewer: {
        readonly deploy: {
            readonly id: string;
            readonly commitSHA: string;
            readonly commitMessage: string;
            readonly startedAt: string;
            readonly finishedAt: string | null;
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
      commitSHA
      commitMessage
      startedAt
      finishedAt
    }
  }
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
            "kind": "Variable",
            "name": "id",
            "variableName": "id"
          }
        ],
        "concreteType": "Deploy",
        "kind": "LinkedField",
        "name": "deploy",
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
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "Id_DeployQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "Id_DeployQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "7f7d9786ee8fc10292c1166a38ab1b98",
    "id": null,
    "metadata": {},
    "name": "Id_DeployQuery",
    "operationKind": "query",
    "text": "query Id_DeployQuery(\n  $id: ID!\n) {\n  viewer {\n    deploy(id: $id) {\n      id\n      commitSHA\n      commitMessage\n      startedAt\n      finishedAt\n    }\n  }\n}\n"
  }
};
})();
(node as any).hash = '40d98e9501650721d16c9dba426e0ba2';
export default node;
