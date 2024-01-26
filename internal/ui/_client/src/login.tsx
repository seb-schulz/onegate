import * as React from 'react'
import ReactDOM from 'react-dom/client';
import 'bootstrap/dist/css/bootstrap.min.css';
import './login.css'
import { Button, Card, Modal } from 'react-bootstrap';
import Provider from './client';
import { useTranslation } from 'react-i18next';

const root = ReactDOM.createRoot(
    document.getElementById('root') as HTMLElement
);

function LoginButton() {
    const { t } = useTranslation();

    return (
        <Button>{t('Login with passkey')}</Button>
    )
}

function CentralCard() {
    const { t } = useTranslation();

    return (
        <Card className="shadow text-center mt-5 login-card m-auto">
            <Card.Body>
                <Card.Title>{t('One Gate')}</Card.Title>
                <LoginButton />
            </Card.Body>
        </Card>
    );
}

root.render(
    <React.StrictMode>
        <Provider><CentralCard /></Provider>
    </React.StrictMode >
);
