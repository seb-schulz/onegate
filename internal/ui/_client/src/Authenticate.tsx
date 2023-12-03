import React, { useState } from "react";
import { Alert, Button, Card, Form } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import { gql, useMutation } from "@apollo/client";

const CREATE_USER_GQL = gql`
mutation {
  createUser {
    challenge
    rp {
      name
      id
    }
    pubKeyCredParams {
      type
      alg
    }
    userID
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

    const [createUser, { loading, error }] = useMutation(CREATE_USER_GQL);
    if (loading) return <p>Loading...</p>;
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

        if (!userName) return;

        createUser()
            .then(result => {
                const data = result.data.createUser;
                const createCredentialOptions = {
                    challenge: b64ToUArray(data.challenge),
                    rp: {
                        name: data.rp.name,
                        id: data.rp.id,
                    },
                    pubKeyCredParams: data.pubKeyCredParams.map((param: { alg: COSEAlgorithmIdentifier, type: PublicKeyCredentialType }) => {
                        return { alg: param.alg, type: param.type }
                    }),
                    user: {
                        id: b64ToUArray(data.userID),
                        name: userName,
                        displayName: ''
                    }
                } as PublicKeyCredentialCreationOptions;

                return navigator.credentials.create({ publicKey: createCredentialOptions });
            })
            .then(newCredentialInfo => console.log(newCredentialInfo))
            .catch(err => console.error(err));
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
