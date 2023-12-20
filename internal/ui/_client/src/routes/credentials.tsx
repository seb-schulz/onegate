import { ApolloError, useMutation, useQuery } from "@apollo/client";
import { useRef, useState } from "react";
import { Button, ButtonGroup, Col, Form, InputGroup, ListGroup, Row, Spinner, Stack, Table } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import Moment from "react-moment";
import { useOutletContext } from "react-router-dom";
import { gql } from '../__generated__/gql';
import * as graphql from '../__generated__/graphql';
import { ContextType } from "./root";



const ME_GQL = gql(`
query credentials {
  credentials {
    id
    description
    createdAt
    updatedAt
  }
}`);

const UPDATE_GQL = gql(`
mutation updateCredential($id: ID!, $description: String) {
  updateCredential(id: $id, description: $description) {
    id
  }
}`);

function InlineEditingText({ value, size = 'sm', onSubmit, loading }: {
    value: string,
    size?: 'sm' | 'lg',
    onSubmit?: (value?: string) => Promise<void>,
    loading?: boolean
}) {
    const [editMode, setEditMode] = useState(false);
    const inputRef = useRef<HTMLInputElement | null>(null);

    if (!!loading) return <Spinner animation="border" />

    let actions;
    if (editMode) {
        actions = <>
            <Button variant="outline-success">
                <i className="bi bi-check" onClick={(e) => {
                    e.preventDefault();
                    e.stopPropagation();
                    setEditMode(false);
                    if (onSubmit) onSubmit(inputRef.current?.value);
                }} />
            </Button>
            <Button variant="outline-secondary" onClick={() => setEditMode(false)} >
                <i className="bi bi-x" />
            </Button >
        </>
    } else {
        actions = <>
            <Button variant="outline-secondary" onClick={() => setEditMode(true)}>
                <i className="bi bi-pencil" />
            </Button>
        </>
    }

    return (
        <InputGroup size={size}>
            <Form.Control ref={inputRef} readOnly={!editMode} defaultValue={value} />{actions}
        </InputGroup>
    );

}

function CredentialEntry({ modus, credential, idx }: {
    modus: "list" | "table"
    credential: graphql.Credential
    idx: number
}) {
    const { t } = useTranslation();
    const [updateCredential, { loading, error, data }] = useMutation(UPDATE_GQL);
    const { setFlashMessage } = useOutletContext<ContextType>()

    const actions = <ButtonGroup size="sm">
        <Button disabled variant="outline-danger"><i className="bi bi-trash" /></Button>
    </ButtonGroup>;

    const handleSubmit = async (value?: string) => {
        if (!value) return;

        const resp = await updateCredential({
            variables: {
                id: credential.id,
                description: value,
            },
            onError: (error) => {
                setFlashMessage({ msg: (error as ApolloError).message, type: "danger" })
            }
        });

        if (!!resp.data?.updateCredential?.id) {
            setFlashMessage({ msg: "Description saved", type: "success" })
        }
    }

    const description = !!credential.description ? credential.description : "Credential " + (idx + 1)

    if (modus === "list") {
        return (
            <ListGroup.Item>
                <Row className="mb-2"><InlineEditingText value={description} onSubmit={handleSubmit} loading={loading} size="lg" /></Row>
                <Row>
                    <Col sm={true}>
                        <Stack direction="horizontal" gap={1}>
                            <strong>{t("Created")}</strong>
                            <Moment fromNow withTitle>{credential!.createdAt}</Moment>
                        </Stack>
                    </Col>
                </Row>
                <Row>
                    <Col sm={true}>
                        <Stack direction="horizontal" gap={1}>
                            <strong>{t("Updated")}</strong>
                            <Moment fromNow withTitle>{credential!.updatedAt}</Moment>
                        </Stack>
                    </Col>
                </Row>
                <div className="d-flex flex-row-reverse">{actions}</div>
            </ListGroup.Item>
        );
    } else if (modus === "table") {
        return (
            <tr>
                <td>
                    <InlineEditingText value={description} onSubmit={handleSubmit} loading={loading} />
                </td>
                <td><Moment fromNow withTitle>{credential!.createdAt}</Moment></td>
                <td><Moment fromNow withTitle>{credential!.updatedAt}</Moment></td>
                <td className="d-flex justify-content-end">
                    {actions}
                </td>
            </tr>
        )
    }

}

export default function Credentials() {
    const { t } = useTranslation();
    const { loading, error, data } = useQuery(ME_GQL, {
        onError: (error) => {
            setFlashMessage({ msg: (error as ApolloError).message, type: "danger" })
        }
    });
    const { setFlashMessage } = useOutletContext<ContextType>()


    if (loading) return <Spinner animation="border" />;
    if (!data?.credentials) return t("You are logged-out!");

    const actions = <ButtonGroup size="sm">
        <Button disabled variant="outline-secondary" size="sm"><i className="bi bi-pencil" /></Button>
        <Button disabled variant="outline-danger"><i className="bi bi-trash" /></Button>
    </ButtonGroup>;

    const credList = data?.credentials!.map((cred, idx) => <CredentialEntry modus="list" credential={cred as graphql.Credential} idx={idx} key={cred!.id} />);

    const credTable = data?.credentials!.map(
        (cred, idx) => <CredentialEntry modus="table" credential={cred as graphql.Credential} idx={idx} key={cred!.id} />

    );

    return (
        <>
            <Row><Col>
                <ListGroup className="d-md-none pb-2">{credList}</ListGroup>
                <Table className="d-none d-md-table" responsive>
                    <thead>
                        <tr>
                            <th>{t("Description")}</th>
                            <th>{t("Created")}</th>
                            <th>{t("Updated")}</th>
                            <th></th>
                        </tr>
                    </thead>
                    <tbody>{credTable}</tbody>
                </Table>
            </Col></Row>
            <Row><Col><Button disabled>Add</Button></Col></Row>
        </>
    )
}
