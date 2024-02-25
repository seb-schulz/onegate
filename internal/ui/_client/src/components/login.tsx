import { Button, Spinner } from 'react-bootstrap';
import * as urql from 'urql';
import { gql } from '../__generated__/gql';
import * as graphql from '../__generated__/graphql';
import { startAuthentication } from '@simplewebauthn/browser';
import { useEffect, useRef, useState } from 'react';

const BEGIN_LOGIN_QGL = gql(`
mutation beginLogin {
  beginLogin
}
`);

const VALIDATE_LOGIN_QGL = gql(`
mutation validateLogin($body: CredentialRequestResponse!) {
    validateLogin(body: $body) {
        redirectURL
    }
}
`);

async function handleLogin({ beginLogin, validateLogin, onSuccess, onError }: {
    beginLogin: urql.UseMutationExecute<graphql.BeginLoginMutation, graphql.Exact<{
        [key: string]: never;
    }>>,
    validateLogin: urql.UseMutationExecute<graphql.ValidateLoginMutation, graphql.Exact<{
        body: any;
    }>>
    onSuccess?: (redirectURL?: string) => void
    onError?: (errMsg: string) => void
}) {
    if (onSuccess === undefined) {
        onSuccess = (redirectURL?: string) => { };
    }

    if (onError === undefined) {
        onError = (_: string) => { };
    }

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

        if (!!verificationResp.data.validateLogin) {
            onSuccess(verificationResp.data.validateLogin?.redirectURL);
        }
    } catch (error) {
        onError((error as urql.CombinedError).message);
    }
}

export function LoginSpinner({ onError, onSuccess }: {
    onSuccess?: (redirectURL?: string) => void
    onError?: (errMsg: string) => void
}) {
    const executed = useRef(false)
    const [, beginLogin] = urql.useMutation(BEGIN_LOGIN_QGL);
    const [, validateLogin] = urql.useMutation(VALIDATE_LOGIN_QGL);

    if (!window.PublicKeyCredential) {
        return (
            <p>WebAuthN is not supported.</p>
        )
    }

    useEffect(() => {
        if (!executed.current) {
            handleLogin({ beginLogin, validateLogin, onSuccess, onError })
        }
        return () => {
            executed.current = true;
        };
    }, [])
    return (<Spinner animation="border" />);
}

export function LoginButton({ onError, onSuccess, children }: {
    onSuccess?: (redirectURL?: string) => void
    onError?: (errMsg: string) => void
    children: string | JSX.Element | JSX.Element[]
}) {
    const [{ fetching: fetchingBeginLogin }, beginLogin] = urql.useMutation(BEGIN_LOGIN_QGL);
    const [{ fetching: fetchingValidateLogin }, validateLogin] = urql.useMutation(VALIDATE_LOGIN_QGL);
    const [spinner, setSpinner] = useState(false);

    const onHookedSuccess = (redirectURL?: string) => {
        if (!!onSuccess) onSuccess(redirectURL);
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
        await handleLogin({ beginLogin, validateLogin, onSuccess: onHookedSuccess, onError: onHookedError })
    };

    return (
        <Button disabled={!window.PublicKeyCredential} onClick={handleSubmit}>{children}</Button>
    )
}
