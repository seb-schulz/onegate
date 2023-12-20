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
  description: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  updatedAt: Scalars['Time']['output'];
};

export type Mutation = {
  __typename?: 'Mutation';
  addCredential: Scalars['Boolean']['output'];
  beginLogin: Scalars['CredentialAssertion']['output'];
  createUser: Scalars['CredentialCreation']['output'];
  removeCredential: Scalars['Boolean']['output'];
  updateCredential: Credential;
  validateLogin: Scalars['Boolean']['output'];
};


export type MutationAddCredentialArgs = {
  body: Scalars['CredentialCreationResponse']['input'];
};


export type MutationCreateUserArgs = {
  name: Scalars['String']['input'];
};


export type MutationRemoveCredentialArgs = {
  id: Scalars['ID']['input'];
};


export type MutationUpdateCredentialArgs = {
  description?: InputMaybe<Scalars['String']['input']>;
  id: Scalars['ID']['input'];
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
  credentials?: Maybe<Array<Maybe<Credential>>>;
  me?: Maybe<User>;
};

export type RelyingParty = {
  __typename?: 'RelyingParty';
  id: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type User = {
  __typename?: 'User';
  displayName: Scalars['String']['output'];
  name: Scalars['String']['output'];
};

export type BeginLoginMutationVariables = Exact<{ [key: string]: never; }>;


export type BeginLoginMutation = { __typename?: 'Mutation', beginLogin: any };

export type ValidateLoginMutationVariables = Exact<{
  body: Scalars['CredentialRequestResponse']['input'];
}>;


export type ValidateLoginMutation = { __typename?: 'Mutation', validateLogin: boolean };

export type MeNavbarQueryVariables = Exact<{ [key: string]: never; }>;


export type MeNavbarQuery = { __typename?: 'Query', me?: { __typename?: 'User', displayName: string, name: string } | null };

export type CreateUserMutationVariables = Exact<{
  name: Scalars['String']['input'];
}>;


export type CreateUserMutation = { __typename?: 'Mutation', createUser: any };

export type AddCredentialMutationVariables = Exact<{
  body: Scalars['CredentialCreationResponse']['input'];
}>;


export type AddCredentialMutation = { __typename?: 'Mutation', addCredential: boolean };

export type CredentialsQueryVariables = Exact<{ [key: string]: never; }>;


export type CredentialsQuery = { __typename?: 'Query', credentials?: Array<{ __typename?: 'Credential', id: string, description: string, createdAt: any, updatedAt: any } | null> | null };

export type UpdateCredentialMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  description?: InputMaybe<Scalars['String']['input']>;
}>;


export type UpdateCredentialMutation = { __typename?: 'Mutation', updateCredential: { __typename?: 'Credential', id: string } };

export type MeIndexQueryVariables = Exact<{ [key: string]: never; }>;


export type MeIndexQuery = { __typename?: 'Query', me?: { __typename?: 'User', displayName: string, name: string } | null };


export const BeginLoginDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"beginLogin"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"beginLogin"}}]}}]} as unknown as DocumentNode<BeginLoginMutation, BeginLoginMutationVariables>;
export const ValidateLoginDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"validateLogin"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"body"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CredentialRequestResponse"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"validateLogin"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"body"},"value":{"kind":"Variable","name":{"kind":"Name","value":"body"}}}]}]}}]} as unknown as DocumentNode<ValidateLoginMutation, ValidateLoginMutationVariables>;
export const MeNavbarDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"meNavbar"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}}]} as unknown as DocumentNode<MeNavbarQuery, MeNavbarQueryVariables>;
export const CreateUserDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"createUser"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"name"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"createUser"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"name"},"value":{"kind":"Variable","name":{"kind":"Name","value":"name"}}}]}]}}]} as unknown as DocumentNode<CreateUserMutation, CreateUserMutationVariables>;
export const AddCredentialDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"addCredential"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"body"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"CredentialCreationResponse"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"addCredential"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"body"},"value":{"kind":"Variable","name":{"kind":"Name","value":"body"}}}]}]}}]} as unknown as DocumentNode<AddCredentialMutation, AddCredentialMutationVariables>;
export const CredentialsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"credentials"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"credentials"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"description"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}},{"kind":"Field","name":{"kind":"Name","value":"updatedAt"}}]}}]}}]} as unknown as DocumentNode<CredentialsQuery, CredentialsQueryVariables>;
export const UpdateCredentialDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"updateCredential"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"description"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"updateCredential"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"description"},"value":{"kind":"Variable","name":{"kind":"Name","value":"description"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}}]}}]}}]} as unknown as DocumentNode<UpdateCredentialMutation, UpdateCredentialMutationVariables>;
export const MeIndexDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"meIndex"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"me"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"name"}}]}}]}}]} as unknown as DocumentNode<MeIndexQuery, MeIndexQueryVariables>;