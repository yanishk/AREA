import { useNavigate, useLocation } from "react-router-dom";
import ButtonIfCustom from "../components/buttons/ButtonIfCustom";
import { useState } from "react";

const WorkflowBuilder = () => {
    const navigate = useNavigate();
    const location = useLocation();

    const [selectedAction, setSelectedAction] = useState(location.state?.selectedAction || null);

    const handleDeleteAction = () => {
        setSelectedAction(null);
        navigate(location.pathname, { state: { ...location.state, selectedAction: null } });
    };

    return (
        <div>
            <h1>Create</h1>
            <div>
                {selectedAction ? (
                    <ButtonIfCustom text={selectedAction.name} />
                ) : (
                    <button onClick={() => navigate("/select-action")}>If This</button>
                )}
                    <button
                        onClick={handleDeleteAction}
                    >
                        delete
                    </button>
            </div>
        </div>
    );
};

export default WorkflowBuilder