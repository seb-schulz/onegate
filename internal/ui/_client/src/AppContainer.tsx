import { Alert, Col, Container, Navbar, Row, Stack, Toast, ToastContainer } from "react-bootstrap";
import { Variant } from "react-bootstrap/types";
import SignupCard from "./Signup";
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


function InfoToast({ bg, setChildren, children }: {
    bg: Variant | undefined
    setChildren: (msg: string) => void
    children: string
}) {
    return <ToastContainer
        className="p-3"
        position="top-center"
        style={{ zIndex: 1 }}
    >
        <Toast className="d-inline-block m-1" bg={bg} onClose={() => setChildren("")} show={!!children} delay={3000} autohide>
            <Toast.Header>
                <strong className="me-auto">Info</strong>
            </Toast.Header>
            <Toast.Body>{children}</Toast.Body>
        </Toast></ToastContainer >
}


function AppContainer() {
    const { t } = useTranslation();
    const { loading, error, data, refetch } = useQuery(ME_GQL);
    const [infoMsg, setInfoMsg] = useState("")
    const [infoType, setInfoType] = useState<Variant>("danger")

    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error : {error.message}</p>;

    const loggedOut = !data.me;

    const handleError = (e: string) => {
        setInfoType("danger")
        setInfoMsg(e)
    };

    return (
        <Stack gap={2}>
            <Navbar expand="lg" className="bg-body-tertiary">
                <Container>
                    <Navbar.Brand>OneGate</Navbar.Brand>
                    <NavbarLogin me={data.me} onError={handleError} onSuccess={() => {
                        refetch()
                        setInfoType("success")
                        setInfoMsg(t("Login succeeded"))
                    }} />
                </Container>
            </Navbar>
            <InfoToast bg={infoType} setChildren={setInfoMsg}>{infoMsg}</InfoToast>
            <Container>
                {!window.PublicKeyCredential ? <Row><Alert variant="danger">{t('This browser does not support WebAuthN.')}</Alert></Row> : ''}
                <Row>
                    <Col md={6} xs={true}>{loggedOut ? <SignupCard
                        onError={handleError}
                        onUserCreated={refetch} onPasskeyAdded={() => {
                            setInfoType("success")
                            setInfoMsg(t("Key creation succeeded"))
                        }} /> : 'You are logged in'}</Col>
                </Row>
            </Container>
        </Stack>
    );
}
export default AppContainer;
