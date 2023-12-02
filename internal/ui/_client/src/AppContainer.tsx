import { Col, Container, Navbar, Row, Stack } from "react-bootstrap";
import AuthenticateCard from "./Authenticate";
import { useTranslation } from "react-i18next";
import { useState } from "react";
import { gql, useQuery } from "@apollo/client";

const ME_GQL = gql`{me{PasskeyID}}`


function AppContainer() {
    const { t } = useTranslation();

    const { loading, error, data } = useQuery(ME_GQL);
    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error : {error.message}</p>;

    const loggedOut = !data.me;

    return (
        <Stack gap={2}>

            <Navbar expand="lg" className="bg-body-tertiary">
                <Container>
                    <Navbar.Brand>OneGate</Navbar.Brand>
                    {loggedOut ? '' : <Navbar.Text>ID: {data.me.PasskeyID}</Navbar.Text>}
                </Container>
            </Navbar>
            <Container>
                <Row>
                    <Col md={6} xs={true}>{loggedOut ? <AuthenticateCard /> : 'You are logged in'}</Col>
                </Row>
            </Container>
        </Stack>
    );
}
export default AppContainer;
