import * as React from 'react'
import ReactDOM from 'react-dom/client';
import { ApolloClient, InMemoryCache, ApolloProvider, gql, createHttpLink } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import i18n from "i18next";
import { initReactI18next } from "react-i18next";

import 'bootstrap/dist/css/bootstrap.min.css';

import AppContainer from './AppContainer';

i18n
    .use(initReactI18next) // passes i18n down to react-i18next
    .init({
        // the translations
        // (tip move them in a JSON file and import them,
        // or even better, manage them via a UI: https://react.i18next.com/guides/multiple-translation-files#manage-your-translations-with-a-management-gui)
        lng: "en",
        fallbackLng: "en",
    });

const JWT_TOKEN_KEY = 'jwt-token';

const JWT_HEADER = (document.getElementsByName("jwt-header")[0] as HTMLMetaElement).content;

const authLink = setContext((_, { headers }) => {
    const token = localStorage.getItem(JWT_TOKEN_KEY);

    const h = {
        ...headers
    }
    if (token) {
        h[JWT_HEADER] = token;
    }
    return {
        headers: h
    }
});

const client = new ApolloClient({
    link: authLink.concat(createHttpLink({
        uri: '/query',
    })),
    cache: new InMemoryCache(),
});

const root = ReactDOM.createRoot(
    document.getElementById('root') as HTMLElement
);

(() => {
    type tokenBody = {
        sub: string,
        exp: number,
    }

    const REDEEM_TOKEN = gql`query { redeemToken }`;

    function extractTokenBody(token: string): tokenBody {
        const [_h, encodedBody, _s] = token.split('.');
        return JSON.parse(atob(encodedBody))

    }

    function validToken(token: string | null): boolean {
        if (token === null) return false;
        const { exp } = extractTokenBody(token)
        return Math.floor(Date.now() / 1000) < exp
    }

    function redeemToken() {
        const newToken = 'abcd';
        // localStorage.setItem(JWT_TOKEN_KEY, newToken);
        // const { exp } = extractTokenBody(newToken);
        console.log('redeem toke')
        setupTimeout()
    }

    function setupTimeout() {
        const token = localStorage.getItem(JWT_TOKEN_KEY) as string;
        const { exp } = extractTokenBody(token)
        setTimeout(redeemToken, 1000 * Math.floor(Math.min((exp - Math.floor(Date.now() / 1000)) / 2, 30 * 60 * 1000)))
    }

    if (!validToken(localStorage.getItem(JWT_TOKEN_KEY))) {
        localStorage.setItem(JWT_TOKEN_KEY, (document.getElementsByName("jwt-init-token")[0] as HTMLMetaElement).content)
    }

    // setupTimeout()
    // setTimeout(redeemToken, 1000)
})()

root.render(
    <React.StrictMode>
        <ApolloProvider client={client}>
            <AppContainer />
        </ApolloProvider>
    </React.StrictMode >
);
