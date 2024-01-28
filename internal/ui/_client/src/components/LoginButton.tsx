import { Button, Spinner } from 'react-bootstrap';
import * as urql from 'urql';
import { gql } from '../__generated__/gql';
import { startAuthentication } from '@simplewebauthn/browser';
import { useState } from 'react';

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

export default function LoginButton({ onError, onSuccess, children }: {
    onSuccess?: () => void
    onError?: (errMsg: string) => void
    children: string | JSX.Element | JSX.Element[]
}) {
    const [{ fetching: fetchingBeginLogin }, beginLogin] = urql.useMutation(BEGIN_LOGIN_QGL);
    const [{ fetching: fetchingValidateLogin }, validateLogin] = urql.useMutation(VALIDATE_LOGIN_QGL);
    const [spinner, setSpinner] = useState(false);

    const onHookedSuccess = () => {
        if (!!onSuccess) onSuccess();
        setSpinner(false)
    }

    const onHookedError = (errMsg: string) => {
        if (!!onError) onError(errMsg);
        setSpinner(false)
    }

    if (!spinner && (fetchingBeginLogin || fetchingValidateLogin)) setSpinner(true);

    if (spinner) return (<Spinner animation="border" />);

    const handleSubmit = async (e: React.SyntheticEvent) => {
        e.preventDefault();
        e.stopPropagation();
        try {
            const result = await beginLogin({});
            if (!result || !result.data) {
                onHookedError("cannot request data")
                return;
            }

            const asseResp = await startAuthentication(result.data.beginLogin.publicKey);
            const verificationResp = await validateLogin({ body: JSON.stringify(asseResp) });

            if (!verificationResp || !verificationResp.data) {
                onHookedError("cannot request data")
                return;
            }

            if (verificationResp.data.validateLogin) {
                onHookedSuccess();
            }
        } catch (error) {
            onHookedError((error as urql.CombinedError).message);
        }
    };

    return (
        <Button disabled={!window.PublicKeyCredential} onClick={handleSubmit}>{children}</Button>
    )
}
