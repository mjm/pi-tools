/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type GoLinksHomePageQueryVariables = {};
export type GoLinksHomePageQueryResponse = {
    readonly viewer: {
        readonly links: {
            readonly __id: string;
            readonly edges: ReadonlyArray<{
                readonly __id: string;
            }>;
            readonly " $fragmentRefs": FragmentRefs<"RecentLinksList_links">;
        } | null;
    } | null;
};
export type GoLinksHomePageQuery = {
    readonly response: GoLinksHomePageQueryResponse;
    readonly variables: GoLinksHomePageQueryVariables;
};



/*
query GoLinksHomePageQuery {
  viewer {
    links(first: 30) {
      ...RecentLinksList_links
      edges {
        cursor
        node {
          __typename
          id
        }
      }
      pageInfo {
        endCursor
        hasNextPage
      }
    }
  }
}

fragment LinkRow_link on Link {
  id
  shortURL
  description
}

fragment RecentLinksList_links on LinkConnection {
  edges {
    node {
      id
      ...LinkRow_link
    }
  }
}
*/

const node: ConcreteRequest = (function(){
var v0 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "cursor",
  "storageKey": null
},
v1 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "__typename",
  "storageKey": null
},
v2 = {
  "kind": "ClientExtension",
  "selections": [
    {
      "alias": null,
      "args": null,
      "kind": "ScalarField",
      "name": "__id",
      "storageKey": null
    }
  ]
},
v3 = {
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
},
v4 = [
  {
    "kind": "Literal",
    "name": "first",
    "value": 30
  }
];
return {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "GoLinksHomePageQuery",
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
            "alias": "links",
            "args": null,
            "concreteType": "LinkConnection",
            "kind": "LinkedField",
            "name": "__RecentLinksList_links_connection",
            "plural": false,
            "selections": [
              {
                "alias": null,
                "args": null,
                "concreteType": "LinkEdge",
                "kind": "LinkedField",
                "name": "edges",
                "plural": true,
                "selections": [
                  (v0/*: any*/),
                  {
                    "alias": null,
                    "args": null,
                    "concreteType": "Link",
                    "kind": "LinkedField",
                    "name": "node",
                    "plural": false,
                    "selections": [
                      (v1/*: any*/)
                    ],
                    "storageKey": null
                  },
                  (v2/*: any*/)
                ],
                "storageKey": null
              },
              (v3/*: any*/),
              {
                "args": null,
                "kind": "FragmentSpread",
                "name": "RecentLinksList_links"
              },
              (v2/*: any*/)
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
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "GoLinksHomePageQuery",
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
            "args": (v4/*: any*/),
            "concreteType": "LinkConnection",
            "kind": "LinkedField",
            "name": "links",
            "plural": false,
            "selections": [
              {
                "alias": null,
                "args": null,
                "concreteType": "LinkEdge",
                "kind": "LinkedField",
                "name": "edges",
                "plural": true,
                "selections": [
                  {
                    "alias": null,
                    "args": null,
                    "concreteType": "Link",
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
                        "name": "shortURL",
                        "storageKey": null
                      },
                      {
                        "alias": null,
                        "args": null,
                        "kind": "ScalarField",
                        "name": "description",
                        "storageKey": null
                      },
                      (v1/*: any*/)
                    ],
                    "storageKey": null
                  },
                  (v0/*: any*/),
                  (v2/*: any*/)
                ],
                "storageKey": null
              },
              (v3/*: any*/),
              (v2/*: any*/)
            ],
            "storageKey": "links(first:30)"
          },
          {
            "alias": null,
            "args": (v4/*: any*/),
            "filters": null,
            "handle": "connection",
            "key": "RecentLinksList_links",
            "kind": "LinkedHandle",
            "name": "links"
          }
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "9126325604c99d803b054de19174cc62",
    "id": null,
    "metadata": {
      "connection": [
        {
          "count": null,
          "cursor": null,
          "direction": "forward",
          "path": [
            "viewer",
            "links"
          ]
        }
      ]
    },
    "name": "GoLinksHomePageQuery",
    "operationKind": "query",
    "text": "query GoLinksHomePageQuery {\n  viewer {\n    links(first: 30) {\n      ...RecentLinksList_links\n      edges {\n        cursor\n        node {\n          __typename\n          id\n        }\n      }\n      pageInfo {\n        endCursor\n        hasNextPage\n      }\n    }\n  }\n}\n\nfragment LinkRow_link on Link {\n  id\n  shortURL\n  description\n}\n\nfragment RecentLinksList_links on LinkConnection {\n  edges {\n    node {\n      id\n      ...LinkRow_link\n    }\n  }\n}\n"
  }
};
})();
(node as any).hash = '9d89da67cb1db10f8853f7d7f3d7e436';
export default node;
