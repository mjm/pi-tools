/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type CreateLinkInput = {
    shortURL: string;
    destinationURL: string;
    description: string;
};
export type CreateLinkMutationVariables = {
    input: CreateLinkInput;
    connections: Array<string>;
};
export type CreateLinkMutationResponse = {
    readonly createLink: {
        readonly link: {
            readonly id: string;
            readonly " $fragmentRefs": FragmentRefs<"LinkRow_link">;
        };
    };
};
export type CreateLinkMutation = {
    readonly response: CreateLinkMutationResponse;
    readonly variables: CreateLinkMutationVariables;
};



/*
mutation CreateLinkMutation(
  $input: CreateLinkInput!
) {
  createLink(input: $input) {
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
  "name": "id",
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
    "name": "CreateLinkMutation",
    "selections": [
      {
        "alias": null,
        "args": (v2/*: any*/),
        "concreteType": "CreateLinkPayload",
        "kind": "LinkedField",
        "name": "createLink",
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
              (v3/*: any*/),
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
    "argumentDefinitions": [
      (v1/*: any*/),
      (v0/*: any*/)
    ],
    "kind": "Operation",
    "name": "CreateLinkMutation",
    "selections": [
      {
        "alias": null,
        "args": (v2/*: any*/),
        "concreteType": "CreateLinkPayload",
        "kind": "LinkedField",
        "name": "createLink",
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
              (v3/*: any*/),
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
          },
          {
            "alias": null,
            "args": null,
            "filters": null,
            "handle": "prependNode",
            "key": "",
            "kind": "LinkedHandle",
            "name": "link",
            "handleArgs": [
              {
                "kind": "Variable",
                "name": "connections",
                "variableName": "connections"
              },
              {
                "kind": "Literal",
                "name": "edgeTypeName",
                "value": "LinkEdge"
              }
            ]
          }
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "7cdec9ebf11948bf5c3be8b3e6e8e410",
    "id": null,
    "metadata": {},
    "name": "CreateLinkMutation",
    "operationKind": "mutation",
    "text": "mutation CreateLinkMutation(\n  $input: CreateLinkInput!\n) {\n  createLink(input: $input) {\n    link {\n      id\n      ...LinkRow_link\n    }\n  }\n}\n\nfragment LinkRow_link on Link {\n  id\n  shortURL\n  description\n}\n"
  }
};
})();
(node as any).hash = '23b6e7e5203801e0f05ee222cc77e939';
export default node;
