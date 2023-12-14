import { Col, Container, Navbar, Row, Stack, Toast, ToastContainer } from "react-bootstrap";
import AuthenticateCard from "./Authenticate";
import { useTranslation } from "react-i18next";
import { useState } from "react";
import { gql, useQuery } from "@apollo/client";
import NavbarLogin from "./NavbarLogin";

const ME_GQL = gql`query me {
    me {
      displayName
      name
    }
  }`


function ErrorAlert({ setChildren, children }: {
    setChildren: (msg: string) => void
    children: string
}) {
    return <ToastContainer
        className="p-3"
        position="top-center"
        style={{ zIndex: 1 }}
    >
        <Toast className="d-inline-block m-1" bg="danger" onClose={() => setChildren("")} show={!!children} delay={3000} autohide>
            <Toast.Header>
                <strong className="me-auto">Error</strong>
            </Toast.Header>
            <Toast.Body>{children}</Toast.Body>
        </Toast></ToastContainer >
}


function AppContainer() {
    const { t } = useTranslation();
    const { loading, error, data, refetch } = useQuery(ME_GQL);
    const [errMsg, setErrMsg] = useState("")

    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error : {error.message}</p>;

    const loggedOut = !data.me;
    console.log(data)

    return (
        <Stack gap={2}>

            <Navbar expand="lg" className="bg-body-tertiary">
                <Container>
                    <Navbar.Brand>OneGate</Navbar.Brand>
                    <NavbarLogin me={data.me} loginError={(e) => {
                        console.log(e)
                        setErrMsg(e)
                    }} loginSucceeded={refetch} />
                </Container>
            </Navbar>
            <ErrorAlert setChildren={setErrMsg}>{errMsg}</ErrorAlert>
            <Container>
                <Row>
                    <Col md={6} xs={true}>{loggedOut ? <AuthenticateCard loginSucceeded={() => {
                        console.log("Login succesfull")
                        refetch()
                    }} /> : 'You are logged in'}</Col>
                </Row>
            </Container>
        </Stack>
    );
}
export default AppContainer;
