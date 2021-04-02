/* tslint:disable */
/* eslint-disable */
// @ts-nocheck

import { ConcreteRequest } from "relay-runtime";
import { FragmentRefs } from "relay-runtime";
export type backups_BackupsPageQueryVariables = {};
export type backups_BackupsPageQueryResponse = {
    readonly viewer: {
        readonly " $fragmentRefs": FragmentRefs<"BackupsList_viewer">;
    } | null;
};
export type backups_BackupsPageQuery = {
    readonly response: backups_BackupsPageQueryResponse;
    readonly variables: backups_BackupsPageQueryVariables;
};



/*
query backups_BackupsPageQuery {
  viewer {
    ...BackupsList_viewer
  }
}

fragment ArchiveRow_archive on Archive {
  id
  name
  createdAt
}

fragment BackupsList_viewer on Viewer {
  backupArchives(first: 10) {
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
    "kind": "Literal",
    "name": "first",
    "value": 10
  }
];
return {
  "fragment": {
    "argumentDefinitions": [],
    "kind": "Fragment",
    "metadata": null,
    "name": "backups_BackupsPageQuery",
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
    "argumentDefinitions": [],
    "kind": "Operation",
    "name": "backups_BackupsPageQuery",
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
            "storageKey": "backupArchives(first:10)"
          },
          {
            "alias": null,
            "args": (v0/*: any*/),
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
    "cacheID": "4f22260262b5662d1dc5e16bd5395979",
    "id": null,
    "metadata": {},
    "name": "backups_BackupsPageQuery",
    "operationKind": "query",
    "text": "query backups_BackupsPageQuery {\n  viewer {\n    ...BackupsList_viewer\n  }\n}\n\nfragment ArchiveRow_archive on Archive {\n  id\n  name\n  createdAt\n}\n\nfragment BackupsList_viewer on Viewer {\n  backupArchives(first: 10) {\n    edges {\n      node {\n        id\n        ...ArchiveRow_archive\n        __typename\n      }\n      cursor\n    }\n    pageInfo {\n      endCursor\n      hasNextPage\n    }\n  }\n}\n"
  }
};
})();
(node as any).hash = '6d08309b81a2a4fad9d83dc6d50e1aee';
export default node;
