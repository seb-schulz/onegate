import { Alert, Button, ButtonGroup, Col, Container, Form, InputGroup, Navbar, Row, Stack, Table, Toast, ToastContainer } from "react-bootstrap";
import { Variant } from "react-bootstrap/types";
import SignupCard from "./Signup";
import { useTranslation } from "react-i18next";
import { useState } from "react";
import { useQuery } from "@apollo/client";
import { gql } from './__generated__/gql';
import { Credential } from "./__generated__/graphql";
import NavbarLogin from "./NavbarLogin";
import Moment from "react-moment";

const ME_GQL = gql(`
query me {
  me {
    displayName
    name
    credentials {
      id
      createdAt
      updatedAt
    }
  }
}`);


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

function Credentials({ credentials }: {
    credentials: Credential[]
}) {
    const { t } = useTranslation();

    const tbody = credentials.map((cred, idx) => <tr key={cred.id}>
        <td>{"Credential " + (idx + 1)}</td>
        <td><Moment fromNow withTitle>{cred.createdAt}</Moment></td>
        <td><Moment fromNow withTitle>{cred.updatedAt}</Moment></td>
        <td><ButtonGroup size="sm"><Button disabled>{t("Edit")}</Button><Button disabled>{t("Remove")}</Button></ButtonGroup></td>
    </tr>
    );


    return (
        <>
            <Table responsive>
                <thead>
                    <tr>
                        <th>{t("Description")}</th>
                        <th>{t("Created at")}</th>
                        <th>{t("Updated at")}</th>
                        <th>{t("Action")}</th>
                    </tr>
                </thead>
                <tbody>{tbody}</tbody>
            </Table>
            <Button disabled>Add</Button>
        </>
    )
}

function AppContainer() {
    const { t } = useTranslation();
    const { loading, error, data, refetch } = useQuery(ME_GQL);
    const [infoMsg, setInfoMsg] = useState("")
    const [infoType, setInfoType] = useState<Variant>("danger")

    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error : {error.message}</p>;
    if (!data) return <p>Cannot load data</p>


    const handleError = (e: string) => {
        setInfoType("danger")
        setInfoMsg(e)
    };

    let row;
    if (!data || !data.me) {
        row = <Col md={6} xs={true}>
            <SignupCard
                onError={handleError}
                onUserCreated={refetch} onPasskeyAdded={() => {
                    setInfoType("success")
                    setInfoMsg(t("Key creation succeeded"))
                }} />
        </Col>;
    } else {
        row = <Col><h1>{t("Credentials")}</h1><Credentials credentials={data.me.credentials as Credential[]} /></Col>;
    }


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
                    {row}
                </Row>
            </Container>
        </Stack>
    );
}
export default AppContainer;
