import { gql, useQuery } from "@apollo/client"

const HELLO = gql`
{
    hello(name: "Waldo") {
        name
    }
}
`
export default function Greet() {
    const { loading, error, data } = useQuery(HELLO);
    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error : {error.message}</p>;

    return (<h1>Hello, {data.hello.name}!</h1>)
}
