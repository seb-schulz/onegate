import { Button, ButtonGroup, Col, ListGroup, Row, Stack, Table } from "react-bootstrap";
import { useTranslation } from "react-i18next";
import { gql } from '../__generated__/gql';
import Moment from "react-moment";
import { useQuery } from "@apollo/client";

const ME_GQL = gql(`
query myCredentials {
  me {
    credentials {
      id
      createdAt
      updatedAt
    }
  }
}`);

export default function Credentials() {
    const { t } = useTranslation();
    const { loading, error, data } = useQuery(ME_GQL);

    if (loading) return <p>Loading...</p>;
    if (error) return <p>Error : {error.message}</p>;
    if (!data?.me) return t("You are logged-out!");

    const credList = data?.me?.credentials!.map((cred, idx) =>
        <ListGroup.Item key={cred!.id}>
            <Row><big>{"Credential " + (idx + 1)}</big></Row>
            <Row>
                <Col sm={true}>
                    <Stack direction="horizontal" gap={1}>
                        <strong>{t("Created")}</strong>
                        <Moment fromNow withTitle>{cred!.createdAt}</Moment>
                    </Stack>
                </Col>
            </Row>
            <Row>
                <Col sm={true}>
                    <Stack direction="horizontal" gap={1}>
                        <strong>{t("Updated")}</strong>
                        <Moment fromNow withTitle>{cred!.updatedAt}</Moment>
                    </Stack>
                </Col>
            </Row>
            <div className="d-flex flex-row-reverse">
                <ButtonGroup size="sm"><Button disabled>{t("Edit")}</Button><Button disabled>{t("Remove")}</Button></ButtonGroup>
            </div>
        </ListGroup.Item>
    );

    const credTable = data?.me?.credentials!.map(
        (cred, idx) =>
            <tr key={cred!.id}>
                <td>{"Credential " + (idx + 1)}</td>
                <td><Moment fromNow withTitle>{cred!.createdAt}</Moment></td>
                <td><Moment fromNow withTitle>{cred!.updatedAt}</Moment></td>
                <td><ButtonGroup size="sm"><Button disabled>{t("Edit")}</Button><Button disabled>{t("Remove")}</Button></ButtonGroup></td>
            </tr>
    );

    return (
        <>
            <Row><Col>
                <ListGroup className="d-sm-none pb-2">{credList}</ListGroup>
                <Table className="d-none d-sm-table" responsive>
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
