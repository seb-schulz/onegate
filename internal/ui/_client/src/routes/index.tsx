import { Col, Row } from "react-bootstrap";
import SignupCard from "../components/Signup";
import { useOutletContext } from "react-router-dom";
import { ContextType } from "./root";
import { useTranslation } from "react-i18next";

export default function Index() {
    const { t } = useTranslation();
    const { me, refetchMe, setFlashMessage } = useOutletContext<ContextType>()

    const handleError = (e: string) => setFlashMessage({ msg: e, type: "danger" })

    if (me) {
        return (
            <Row>
                <Col>
                    <h1>Welcome {me.displayName ? me.displayName : me.name}</h1>
                </Col>
            </Row>
        );
    }

    return (
        <Row>
            <Col md={6} xs={true}>
                <SignupCard
                    onError={handleError}
                    onUserCreated={refetchMe}
                    onPasskeyAdded={() => {
                        setFlashMessage({ msg: t("Key creation succeeded"), type: "success" })
                    }} />
            </Col>
        </Row>
    );
}
