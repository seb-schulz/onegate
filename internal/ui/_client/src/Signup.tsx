import React, { useState, useRef } from "react";
import { Button, Card, Form } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import { ApolloError, useMutation } from "@apollo/client";
import { startRegistration } from '@simplewebauthn/browser';
import { gql } from "./__generated__/gql";

const CREATE_USER_GQL = gql(`
mutation createUser($name: String!) {
  createUser(name: $name)
}
`)

const ADD_PASSKEY_QGL = gql(`
mutation addPasskey($body: CredentialCreationResponse!) {
  addPasskey(body: $body)
}
`)

function SignupCard({ onUserCreated, onError, onPasskeyAdded }: {
    onUserCreated: () => void
    onError: (errMsg: string) => void
    onPasskeyAdded: () => void
}) {
    const { t } = useTranslation();
    const [validated, setValidated] = useState(false);
    const userNameRef = useRef<HTMLInputElement | null>(null);
    const [createUser, { loading: loadingCreateUser }] = useMutation(CREATE_USER_GQL);
    const [addPasskey, { loading: loadingAddPasskey }] = useMutation(ADD_PASSKEY_QGL);

    if (loadingCreateUser || loadingAddPasskey) return <p>Loading...</p>;

    const handleSubmit = async (e: React.SyntheticEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setValidated(true);

        if (userNameRef.current === null) return;
        if (!userNameRef.current.value) return;

        const userName = userNameRef.current.value;

        const result = await createUser({
            variables: {
                name: userName,
            },
            onError: (error) => {
                onError((error as ApolloError).message)
            }
        });

        if (!result.data || !result.data.createUser) {
            onError("cannot load data");
            return;
        }
        setTimeout(onUserCreated, 0);

        try {
            const attResp = await startRegistration(result.data.createUser.publicKey);

            const addPasskeyData = await addPasskey({
                variables: {
                    body: JSON.stringify(attResp),
                },
            });

            if (!addPasskeyData.data || !addPasskeyData.data.addPasskey) {
                onError("cannot load data");
                return;
            }
            setTimeout(onPasskeyAdded, 0);
        } catch (error) {
            onError((error as ApolloError).message);
        }

    };

    return (
        <Card>
            <Card.Body>
                <Form noValidate validated={validated} onSubmit={handleSubmit}>
                    <Card.Text>
                        <Form.Label htmlFor="inputUserName">{t('Username')}</Form.Label>
                        <Form.Control required type="text" id="inputUserName" ref={userNameRef} autoComplete="username webauthn" />
                    </Card.Text>
                    <Button type="submit" disabled={!window.PublicKeyCredential || loadingCreateUser || loadingAddPasskey}>{t('Register')}</Button>{' '}
                </Form>
            </Card.Body>
        </Card>
    );
}

export default SignupCard;
