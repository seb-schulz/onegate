import React, { useState } from "react";
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

function b64ToUArray(x: string): Uint8Array {
    return Uint8Array.from(atob(x), c =>
        c.charCodeAt(0))
};

function AuthenticateCard() {
    const { t } = useTranslation();
    const hasWebAuthN = !!window.PublicKeyCredential;
    const [validated, setValidated] = useState(false);
    const [userName, setUserName] = useState<string>("");
    const [errorMsg, setErrorMsg] = useState<string>("");

    if (!hasWebAuthN) {
        setErrorMsg(t('This browser does not support WebAuthN.'))
    }

    const [createUser, { loading: loadingCreateUser, error: errorCreateUser }] = useMutation(CREATE_USER_GQL);
    if (loadingCreateUser) return <p>Loading...</p>;
    if (errorCreateUser) return <p>Error : {errorCreateUser.message}</p>;

    const [addPasskey, { loading: loadingAddPasskey, error: errorAddPasskey }] = useMutation(ADD_PASSKEY_QGL);
    if (loadingAddPasskey) return <p>Loading...</p>;
    if (errorAddPasskey) setErrorMsg(errorAddPasskey.graphQLErrors.map(({ message }) => message).join(', '));

    const handleSubmit = (event: React.SyntheticEvent) => {
        event.preventDefault();
        event.stopPropagation();
        console.log("login", userName)
        setValidated(true);
    };

    const handleRegistration = async (event: React.SyntheticEvent) => {
        event.preventDefault();
        event.stopPropagation();
        setValidated(true);

        if (!userName) return;

        const result = await createUser({
            variables: {
                name: userName,
            },
        });

        let attResp;
        try {
            attResp = await startRegistration(result.data.createUser.publicKey);
        } catch (error) {
            setErrorMsg(error as string);
            throw error;
        }

        console.log("attResp", attResp)

        try {

            const verificationResp = await addPasskey({
                variables: {
                    body: JSON.stringify(attResp),
                },
            });

            console.log(verificationResp);
        } catch (error) {
            setErrorMsg(error as string);
            throw error;
        }

    };



    return (
        <Card>
            <Card.Body>
                {!errorMsg ? '' : <Alert variant="danger">{errorMsg}</Alert>}
                <Form noValidate validated={validated} onSubmit={handleSubmit}>

                    <Card.Text>
                        <Form.Control required type="text" id="inputUserName" placeholder={t('user name')} value={userName} onChange={e => setUserName(e.target.value)} />
                    </Card.Text>
                    <Button onClick={handleRegistration} disabled={!hasWebAuthN || loadingCreateUser || loadingAddPasskey}>{t('Register')}</Button>{' '}
                    <Button type="submit" disabled={!hasWebAuthN || loadingCreateUser || loadingAddPasskey}>{t('Login')}</Button>
                </Form>
            </Card.Body>
        </Card>
    );
}

export default AuthenticateCard;
