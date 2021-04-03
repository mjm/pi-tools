/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
export type IgnoreTripInput = {
    id: string;
};
export type IgnoreTripMutationVariables = {
    input: IgnoreTripInput;
    connections: Array<string>;
};
export type IgnoreTripMutationResponse = {
    readonly ignoreTrip: {
        readonly ignoredTripID: string;
    };
};
export type IgnoreTripMutation = {
    readonly response: IgnoreTripMutationResponse;
    readonly variables: IgnoreTripMutationVariables;
};



/*
mutation IgnoreTripMutation(
  $input: IgnoreTripInput!
) {
  ignoreTrip(input: $input) {
    ignoredTripID
  }
}
*/

const node: ConcreteRequest = (function(){
var v0 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "connections"
},
v1 = {
  "defaultValue": null,
  "kind": "LocalArgument",
  "name": "input"
},
v2 = [
  {
    "kind": "Variable",
    "name": "input",
    "variableName": "input"
  }
],
v3 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "ignoredTripID",
  "storageKey": null
};
return {
  "fragment": {
    "argumentDefinitions": [
      (v0/*: any*/),
      (v1/*: any*/)
    ],
    "kind": "Fragment",
    "metadata": null,
    "name": "IgnoreTripMutation",
    "selections": [
      {
        "alias": null,
        "args": (v2/*: any*/),
        "concreteType": "IgnoreTripPayload",
        "kind": "LinkedField",
        "name": "ignoreTrip",
        "plural": false,
        "selections": [
          (v3/*: any*/)
        ],
        "storageKey": null
      }
    ],
    "type": "Mutation",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": [
      (v1/*: any*/),
      (v0/*: any*/)
    ],
    "kind": "Operation",
    "name": "IgnoreTripMutation",
    "selections": [
      {
        "alias": null,
        "args": (v2/*: any*/),
        "concreteType": "IgnoreTripPayload",
        "kind": "LinkedField",
        "name": "ignoreTrip",
        "plural": false,
        "selections": [
          (v3/*: any*/),
          {
            "alias": null,
            "args": null,
            "filters": null,
            "handle": "deleteEdge",
            "key": "",
            "kind": "ScalarHandle",
            "name": "ignoredTripID",
            "handleArgs": [
              {
                "kind": "Variable",
                "name": "connections",
                "variableName": "connections"
              }
            ]
          }
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "3f3ee2e20ab97bc064c3624b11b2f85d",
    "id": null,
    "metadata": {},
    "name": "IgnoreTripMutation",
    "operationKind": "mutation",
    "text": "mutation IgnoreTripMutation(\n  $input: IgnoreTripInput!\n) {\n  ignoreTrip(input: $input) {\n    ignoredTripID\n  }\n}\n"
  }
};
})();
(node as any).hash = '6f54b85368bc7f017f3ffd2a8f18c5fb';
export default node;
