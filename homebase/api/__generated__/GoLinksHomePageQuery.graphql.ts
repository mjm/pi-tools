/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type GoLinksHomePageQueryVariables = {};
export type GoLinksHomePageQueryResponse = {
    readonly viewer: {
        readonly " $fragmentRefs": FragmentRefs<"RecentLinksList_viewer">;
    } | null;
};
export type GoLinksHomePageQuery = {
    readonly response: GoLinksHomePageQueryResponse;
    readonly variables: GoLinksHomePageQueryVariables;
};



/*
query GoLinksHomePageQuery {
  viewer {
    ...RecentLinksList_viewer
  }
}

fragment LinkRow_link on Link {
  id
  shortURL
  description
}

fragment RecentLinksList_viewer on Viewer {
  links(first: 30) {
    edges {
      node {
        id
        ...LinkRow_link
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
            "args": null,
            "kind": "FragmentSpread",
            "name": "RecentLinksList_viewer"
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
            "args": (v0/*: any*/),
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
            "storageKey": "links(first:30)"
          },
          {
            "alias": null,
            "args": (v0/*: any*/),
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
    "cacheID": "57a681eb021931db30bd3f85bf9ea90e",
    "id": null,
    "metadata": {},
    "name": "GoLinksHomePageQuery",
    "operationKind": "query",
    "text": "query GoLinksHomePageQuery {\n  viewer {\n    ...RecentLinksList_viewer\n  }\n}\n\nfragment LinkRow_link on Link {\n  id\n  shortURL\n  description\n}\n\nfragment RecentLinksList_viewer on Viewer {\n  links(first: 30) {\n    edges {\n      node {\n        id\n        ...LinkRow_link\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n"
  }
};
})();
(node as any).hash = '7c7f50c485742718196025c93cb618b5';
export default node;
