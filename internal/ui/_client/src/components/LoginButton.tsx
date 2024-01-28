import { Button, Spinner } from 'react-bootstrap';
import * as urql from 'urql';
import { gql } from '../__generated__/gql';
import { startAuthentication } from '@simplewebauthn/browser';

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
    onSuccess: () => void
    onError: (errMsg: string) => void
    children: string | JSX.Element | JSX.Element[]
}) {
    const [{ fetching: fetchingBeginLogin }, beginLogin] = urql.useMutation(BEGIN_LOGIN_QGL);
    const [{ fetching: fetchingValidateLogin }, validateLogin] = urql.useMutation(VALIDATE_LOGIN_QGL);

    if (fetchingBeginLogin || fetchingValidateLogin) return (<Spinner animation="border" />);


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
                onSuccess();
            }
        } catch (error) {
            onError((error as urql.CombinedError).message);
        }
    };

    return (
        <Button disabled={!window.PublicKeyCredential} onClick={handleSubmit}>{children}</Button>
    )
}
