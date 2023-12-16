import { Col, Row, Spinner } from "react-bootstrap";
import SignupCard from "../components/Signup";
import { useOutletContext } from "react-router-dom";
import { ContextType } from "./root";
import { useTranslation } from "react-i18next";
import { gql } from '../__generated__/gql';
import { useQuery } from "@apollo/client";

const ME_GQL = gql(`
query meIndex {
  me {
    displayName
    name
  }
}`);


export default function Index() {
    const { t } = useTranslation();
    const { setFlashMessage } = useOutletContext<ContextType>()
    const { data, loading, refetch } = useQuery(ME_GQL);

    const handleError = (e: string) => setFlashMessage({ msg: e, type: "danger" })

    if (loading) return <Spinner animation="border" />;

    if (data?.me) {
        return (<h1>Welcome {data.me.displayName ? data.me.displayName : data.me.name}</h1>);
    }

    return (
        <Row>
            <Col md={6} xs={true}>
                <SignupCard
                    onError={handleError}
                    onUserCreated={refetch}
                    onPasskeyAdded={() => {
                        setFlashMessage({ msg: t("Key creation succeeded"), type: "success" })
                    }} />
            </Col>
        </Row>
    );
}
