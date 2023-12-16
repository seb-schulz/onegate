import { Button, ButtonGroup, Table } from "react-bootstrap";
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

    const tbody = data?.me?.credentials!.map((cred, idx) =>
        <tr key={cred!.id}>
            <td>{"Credential " + (idx + 1)}</td>
            <td><Moment fromNow withTitle>{cred!.createdAt}</Moment></td>
            <td><Moment fromNow withTitle>{cred!.updatedAt}</Moment></td>
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
