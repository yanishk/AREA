import { useNavigate, useLocation } from "react-router-dom";
import ButtonIfCustom from "../components/buttons/ButtonIfCustom";
import { useState } from "react";
import ButtonIfDefault from "../components/buttons/ButtonIfDefault";
import CancelButton from "../components/buttons/CancelButton";
import ButtonThenCustom from "../components/buttons/ButtonThenCustom";
import ButtonThenDefault from "../components/buttons/ButtonThenDefault";

const WorkflowBuilder = () => {
    const navigate = useNavigate();
    const location = useLocation();

    const [selectedAction, setSelectedAction] = useState(location.state?.selectedAction || null);
    const [selectedReaction, setSelectedReaction] = useState(location.state?.selectedReaction || null);

    const handleDeleteAction = () => {
        setSelectedAction(null);
        navigate(location.pathname, { state: { ...location.state, selectedAction: null } });
    };

    const handleDeleteReaction = () => {
        setSelectedReaction(null);
        navigate(location.pathname, { state: { ...location.state, selectedReaction: null } });
    };

    const handleOnEditAction = () => {
        navigate("/select-action", {
            state: {
                selectedReaction,
                selectedAction
            }
        })
    }

    const handleOnEditReaction = () => {
        navigate("/select-reaction", {
            state: {
                selectedReaction,
                selectedAction
            }
        })
    }

    const goBack = () => {
        navigate(-1);
    }

    return (
        <div>
            <h1>Create</h1>
            <CancelButton
                text="Cancel"
                onClick={goBack}
            />
            <div>
                {selectedAction ? (
                    <ButtonIfCustom
                        text={selectedAction.name}
                        color={selectedAction.color}
                        onEdit={handleOnEditAction}
                        onDelete={handleDeleteAction}
                    />
                ) : (
                    <ButtonIfDefault
                        onClick={() =>
                            navigate("/select-action", {
                                state: { selectedReaction },
                            })
                        }
                        text="If This"
                    />
                )}
            </div>
            <div>
                {selectedReaction ? (
                    <ButtonThenCustom
                        text={selectedReaction.name}
                        color={selectedReaction.color}
                        onEdit={handleOnEditReaction}
                        onDelete={handleDeleteReaction}
                    />
                ) : (
                    <ButtonThenDefault
                        onClick={() =>
                            navigate("/select-reaction", {
                                state: { selectedAction },
                            })
                        }
                        text="Then That"
                    />
                )}
            </div>
        </div>
    );
};

export default WorkflowBuilder