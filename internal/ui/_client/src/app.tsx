import * as React from 'react'
import ReactDOM from 'react-dom/client';
import { ApolloClient, InMemoryCache, ApolloProvider, gql } from '@apollo/client';

import Greet from './Greet';

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
            <Greet />
        </ApolloProvider>
    </React.StrictMode>
);
