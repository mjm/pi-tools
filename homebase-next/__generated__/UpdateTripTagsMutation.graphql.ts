/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
export type UpdateTripTagsInput = {
    tripID: string;
    tagsToAdd: Array<string>;
    tagsToRemove: Array<string>;
};
export type UpdateTripTagsMutationVariables = {
    input: UpdateTripTagsInput;
};
export type UpdateTripTagsMutationResponse = {
    readonly updateTripTags: {
        readonly trip: {
            readonly id: string;
            readonly tags: ReadonlyArray<string>;
        } | null;
    };
};
export type UpdateTripTagsMutation = {
    readonly response: UpdateTripTagsMutationResponse;
    readonly variables: UpdateTripTagsMutationVariables;
};



/*
mutation UpdateTripTagsMutation(
  $input: UpdateTripTagsInput!
) {
  updateTripTags(input: $input) {
    trip {
      id
      tags
    }
  }
}
*/

const node: ConcreteRequest = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "input"
  }
],
v1 = [
  {
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "input",
        "variableName": "input"
      }
    ],
    "concreteType": "UpdateTripTagsPayload",
    "kind": "LinkedField",
    "name": "updateTripTags",
    "plural": false,
    "selections": [
      {
        "alias": null,
        "args": null,
        "concreteType": "Trip",
        "kind": "LinkedField",
        "name": "trip",
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
            "name": "tags",
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
    "name": "UpdateTripTagsMutation",
    "selections": (v1/*: any*/),
    "type": "Mutation",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "UpdateTripTagsMutation",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "2ad3e39a10cdf4985170076323e67efd",
    "id": null,
    "metadata": {},
    "name": "UpdateTripTagsMutation",
    "operationKind": "mutation",
    "text": "mutation UpdateTripTagsMutation(\n  $input: UpdateTripTagsInput!\n) {\n  updateTripTags(input: $input) {\n    trip {\n      id\n      tags\n    }\n  }\n}\n"
  }
};
})();
(node as any).hash = 'ac1b907ff5ede12f2c61c1511836c4bb';
export default node;
