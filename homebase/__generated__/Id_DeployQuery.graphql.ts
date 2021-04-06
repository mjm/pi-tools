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
            readonly " $fragmentRefs": FragmentRefs<"DeploymentDetails_deploy" | "DeploymentEvents_deploy">;
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
      ...DeploymentEvents_deploy
    }
  }
}

fragment DeploymentDetails_deploy on Deploy {
  commitSHA
  commitMessage
  startedAt
  finishedAt
}

fragment DeploymentEvent_deploy on Deploy {
  startedAt
}

fragment DeploymentEvent_event on DeployEvent {
  timestamp
  level
  summary
  description
}

fragment DeploymentEvents_deploy on Deploy {
  ...DeploymentEvent_deploy
  report {
    events {
      ...DeploymentEvent_event
    }
    id
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
              },
              {
                "args": null,
                "kind": "FragmentSpread",
                "name": "DeploymentEvents_deploy"
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
              },
              {
                "alias": null,
                "args": null,
                "concreteType": "DeployReport",
                "kind": "LinkedField",
                "name": "report",
                "plural": false,
                "selections": [
                  {
                    "alias": null,
                    "args": null,
                    "concreteType": "DeployEvent",
                    "kind": "LinkedField",
                    "name": "events",
                    "plural": true,
                    "selections": [
                      {
                        "alias": null,
                        "args": null,
                        "kind": "ScalarField",
                        "name": "timestamp",
                        "storageKey": null
                      },
                      {
                        "alias": null,
                        "args": null,
                        "kind": "ScalarField",
                        "name": "level",
                        "storageKey": null
                      },
                      {
                        "alias": null,
                        "args": null,
                        "kind": "ScalarField",
                        "name": "summary",
                        "storageKey": null
                      },
                      {
                        "alias": null,
                        "args": null,
                        "kind": "ScalarField",
                        "name": "description",
                        "storageKey": null
                      }
                    ],
                    "storageKey": null
                  },
                  (v2/*: any*/)
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
    "cacheID": "44dd98bf0e22682df545502984586db6",
    "id": null,
    "metadata": {},
    "name": "Id_DeployQuery",
    "operationKind": "query",
    "text": "query Id_DeployQuery(\n  $id: ID!\n) {\n  viewer {\n    deploy(id: $id) {\n      id\n      ...DeploymentDetails_deploy\n      ...DeploymentEvents_deploy\n    }\n  }\n}\n\nfragment DeploymentDetails_deploy on Deploy {\n  commitSHA\n  commitMessage\n  startedAt\n  finishedAt\n}\n\nfragment DeploymentEvent_deploy on Deploy {\n  startedAt\n}\n\nfragment DeploymentEvent_event on DeployEvent {\n  timestamp\n  level\n  summary\n  description\n}\n\nfragment DeploymentEvents_deploy on Deploy {\n  ...DeploymentEvent_deploy\n  report {\n    events {\n      ...DeploymentEvent_event\n    }\n    id\n  }\n}\n"
  }
};
})();
(node as any).hash = 'dbd4c23533cc1a755e77ecf103e5f2b0';
export default node;
