/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type UpdateLinkInput = {
    id: string;
    shortURL: string;
    destinationURL: string;
    description: string;
};
export type UpdateLinkMutationVariables = {
    input: UpdateLinkInput;
};
export type UpdateLinkMutationResponse = {
    readonly updateLink: {
        readonly link: {
            readonly id: string;
            readonly " $fragmentRefs": FragmentRefs<"LinkRow_link">;
        };
    };
};
export type UpdateLinkMutation = {
    readonly response: UpdateLinkMutationResponse;
    readonly variables: UpdateLinkMutationVariables;
};



/*
mutation UpdateLinkMutation(
  $input: UpdateLinkInput!
) {
  updateLink(input: $input) {
    link {
      id
      ...LinkRow_link
    }
  }
}

fragment LinkRow_link on Link {
  id
  shortURL
  description
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
    "kind": "Variable",
    "name": "input",
    "variableName": "input"
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
    "name": "UpdateLinkMutation",
    "selections": [
      {
        "alias": null,
        "args": (v1/*: any*/),
        "concreteType": "UpdateLinkPayload",
        "kind": "LinkedField",
        "name": "updateLink",
        "plural": false,
        "selections": [
          {
            "alias": null,
            "args": null,
            "concreteType": "Link",
            "kind": "LinkedField",
            "name": "link",
            "plural": false,
            "selections": [
              (v2/*: any*/),
              {
                "args": null,
                "kind": "FragmentSpread",
                "name": "LinkRow_link"
              }
            ],
            "storageKey": null
          }
        ],
        "storageKey": null
      }
    ],
    "type": "Mutation",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "UpdateLinkMutation",
    "selections": [
      {
        "alias": null,
        "args": (v1/*: any*/),
        "concreteType": "UpdateLinkPayload",
        "kind": "LinkedField",
        "name": "updateLink",
        "plural": false,
        "selections": [
          {
            "alias": null,
            "args": null,
            "concreteType": "Link",
            "kind": "LinkedField",
            "name": "link",
            "plural": false,
            "selections": [
              (v2/*: any*/),
              {
                "alias": null,
                "args": null,
                "kind": "ScalarField",
                "name": "shortURL",
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
          }
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "c832435a8aa5484a84815a9784aa696f",
    "id": null,
    "metadata": {},
    "name": "UpdateLinkMutation",
    "operationKind": "mutation",
    "text": "mutation UpdateLinkMutation(\n  $input: UpdateLinkInput!\n) {\n  updateLink(input: $input) {\n    link {\n      id\n      ...LinkRow_link\n    }\n  }\n}\n\nfragment LinkRow_link on Link {\n  id\n  shortURL\n  description\n}\n"
  }
};
})();
(node as any).hash = '6202dcc471a169ced48632904548fb02';
export default node;
