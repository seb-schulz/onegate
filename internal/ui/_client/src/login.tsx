import * as React from 'react'
import ReactDOM from 'react-dom/client';
import 'bootstrap/dist/css/bootstrap.min.css';
import './login.css'
import { Card } from 'react-bootstrap';
import Provider from './client';
import { useTranslation } from 'react-i18next';
import LoginButton from './components/LoginButton';

const root = ReactDOM.createRoot(
    document.getElementById('root') as HTMLElement
);

function CentralCard() {
    const { t } = useTranslation();

    return (
        <Card className="shadow text-center mt-5 login-card m-auto">
            <Card.Body>
                <Card.Title>{t('One Gate')}</Card.Title>
                <LoginButton onError={console.error} onSuccess={() => {
                    window.location.href = "/";
                }}>{t('Login with passkey')}</LoginButton>
            </Card.Body>
        </Card>
    );
}

root.render(
    <React.StrictMode>
        <Provider><CentralCard /></Provider>
    </React.StrictMode >
);
