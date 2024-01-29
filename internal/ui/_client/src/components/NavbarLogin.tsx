import { Navbar } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import * as graphql from '../__generated__/graphql';
import { LoginButton } from "./login";


function NavbarLogin({ me, onError, onSuccess }: {
    me?: graphql.User
    onSuccess: () => void
    onError: (errMsg: string) => void
}) {
    const { t } = useTranslation();

    if (!me) return <>
        <LoginButton onSuccess={onSuccess} onError={onError}>{t('Login')}</LoginButton>
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
