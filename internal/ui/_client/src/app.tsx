import * as React from 'react'
import ReactDOM from 'react-dom/client';
import { ApolloClient, InMemoryCache, ApolloProvider, gql, createHttpLink } from '@apollo/client';
import { setContext } from '@apollo/client/link/context';
import i18n from "i18next";
import { initReactI18next } from "react-i18next";

import 'bootstrap/dist/css/bootstrap.min.css';

import { Button, Col, Container, Navbar, Row, Stack } from 'react-bootstrap';
import AuthenticateCard from './Authenticate';

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
    console.log(JWT_HEADER);

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

if (localStorage.getItem(JWT_TOKEN_KEY) === null) {
    localStorage.setItem(JWT_TOKEN_KEY, (document.getElementsByName("jwt-init-token")[0] as HTMLMetaElement).content)
}

root.render(
    <React.StrictMode>
        <ApolloProvider client={client}>
            <Stack gap={2}>

                <Navbar expand="lg" className="bg-body-tertiary">
                    <Container>
                        <Navbar.Brand>OneGate</Navbar.Brand>
                    </Container>
                </Navbar>
                <Container>
                    <Row>
                        <Col md={6} xs={true}><AuthenticateCard /></Col>
                    </Row>
                </Container>
            </Stack>


        </ApolloProvider>
    </React.StrictMode >
);
