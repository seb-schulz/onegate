import * as React from 'react'
import ReactDOM from 'react-dom/client';
import { ApolloClient, InMemoryCache, ApolloProvider, gql, createHttpLink } from '@apollo/client';
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import 'bootstrap/dist/css/bootstrap.min.css';

import ErrorPage from './error-page';
import Root from './routes/root';
import Index from './routes';
import Credentials from './routes/credentials';

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
    link: createHttpLink({
        uri: '/query',
    }),
    cache: new InMemoryCache(),
});

const router = createBrowserRouter([
    {
        path: "/",
        element: <Root />,
        errorElement: <ErrorPage />,
        children: [
            { index: true, element: <Index /> },
            { path: "credentials", element: <Credentials /> }
        ]
    },
]);

const root = ReactDOM.createRoot(
    document.getElementById('root') as HTMLElement
);

root.render(
    <React.StrictMode>
        <ApolloProvider client={client}>
            <RouterProvider router={router} />
        </ApolloProvider>
    </React.StrictMode >
);
