import * as React from 'react'
import ReactDOM from 'react-dom/client';
import { ApolloClient, InMemoryCache, ApolloProvider, gql } from '@apollo/client';
import i18n from "i18next";
import { initReactI18next } from "react-i18next";

import 'bootstrap/dist/css/bootstrap.min.css';

import Greet from './Greet';
import { Button, Col, Container, Row } from 'react-bootstrap';
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

const client = new ApolloClient({
    uri: '/query',
    cache: new InMemoryCache(),
});

// client.query({
//     query: gql`
//         {
//         hello(name: "Waldo") {
//             name
//         }
//         }
// `}).then((result) => console.log(result));

const root = ReactDOM.createRoot(
    document.getElementById('root') as HTMLElement
);

root.render(
    <React.StrictMode>
        <ApolloProvider client={client}>
            <Container>
                <Row>
                    <Col><Greet /></Col>
                </Row>
                <Row>
                    <Col xs={6}><AuthenticateCard /></Col>
                </Row>
            </Container>


        </ApolloProvider>
    </React.StrictMode >
);
