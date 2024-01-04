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
    "\nmutation beginLogin {\n  beginLogin\n}\n": types.BeginLoginDocument,
    "\nmutation validateLogin($body: CredentialRequestResponse!) {\n    validateLogin(body: $body)\n}\n": types.ValidateLoginDocument,
    "\nmutation createUser($name: String!) {\n  createUser(name: $name)\n}\n": types.CreateUserDocument,
    "\nmutation addCredential($body: CredentialCreationResponse!) {\n    addCredential(body: $body)\n}\n": types.AddCredentialDocument,
    "\nquery credentials {\n  credentials {\n    id\n    description\n    createdAt\n    updatedAt\n    lastLogin\n  }\n}": types.CredentialsDocument,
    "\nmutation updateCredential($id: ID!, $description: String) {\n  updateCredential(id: $id, description: $description) {\n    id\n  }\n}": types.UpdateCredentialDocument,
    "\nmutation initCredential {\n    initCredential\n}\n": types.InitCredentialDocument,
    "\nmutation removeCredential($id: ID!) {\n    removeCredential(id: $id)\n}\n": types.RemoveCredentialDocument,
    "\nmutation updateMe($name: String, $displayName: String) {\n  updateMe(name: $name, displayName: $displayName) {\n    name\n  }\n}\n": types.UpdateMeDocument,
    "\nquery me {\n  me {\n    displayName\n    name\n  }\n}": types.MeDocument,
    "\nquery sessions {\n  sessions {\n    id\n    createdAt\n    updatedAt\n    isActive\n    isCurrent\n  }\n}\n": types.SessionsDocument,
    "\nmutation removeSession($id: ID!) {\n  removeSession(id: $id)\n}\n": types.RemoveSessionDocument,
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
export function gql(source: "\nmutation addCredential($body: CredentialCreationResponse!) {\n    addCredential(body: $body)\n}\n"): (typeof documents)["\nmutation addCredential($body: CredentialCreationResponse!) {\n    addCredential(body: $body)\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nquery credentials {\n  credentials {\n    id\n    description\n    createdAt\n    updatedAt\n    lastLogin\n  }\n}"): (typeof documents)["\nquery credentials {\n  credentials {\n    id\n    description\n    createdAt\n    updatedAt\n    lastLogin\n  }\n}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation updateCredential($id: ID!, $description: String) {\n  updateCredential(id: $id, description: $description) {\n    id\n  }\n}"): (typeof documents)["\nmutation updateCredential($id: ID!, $description: String) {\n  updateCredential(id: $id, description: $description) {\n    id\n  }\n}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation initCredential {\n    initCredential\n}\n"): (typeof documents)["\nmutation initCredential {\n    initCredential\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation removeCredential($id: ID!) {\n    removeCredential(id: $id)\n}\n"): (typeof documents)["\nmutation removeCredential($id: ID!) {\n    removeCredential(id: $id)\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation updateMe($name: String, $displayName: String) {\n  updateMe(name: $name, displayName: $displayName) {\n    name\n  }\n}\n"): (typeof documents)["\nmutation updateMe($name: String, $displayName: String) {\n  updateMe(name: $name, displayName: $displayName) {\n    name\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nquery me {\n  me {\n    displayName\n    name\n  }\n}"): (typeof documents)["\nquery me {\n  me {\n    displayName\n    name\n  }\n}"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nquery sessions {\n  sessions {\n    id\n    createdAt\n    updatedAt\n    isActive\n    isCurrent\n  }\n}\n"): (typeof documents)["\nquery sessions {\n  sessions {\n    id\n    createdAt\n    updatedAt\n    isActive\n    isCurrent\n  }\n}\n"];
/**
 * The gql function is used to parse GraphQL queries into a document that can be used by GraphQL clients.
 */
export function gql(source: "\nmutation removeSession($id: ID!) {\n  removeSession(id: $id)\n}\n"): (typeof documents)["\nmutation removeSession($id: ID!) {\n  removeSession(id: $id)\n}\n"];

export function gql(source: string) {
  return (documents as any)[source] ?? {};
}

export type DocumentType<TDocumentNode extends DocumentNode<any, any>> = TDocumentNode extends DocumentNode<  infer TType,  any>  ? TType  : never;