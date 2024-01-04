import { Button, Col, Form, Row, Spinner } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import { useOutletContext } from "react-router-dom";
import { ContextType } from "./root";
import { useRef } from "react";
import { gql } from '../__generated__/gql';
import * as graphql from '../__generated__/graphql';
import { ApolloError, useMutation } from "@apollo/client";

const UPDATE_GQL = gql(`
mutation updateMe($name: String, $displayName: String) {
  updateMe(name: $name, displayName: $displayName) {
    name
  }
}
`);

export default function Me() {
    const { t } = useTranslation();
    const { me, refetchMe, setFlashMessage } = useOutletContext<ContextType>()
    const nameRef = useRef<HTMLInputElement | null>(null)
    const displayNameRef = useRef<HTMLInputElement | null>(null)
    const [updateMe, { loading }] = useMutation(UPDATE_GQL);


    if (!me) return (
        <p>Logged out</p>
    )

    if (loading) return <Spinner animation="border" />;


    const handleSubmit = async (e: React.SyntheticEvent) => {
        e.preventDefault();
        e.stopPropagation();

        try {
            await updateMe({
                variables: {
                    name: nameRef.current?.value,
                    displayName: displayNameRef.current?.value,
                },
            });
            await refetchMe();
        } catch (error) {
            setFlashMessage({ msg: (error as ApolloError).message, type: "danger" });
        }
    };

    return (
        <>
            <Form onSubmit={handleSubmit}>
                <Form.Group as={Row} className="mb-3" controlId="userName">
                    <Form.Label column sm="2">
                        {t('Name')}
                    </Form.Label>
                    <Col sm="10">
                        <Form.Control defaultValue={me.name} ref={nameRef} />
                    </Col>
                </Form.Group>
                <Form.Group as={Row} className="mb-3" controlId="userDisplayName">
                    <Form.Label column sm="2">
                        {t('Display Name')}
                    </Form.Label>
                    <Col sm="10">
                        <Form.Control defaultValue={me.displayName} ref={displayNameRef} />
                    </Col>
                </Form.Group>
                <Button type="submit">{t('Save')}</Button>
            </Form>
        </>
    );
}
