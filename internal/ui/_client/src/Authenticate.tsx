import React, { useState } from "react";
import { Alert, Button, Card, Form } from "react-bootstrap";
import { useTranslation } from "react-i18next";

function AuthenticateCard() {
    const { t } = useTranslation();
    const hasWebAuthN = !!window.PublicKeyCredential;
    const [validated, setValidated] = useState(false);
    const [userName, setUserName] = useState<string>("");

    const handleSubmit = (event: React.SyntheticEvent) => {
        event.preventDefault();
        event.stopPropagation();
        console.log("login", userName)
        setValidated(true);
    };

    const handleRegistration = (event: React.SyntheticEvent) => {
        event.preventDefault();
        event.stopPropagation();

        console.log("register", userName)
        setValidated(true);
    };


    return (
        <Card>
            <Card.Body>
                {hasWebAuthN ? '' : <Alert variant="danger">{t('This browser does not support WebAuthN.')}</Alert>}
                <Form noValidate validated={validated} onSubmit={handleSubmit}>

                    <Card.Text>
                        <Form.Control required type="text" id="inputUserName" placeholder={t('user name')} value={userName} onChange={e => setUserName(e.target.value)} />
                    </Card.Text>
                    <Button onClick={handleRegistration} disabled={!hasWebAuthN}>{t('Register')}</Button>{' '}
                    <Button type="submit" disabled={!hasWebAuthN}>{t('Login')}</Button>
                </Form>
            </Card.Body>
        </Card>
    );
}

export default AuthenticateCard;
