import { useState } from "react";
import { Variant } from "react-bootstrap/types";
import { useTranslation } from "react-i18next";
import { Outlet, NavLink } from "react-router-dom";
import { Alert, Container, Nav, Navbar, Row, Spinner, Stack, Toast, ToastContainer } from "react-bootstrap";
import NavbarLogin from "../components/NavbarLogin";
import * as graphql from '../__generated__/graphql';
import { gql } from '../__generated__/gql';
import * as urql from 'urql';

const ME_GQL = gql(`
query me {
  me {
    displayName
    name
  }
}`);

function FlashMessageToast({ bg, emptyChildren, children }: {
    bg: Variant | undefined
    emptyChildren: () => void
    children: string
}) {
    return (
        <ToastContainer className="p-3" position="bottom-center" style={{ zIndex: 1 }}>
            <Toast className="d-inline-block m-1" bg={bg} onClose={emptyChildren} show={!!children} delay={3000} autohide>
                <Toast.Header>
                    <strong className="me-auto">Info</strong>
                </Toast.Header>
                <Toast.Body>{children}</Toast.Body>
            </Toast>
        </ToastContainer >
    );
}

interface FlashMessageType {
    msg: string
    type: Variant
}

export type ContextType = {
    setFlashMessage: (value: FlashMessageType) => void, me?: graphql.User,
    refetchMe: urql.UseQueryExecute,
};

export default function Root() {
    const { t } = useTranslation();
    const [flashMessage, setFlashMessage] = useState<FlashMessageType>({ msg: "", type: "danger" });
    const [{ fetching, data }, refetch] = urql.useQuery({ query: ME_GQL });

    const handleError = (e: string) => setFlashMessage({ msg: e, type: "danger" });

    if (fetching) return <Spinner animation="border" />;

    return (
        <Stack gap={2}>
            <Navbar expand="lg" className="bg-body-tertiary">
                <Container>
                    <Navbar.Brand href="/">OneGate</Navbar.Brand>
                    <Navbar.Collapse>
                        <Nav className="me-auto">
                            <NavLink
                                to="/credentials"
                                className={({ isActive, isPending }) =>
                                    [isPending ? "pending" : isActive ? "active" : "", "nav-link"].join(" ").trim()
                                }
                            >
                                {t("Credentials")}
                            </NavLink>
                            <NavLink
                                to="/sessions"
                                className={({ isActive, isPending }) =>
                                    [isPending ? "pending" : isActive ? "active" : "", "nav-link"].join(" ").trim()
                                }
                            >
                                {t("Sessions")}
                            </NavLink>
                        </Nav>
                    </Navbar.Collapse>
                    <NavbarLogin me={data?.me as graphql.User} onError={handleError} onSuccess={() => {
                        refetch()
                        setFlashMessage({ msg: t("Login succeeded"), type: "success" })
                    }} />
                    <Navbar.Toggle aria-controls="basic-navbar-nav" />
                </Container>
            </Navbar>
            <FlashMessageToast bg={flashMessage.type} emptyChildren={() => setFlashMessage({ ...flashMessage, msg: "" })}>{flashMessage.msg}</FlashMessageToast>
            <Container>
                {!window.PublicKeyCredential ? <Row><Alert variant="danger">{t('This browser does not support WebAuthN.')}</Alert></Row> : ''}
                <Outlet context={{ setFlashMessage, me: data?.me as graphql.User, refetchMe: refetch } satisfies ContextType} />
            </Container>
        </Stack>
    );
}
