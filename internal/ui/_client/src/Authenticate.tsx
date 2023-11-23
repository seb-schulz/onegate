import React, { useState } from "react";
import { Alert, Button, Card, Form } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import { gql, useQuery } from "@apollo/client";

const CREATE_CREDENTIAL_OPTIONS_GQL = gql`
query {
  createCredentialOptions {
    challenge
    rp {
      name
      id
    }
    pubKeyCredParams {
      type
      alg
    }
    user {
      id
    }
  }
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

    const { loading, error, data } = useQuery(CREATE_CREDENTIAL_OPTIONS_GQL);
    // if (loading) return <p>Loading...</p>;
    if (error) return <p>Error : {error.message}</p>;

    const handleSubmit = (event: React.SyntheticEvent) => {
        event.preventDefault();
        event.stopPropagation();
        console.log("login", userName)
        setValidated(true);
    };

    const handleRegistration = (event: React.SyntheticEvent) => {
        event.preventDefault();
        event.stopPropagation();
        setValidated(true);

        console.log(data.createCredentialOptions.pubKeyCredParams.map((param: { alg: COSEAlgorithmIdentifier, type: PublicKeyCredentialType }) => {
            return { alg: param.alg, type: param.type }
        }))
        if (!userName) return;


        const createCredentialOptions: PublicKeyCredentialCreationOptions = {
            challenge: b64ToUArray(data.createCredentialOptions.challenge),
            rp: {
                name: data.createCredentialOptions.rp.name,
                id: data.createCredentialOptions.rp.id,
            },
            pubKeyCredParams: data.createCredentialOptions.pubKeyCredParams.map((param: { alg: COSEAlgorithmIdentifier, type: PublicKeyCredentialType }) => {
                return { alg: param.alg, type: param.type }
            }),
            user: {
                id: b64ToUArray(data.createCredentialOptions.user.id),
                name: userName,
                displayName: ''
            }


        };

        console.log("register", createCredentialOptions)

        navigator.credentials
            .create({ publicKey: createCredentialOptions })
            .then(newCredentialInfo => console.log(newCredentialInfo))
            .catch((err) => {
                console.error(err);
            });

    };

    return (
        <Card>
            <Card.Body>
                {hasWebAuthN ? '' : <Alert variant="danger">{t('This browser does not support WebAuthN.')}</Alert>}
                <Form noValidate validated={validated} onSubmit={handleSubmit}>

                    <Card.Text>
                        <Form.Control required type="text" id="inputUserName" placeholder={t('user name')} value={userName} onChange={e => setUserName(e.target.value)} />
                    </Card.Text>
                    <Button onClick={handleRegistration} disabled={!hasWebAuthN || loading}>{t('Register')}</Button>{' '}
                    <Button type="submit" disabled={!hasWebAuthN || loading}>{t('Login')}</Button>
                </Form>
            </Card.Body>
        </Card>
    );
}

export default AuthenticateCard;
