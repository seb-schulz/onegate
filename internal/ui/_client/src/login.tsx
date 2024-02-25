import * as React from 'react'
import ReactDOM from 'react-dom/client';
import 'bootstrap/dist/css/bootstrap.min.css';
import './login.css'
import { Alert, Card, Spinner } from 'react-bootstrap';
import Provider from './client';
import { useTranslation } from 'react-i18next';
import { LoginButton, LoginSpinner } from './components/login';

const rootDom = document.getElementById('root') as HTMLElement;
const root = ReactDOM.createRoot(rootDom);
const startLogin = rootDom.dataset['startLogin'] == "1"

function CentralCard() {
    const { t } = useTranslation();
    const [error, setError] = React.useState("")

    const onSuccess = (redirectURL?: string) => {
        console.log(redirectURL)
        if (!redirectURL) {
            window.location.href = "/";
        } else {
            window.location.href = redirectURL;
        }
    }

    const login = (
        <LoginButton onError={setError} onSuccess={onSuccess}>{t('Login with passkey')}</LoginButton>
    )

    const spinner = (
        <LoginSpinner onError={setError} onSuccess={onSuccess} />
    );

    return (
        <Card className="shadow text-center mt-5 login-card m-auto">
            <Card.Body>
                <Card.Title>{t('One Gate')}</Card.Title>
                {error ? <Alert variant="danger">{error}</Alert> : ""}
                {startLogin && !error ? spinner : login}
            </Card.Body>
        </Card>
    );
}

root.render(
    <React.StrictMode>
        <Provider><CentralCard /></Provider>
    </React.StrictMode >
);
