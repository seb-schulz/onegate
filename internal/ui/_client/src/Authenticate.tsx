import React, { useState, useRef } from "react";
import { Alert, Button, Card, Form } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import { gql, useMutation } from "@apollo/client";
import { startRegistration } from '@simplewebauthn/browser';

const CREATE_USER_GQL = gql`
mutation createUser($name: String!) {
  createUser(name: $name)
}
`

const ADD_PASSKEY_QGL = gql`
mutation addPasskey($body: CredentialCreationResponse!) {
  addPasskey(body: $body)
}
`

function AuthenticateCard({ loginSucceeded }: {
    loginSucceeded: () => void
}) {
    const { t } = useTranslation();
    const [validated, setValidated] = useState(false);
    const [errorMsg, setErrorMsg] = useState<string>("");
    const userNameRef = useRef<HTMLInputElement | null>(null);
    const [createUser, { loading: loadingCreateUser, error: errorCreateUser }] = useMutation(CREATE_USER_GQL);
    const [addPasskey, { loading: loadingAddPasskey, error: errorAddPasskey }] = useMutation(ADD_PASSKEY_QGL);

    const hasWebAuthN = !!window.PublicKeyCredential;
    const errorMsgList = [];

    if (!hasWebAuthN) errorMsgList.push(t('This browser does not support WebAuthN.'));
    if (!!errorMsg) errorMsgList.push(errorMsg);
    if (loadingCreateUser) return <p>Loading...</p>;
    if (errorCreateUser) errorMsgList.push(errorCreateUser.message);
    if (loadingAddPasskey) return <p>Loading...</p>;
    if (errorAddPasskey) errorMsgList.push(errorAddPasskey.message);

    const handleSubmit = (event: React.SyntheticEvent) => {
        event.preventDefault();
        event.stopPropagation();
        // console.log("login", userName)
        setValidated(true);
    };

    const handleRegistration = async (e: React.SyntheticEvent) => {
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
        });
        setTimeout(loginSucceeded, 0);

        let attResp;
        try {
            attResp = await startRegistration(result.data.createUser.publicKey);
        } catch (error) {
            setErrorMsg(error as string);
            throw error;
        }

        try {
            const addPasskeyData = await addPasskey({
                variables: {
                    body: JSON.stringify(attResp),
                },
            });

            if (addPasskeyData.data.addPasskey) {
                // TODO: Add alert on top
            }
        } catch (error) {
            setErrorMsg(error as string);
            throw error;
        }

    };

    return (
        <Card>
            <Card.Body>
                {errorMsgList.length > 0 ? <Alert variant="danger">{errorMsgList.join(',`')}</Alert> : ''}
                <Form noValidate validated={validated} onSubmit={handleSubmit}>

                    <Card.Text>
                        <Form.Control required type="text" id="inputUserName" placeholder={t('user name')} ref={userNameRef} />
                    </Card.Text>
                    <Button onClick={handleRegistration} disabled={!hasWebAuthN || loadingCreateUser || loadingAddPasskey}>{t('Register')}</Button>{' '}
                    <Button type="submit" disabled={!hasWebAuthN || loadingCreateUser || loadingAddPasskey}>{t('Login')}</Button>
                </Form>
            </Card.Body>
        </Card>
    );
}

export default AuthenticateCard;
