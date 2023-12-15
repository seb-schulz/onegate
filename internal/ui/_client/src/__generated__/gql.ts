/* eslint-disable */
import * as types from './graphql';
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';

/**
 * Map of all GraphQL operations in the project.
 *
 * This map has several performance disadvantages:
 * 1. It is not tree-shakeable, so it will include all operations in the project.
 * 2. It is not minifiable, so the string of a GraphQL query will be multiple times inside the bundle.
 * 3. It does not support dead code elimination, so it will add unused operations.
 *
 * Therefore it is highly recommended to use the babel or swc plugin for production.
 */
const documents = {
    "\nquery me {\n  me {\n    displayName\n    name\n    credentials {\n      id\n      createdAt\n      updatedAt\n    }\n  }\n}": types.MeDocument,
    "\nmutation beginLogin {\n  beginLogin\n}\n": types.BeginLoginDocument,
    "\nmutation validateLogin($body: CredentialRequestResponse!) {\n    validateLogin(body: $body)\n}\n": types.ValidateLoginDocument,
    "\nmutation createUser($name: String!) {\n  createUser(name: $name)\n}\n": types.CreateUserDocument,
    "\nmutation addPasskey($body: CredentialCreationResponse!) {\n  addPasskey(body: $body)\n}\n": types.AddPasskeyDocument,
};

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 *
 *
 * @example
 * ```ts
 * const query = gql(`query GetUser($id: ID!) { user(id: $id) { name } }`);
 * ```
 *
 * The query argument is unknown!
 * Please regenerate the types.
 */
export function gql(source: string): unknown;

/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nquery me {\n  me {\n    displayName\n    name\n    credentials {\n      id\n      createdAt\n      updatedAt\n    }\n  }\n}"): (typeof documents)["\nquery me {\n  me {\n    displayName\n    name\n    credentials {\n      id\n      createdAt\n      updatedAt\n    }\n  }\n}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation beginLogin {\n  beginLogin\n}\n"): (typeof documents)["\nmutation beginLogin {\n  beginLogin\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation validateLogin($body: CredentialRequestResponse!) {\n    validateLogin(body: $body)\n}\n"): (typeof documents)["\nmutation validateLogin($body: CredentialRequestResponse!) {\n    validateLogin(body: $body)\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation createUser($name: String!) {\n  createUser(name: $name)\n}\n"): (typeof documents)["\nmutation createUser($name: String!) {\n  createUser(name: $name)\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation addPasskey($body: CredentialCreationResponse!) {\n  addPasskey(body: $body)\n}\n"): (typeof documents)["\nmutation addPasskey($body: CredentialCreationResponse!) {\n  addPasskey(body: $body)\n}\n"];

export function gql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;