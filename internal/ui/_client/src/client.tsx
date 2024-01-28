import * as React from 'react'
import * as urql from 'urql';
import { retryExchange } from '@urql/exchange-retry';


const client = new urql.Client({
    url: '/query',
    exchanges: [urql.cacheExchange, retryExchange({
        initialDelayMs: 1000,
        maxDelayMs: 15000,
        randomDelay: true,
        maxNumberAttempts: 2,
        retryIf: err => !!(err && err.networkError),
    }), urql.fetchExchange],
    fetchOptions: () => {
        return {
            headers: {
                'X-Onegate-Csrf-Protection': '1'
            },
        };
    },
    requestPolicy: 'cache-and-network',
});


export default function Provider({ children }: {
    children: string | JSX.Element | JSX.Element[]
}) {
    return (
        <urql.Provider value={client}>{children}</urql.Provider>
    );
};
