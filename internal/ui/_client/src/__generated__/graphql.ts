/* eslint-disable */
import { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  CredentialAssertion: { input: any; output: any; }
  CredentialCreation: { input: any; output: any; }
  CredentialCreationResponse: { input: any; output: any; }
  CredentialRequestResponse: { input: any; output: any; }
  Time: { input: any; output: any; }
};

export type Credential = {
  __typename?: 'Credential';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  addPasskey: Scalars['Boolean']['output'];
  beginLogin: Scalars['CredentialAssertion']['output'];
  createUser: Scalars['CredentialCreation']['output'];
  validateLogin: Scalars['Boolean']['output'];
};


export type MutationAddPasskeyArgs = {
  body: Scalars['CredentialCreationResponse']['input'];
};


export type MutationCreateUserArgs = {
  name: Scalars['String']['input'];
};


export type MutationValidateLoginArgs = {
  body: Scalars['CredentialRequestResponse']['input'];
};

export type PubKeyCredParam = {
  __typename?: 'PubKeyCredParam';
  alg: Scalars['Int']['output'];
  type: Scalars['String']['output'];
};

export type Query = {
  __typename?: 'Query';
  me?: Maybe<User>;
};

export type RelyingParty = {
  __typename?: 'RelyingParty';
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type User = {
  __typename?: 'User';
  credentials: Array<Maybe<Credential>>;
  displayName: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type MeQueryVariables = Exact<{ [key: string]: never; }>;


export type MeQuery = { __typename?: 'Query', me?: { __typename?: 'User', displayName: string, name: string, credentials: Array<{ __typename?: 'Credential', id: string, createdAt: any, updatedAt: any } | null> } | null };

export type BeginLoginMutationVariables = Exact<{ [key: string]: never; }>;


export type BeginLoginMutation = { __typename?: 'Mutation', beginLogin: any };

export type ValidateLoginMutationVariables = Exact<{
  body: Scalars['CredentialRequestResponse']['input'];
}>;


export type ValidateLoginMutation = { __typename?: 'Mutation', validateLogin: boolean };

export type CreateUserMutationVariables = Exact<{
  name: Scalars['String']['input'];
}>;


export type CreateUserMutation = { __typename?: 'Mutation', createUser: any };

export type AddPasskeyMutationVariables = Exact<{
  body: Scalars['CredentialCreationResponse']['input'];
}>;


export type AddPasskeyMutation = { __typename?: 'Mutation', addPasskey: boolean };


export const MeDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"name"}},{"kind":"Field","name":{"kind":"Name","value":"credentials"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}}]}}]} as unknown as DocumentNode<MeQuery, MeQueryVariables>;
export const BeginLoginDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"beginLogin"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"beginLogin"}}]}}]} as unknown as DocumentNode<BeginLoginMutation, BeginLoginMutationVariables>;
export const ValidateLoginDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"validateLogin"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"body"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CredentialRequestResponse"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"validateLogin"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"body"},"value":{"kind":"Variable","name":{"kind":"Name","value":"body"}}}]}]}}]} as unknown as DocumentNode<ValidateLoginMutation, ValidateLoginMutationVariables>;
export const CreateUserDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"createUser"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"name"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createUser"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"name"},"value":{"kind":"Variable","name":{"kind":"Name","value":"name"}}}]}]}}]} as unknown as DocumentNode<CreateUserMutation, CreateUserMutationVariables>;
export const AddPasskeyDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"addPasskey"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"body"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CredentialCreationResponse"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"addPasskey"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"body"},"value":{"kind":"Variable","name":{"kind":"Name","value":"body"}}}]}]}}]} as unknown as DocumentNode<AddPasskeyMutation, AddPasskeyMutationVariables>;