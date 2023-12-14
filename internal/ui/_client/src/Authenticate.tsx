import React, { useState, useRef } from "react";
import { Alert, Button, Card, Form } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import { ApolloError, gql, useMutation } from "@apollo/client";
import { startRegistration, startAuthentication } from '@simplewebauthn/browser';

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

const BEGIN_LOGIN_QGL = gql`
mutation beginLogin {
  beginLogin
}
`

const VALIDATE_LOGIN_QGL = gql`
mutation validateLogin($body: CredentialRequestResponse!) {
    validateLogin(body: $body)
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
    const [beginLogin, { loading: loadingBeginLogin, error: errorBeginLogin }] = useMutation(BEGIN_LOGIN_QGL);
    const [validateLogin, { loading: loadingValidateLogin, error: errorValidateLogin }] = useMutation(VALIDATE_LOGIN_QGL);

    const hasWebAuthN = !!window.PublicKeyCredential;
    const errorMsgList = [];

    if (!hasWebAuthN) errorMsgList.push(t('This browser does not support WebAuthN.'));
    if (!!errorMsg) errorMsgList.push(errorMsg);
    if (loadingCreateUser || loadingAddPasskey || loadingBeginLogin || loadingValidateLogin) return <p>Loading...</p>;
    if (errorCreateUser) errorMsgList.push(errorCreateUser.message);
    if (errorAddPasskey) errorMsgList.push(errorAddPasskey.message);
    if (errorBeginLogin) errorMsgList.push(errorBeginLogin.message);
    if (errorValidateLogin) errorMsgList.push(errorValidateLogin.message);

    const handleSubmit = async (e: React.SyntheticEvent) => {
        e.preventDefault();
        e.stopPropagation();
        const result = await beginLogin();

        try {
            const asseResp = await startAuthentication(result.data.beginLogin.publicKey);

            const verificationResp = await validateLogin({
                variables: {
                    body: JSON.stringify(asseResp),
                },
            });

            if (verificationResp.data.validateLogin) {
                setTimeout(loginSucceeded, 0);
            }
        } catch (error) {
            setErrorMsg((error as ApolloError).message);
        }
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

        try {
            const attResp = await startRegistration(result.data.createUser.publicKey);

            const addPasskeyData = await addPasskey({
                variables: {
                    body: JSON.stringify(attResp),
                },
            });

            if (addPasskeyData.data.addPasskey) {
                // TODO: Add alert on top
            }
        } catch (error) {
            setErrorMsg((error as ApolloError).message);
        }

    };

    return (
        <Card>
            <Card.Body>
                {errorMsgList.length > 0 ? <Alert variant="danger">{errorMsgList.join(',`')}</Alert> : ''}
                <Form noValidate validated={validated} onSubmit={handleSubmit}>

                    <Card.Text>
                        <Form.Control required type="text" id="inputUserName" placeholder={t('user name')} ref={userNameRef} autoComplete="username webauthn" />
                    </Card.Text>
                    <Button onClick={handleRegistration} disabled={!hasWebAuthN || loadingCreateUser || loadingAddPasskey}>{t('Register')}</Button>{' '}
                    <Button type="submit" disabled={!hasWebAuthN || loadingCreateUser || loadingAddPasskey}>{t('Login')}</Button>
                </Form>
            </Card.Body>
        </Card>
    );
}

export default AuthenticateCard;
