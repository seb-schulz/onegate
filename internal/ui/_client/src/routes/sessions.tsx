import { useTranslation } from "react-i18next";
import { useOutletContext } from "react-router-dom";
import { ContextType } from "./root";
import { ApolloError, useQuery, useMutation } from "@apollo/client";
import { gql } from '../__generated__/gql';
import * as graphql from '../__generated__/graphql';
import { Button, ButtonGroup, Col, ListGroup, Row, Spinner, Stack, Table } from "react-bootstrap";
import Moment from "react-moment";
import React from "react";

const SESSIONS_GQL = gql(`
query sessions {
  sessions {
    id
    createdAt
    updatedAt
    isActive
  }
}
`)

const REMOVE_SESSION_GQL = gql(`
mutation removeSession($id: ID!) {
  removeSession(id: $id)
}
`)

function Entry({ variant, item, idx, onRemoval }: {
    variant: "list" | "table"
    item: graphql.Session
    idx: number
    onRemoval: (e: React.SyntheticEvent) => Promise<void>
}) {
    const { t } = useTranslation();

    const description = t("Session {{ n }}", { n: idx + 1 })
    const delAction = <Button variant="outline-danger" onClick={onRemoval}><i className="bi bi-trash" /></Button>

    if (variant === "list") {
        return (
            <ListGroup.Item variant={item.isActive ? '' : 'secondary'}>
                <Row><big>{description}</big></Row>
                <Row>
                    <Col sm={true}>
                        <Stack direction="horizontal" gap={1}>
                            <strong>{t("Created")}</strong>
                            <Moment fromNow withTitle>{item!.createdAt}</Moment>
                        </Stack>
                    </Col>
                </Row>
                <Row>
                    <Col sm={true}>
                        <Stack direction="horizontal" gap={1}>
                            <strong>{t("Updated")}</strong>
                            <Moment fromNow withTitle>{item!.updatedAt}</Moment>
                        </Stack>
                    </Col>
                </Row>
                <div className="d-flex flex-row-reverse"><ButtonGroup size="sm">
                    {delAction}
                </ButtonGroup></div>
            </ListGroup.Item>
        );
    } else if (variant === "table") {
        return (
            <tr>
                <td className={item.isActive ? '' : 'bg-secondary-subtle'}>{description}</td>
                <td><Moment fromNow withTitle>{item.createdAt}</Moment></td>
                <td><Moment fromNow withTitle>{item.updatedAt}</Moment></td>
                <td className="d-flex justify-content-end">
                    {delAction}
                </td>
            </tr>
        )
    }
}

export default function Sessions() {
    const { t } = useTranslation();
    const { setFlashMessage } = useOutletContext<ContextType>()
    const { loading, data, refetch } = useQuery(SESSIONS_GQL, {
        onError: (error) => {
            setFlashMessage({ msg: (error as ApolloError).message, type: "danger" })
        }
    });
    const [removeSession, { loading: loadingRemoveSession }] = useMutation(REMOVE_SESSION_GQL);

    const RemoveHandler = (id: string) => (async (e: React.SyntheticEvent) => {
        e.preventDefault();
        e.stopPropagation();

        try {
            const result = await removeSession({
                variables: {
                    id: id
                }
            })
            if (!!result?.data?.removeSession) {
                refetch();
            }
        } catch (error) {
            setFlashMessage({ msg: (error as ApolloError).message, type: "danger" });
        }
    });

    if (loading || loadingRemoveSession) return <Spinner animation="border" />;

    const sessionList = data?.sessions!.map((session, idx) => <Entry variant="list" item={session as graphql.Session} idx={idx} key={session!.id} onRemoval={RemoveHandler(session!.id)} />);

    const sessionTable = data?.sessions!.map((session, idx) => <Entry variant="table" item={session as graphql.Session} idx={idx} key={session!.id} onRemoval={RemoveHandler(session!.id)} />);

    return (
        <Row><Col>
            <ListGroup className="d-md-none pb-2">{sessionList}</ListGroup>
            <Table className="d-none d-md-table" responsive>
                <thead>
                    <tr>
                        <th>{t("Description")}</th>
                        <th>{t("Created")}</th>
                        <th>{t("Updated")}</th>
                        <th></th>
                    </tr>
                </thead>
                <tbody>{sessionTable}</tbody>
            </Table>
        </Col></Row>
    )
}
