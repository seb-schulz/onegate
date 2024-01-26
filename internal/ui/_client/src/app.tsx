import * as React from 'react'
import ReactDOM from 'react-dom/client';
import { Client, Provider, cacheExchange, fetchExchange } from 'urql';
import { retryExchange } from '@urql/exchange-retry';
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-icons/font/bootstrap-icons.min.css';
import ErrorPage from './error-page';
import Root from './routes/root';
import Index from './routes';
import Credentials from './routes/credentials';
import Sessions from './routes/sessions';
import Me from './routes/me';

i18n
    .use(initReactI18next) // passes i18n down to react-i18next
    .init({
        // the translations
        // (tip move them in a JSON file and import them,
        // or even better, manage them via a UI: https://react.i18next.com/guides/multiple-translation-files#manage-your-translations-with-a-management-gui)
        lng: "en",
        fallbackLng: "en",
    });


const client = new Client({
    url: '/query',
    exchanges: [cacheExchange, retryExchange({
        initialDelayMs: 1000,
        maxDelayMs: 15000,
        randomDelay: true,
        maxNumberAttempts: 2,
        retryIf: err => !!(err && err.networkError),
    }), fetchExchange],
    fetchOptions: () => {
        return {
            headers: {
                'X-Onegate-Csrf-Protection': '1'
            },
        };
    },
    requestPolicy: 'cache-and-network',
});

const router = createBrowserRouter([
    {
        path: "/",
        element: <Root />,
        errorElement: <ErrorPage />,
        children: [
            { index: true, element: <Index /> },
            { path: "me", element: <Me /> },
            { path: "credentials", element: <Credentials /> },
            { path: "sessions", element: <Sessions /> }
        ]
    },
]);

const root = ReactDOM.createRoot(
    document.getElementById('root') as HTMLElement
);

root.render(
    <React.StrictMode>
        <Provider value={client}>
            <RouterProvider router={router} />
        </Provider>
    </React.StrictMode >
);
