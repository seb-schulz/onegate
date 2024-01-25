import React, { useState, useRef } from "react";
import { Button, Card, Form } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import * as urql from 'urql';
import { startRegistration } from '@simplewebauthn/browser';
import { gql } from "../__generated__/gql";

const CREATE_USER_GQL = gql(`
mutation createUser($name: String!) {
  createUser(name: $name)
}
`)

const ADD_PASSKEY_QGL = gql(`
mutation addCredential($body: CredentialCreationResponse!) {
    addCredential(body: $body)
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
    const [{ fetching: loadingCreateUser }, createUser] = urql.useMutation(CREATE_USER_GQL);
    const [{ fetching: loadingAddPasskey }, addCredential] = urql.useMutation(ADD_PASSKEY_QGL);

    if (loadingCreateUser || loadingAddPasskey) return <p>Loading...</p>;

    const handleSubmit = async (e: React.SyntheticEvent) => {
        e.preventDefault();
        e.stopPropagation();
        setValidated(true);

        if (userNameRef.current === null) return;
        if (!userNameRef.current.value) return;

        const userName = userNameRef.current.value;
        try {
            const result = await createUser({ name: userName });

            if (!result.data || !result.data.createUser) {
                onError("cannot load data");
                return;
            }
            setTimeout(onUserCreated, 0);

            const attResp = await startRegistration(result.data.createUser.publicKey);

            const addPasskeyData = await addCredential({ body: JSON.stringify(attResp) });

            if (!addPasskeyData?.data?.addCredential) {
                onError("cannot load data");
                return;
            }
            setTimeout(onPasskeyAdded, 0);
        } catch (error) {
            onError((error as urql.CombinedError).message)
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
