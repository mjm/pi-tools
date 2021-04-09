/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type BackupsListPaginationQueryVariables = {
    count?: number | null;
    cursor?: string | null;
};
export type BackupsListPaginationQueryResponse = {
    readonly viewer: {
        readonly " $fragmentRefs": FragmentRefs<"BackupsList_viewer">;
    } | null;
};
export type BackupsListPaginationQuery = {
    readonly response: BackupsListPaginationQueryResponse;
    readonly variables: BackupsListPaginationQueryVariables;
};



/*
query BackupsListPaginationQuery(
  $count: Int = 10
  $cursor: Cursor
) {
  viewer {
    ...BackupsList_viewer_1G22uz
  }
}

fragment ArchiveRow_archive on Archive {
  id
  name
  createdAt
}

fragment BackupsList_viewer_1G22uz on Viewer {
  backupArchives(first: $count, after: $cursor) {
    edges {
      node {
        id
        ...ArchiveRow_archive
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
    "defaultValue": 10,
    "kind": "LocalArgument",
    "name": "count"
  },
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "cursor"
  }
],
v1 = [
  {
    "kind": "Variable",
    "name": "after",
    "variableName": "cursor"
  },
  {
    "kind": "Variable",
    "name": "first",
    "variableName": "count"
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "BackupsListPaginationQuery",
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
            "args": [
              {
                "kind": "Variable",
                "name": "count",
                "variableName": "count"
              },
              {
                "kind": "Variable",
                "name": "cursor",
                "variableName": "cursor"
              }
            ],
            "kind": "FragmentSpread",
            "name": "BackupsList_viewer"
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
    "name": "BackupsListPaginationQuery",
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
            "concreteType": "ArchiveConnection",
            "kind": "LinkedField",
            "name": "backupArchives",
            "plural": false,
            "selections": [
              {
                "alias": null,
                "args": null,
                "concreteType": "ArchiveEdge",
                "kind": "LinkedField",
                "name": "edges",
                "plural": true,
                "selections": [
                  {
                    "alias": null,
                    "args": null,
                    "concreteType": "Archive",
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
                        "name": "name",
                        "storageKey": null
                      },
                      {
                        "alias": null,
                        "args": null,
                        "kind": "ScalarField",
                        "name": "createdAt",
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
            "args": (v1/*: any*/),
            "filters": [
              "kind"
            ],
            "handle": "connection",
            "key": "BackupsList_backupArchives",
            "kind": "LinkedHandle",
            "name": "backupArchives"
          }
        ],
        "storageKey": null
      }
    ]
  },
  "params": {
    "cacheID": "57671810b2b95aac316d541f96a98758",
    "id": null,
    "metadata": {},
    "name": "BackupsListPaginationQuery",
    "operationKind": "query",
    "text": "query BackupsListPaginationQuery(\n  $count: Int = 10\n  $cursor: Cursor\n) {\n  viewer {\n    ...BackupsList_viewer_1G22uz\n  }\n}\n\nfragment ArchiveRow_archive on Archive {\n  id\n  name\n  createdAt\n}\n\nfragment BackupsList_viewer_1G22uz on Viewer {\n  backupArchives(first: $count, after: $cursor) {\n    edges {\n      node {\n        id\n        ...ArchiveRow_archive\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n"
  }
};
})();
(node as any).hash = '5574c804094681e85ccb9436565317d6';
export default node;
