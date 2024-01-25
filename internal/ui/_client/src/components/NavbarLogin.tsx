import { Button, Navbar, Spinner } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import { startAuthentication } from '@simplewebauthn/browser';
import * as urql from 'urql';
import { gql } from '../__generated__/gql';
import * as graphql from '../__generated__/graphql';

const BEGIN_LOGIN_QGL = gql(`
mutation beginLogin {
  beginLogin
}
`);

const VALIDATE_LOGIN_QGL = gql(`
mutation validateLogin($body: CredentialRequestResponse!) {
    validateLogin(body: $body)
}
`);

function NavbarLogin({ me, onError, onSuccess }: {
    me?: graphql.User
    onSuccess: () => void
    onError: (errMsg: string) => void
}) {
    const { t } = useTranslation();
    const [{ fetching: fetchingBeginLogin }, beginLogin] = urql.useMutation(BEGIN_LOGIN_QGL);
    const [{ fetching: fetchingValidateLogin }, validateLogin] = urql.useMutation(VALIDATE_LOGIN_QGL);

    if (fetchingBeginLogin || fetchingValidateLogin) return <Spinner animation="border" />;

    const handleSubmit = async (e: React.SyntheticEvent) => {
        e.preventDefault();
        e.stopPropagation();
        try {
            const result = await beginLogin({});
            if (!result || !result.data) {
                onError("cannot request data")
                return;
            }

            const asseResp = await startAuthentication(result.data.beginLogin.publicKey);
            const verificationResp = await validateLogin({ body: JSON.stringify(asseResp) });

            if (!verificationResp || !verificationResp.data) {
                onError("cannot request data")
                return;
            }

            if (verificationResp.data.validateLogin) {
                setTimeout(onSuccess, 0);
            }
        } catch (error) {
            onError((error as urql.CombinedError).message);
        }
    };

    if (!me) return <>
        <Button disabled={!window.PublicKeyCredential} onClick={handleSubmit}>{t('Login')}</Button>
    </>;

    return (
        <Navbar.Text>
            Name: <a href="/me">
                {me.displayName ? me.displayName : me.name}
            </a>
        </Navbar.Text>
    )
}

export default NavbarLogin;
