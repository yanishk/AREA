import { useNavigate, useLocation } from "react-router-dom";
import ButtonIfCustom from "../components/buttons/ButtonIfCustom";
import { useState } from "react";
import ButtonIfDefault from "../components/buttons/ButtonIfDefault";
import CancelButton from "../components/buttons/CancelButton";
import ButtonThenCustom from "../components/buttons/ButtonThenCustom";
import ButtonThenDefault from "../components/buttons/ButtonThenDefault";
import TriggerButton from "../components/buttons/TriggerButton";
import axios from "axios";

type WorkflowData = {
    selectedAction: {
        name: string;
        id: string;
        color: string;
        icon: string;
    } | null;
    selectedReaction: {
        name: string;
        id: string;
        color: string;
        icon: string;
    } | null;
} | null;

const WorkflowBuilder = () => {
    const navigate = useNavigate();
    const location = useLocation();


    const savedData: WorkflowData = JSON.parse(localStorage.getItem("workflowData") ?? "null");

    const [selectedAction, setSelectedAction] = useState(
        location.state?.selectedAction || savedData?.selectedAction || null
    );

    const [selectedReaction, setSelectedReaction] = useState(
        location.state?.selectedReaction || savedData?.selectedReaction || null
    );


    const handleDeleteAction = () => {
        setSelectedAction(null);
        navigate(location.pathname, {
            state: {
                ...location.state,
                selectedAction: null
            }
        });
    };

    const handleDeleteReaction = () => {
        setSelectedReaction(null);
        navigate(location.pathname, {
            state: {
                ...location.state,
                selectedReaction: null
            }
        });
    };

    const handleOnEditAction = () => {
        return;
    }

    const handleOnEditReaction = () => {
        return;
    }

    const goHome = () => {
        navigate("/");
    }

    const handleContinue = async () => {
        let timeValue;
        let quarterValue;
        if (selectedAction?.fields) {
            const { time, quarter } = selectedAction.fields;
            timeValue = parseInt(time) || 0;
            quarterValue = parseInt(quarter) || 0;
        }

        await axios.post(
            `${import.meta.env.VITE_API_URL}workflow/create`,
            {
                actionid: selectedAction?.id || "",
                actionparam: {
                    hour: timeValue,
                    minute: quarterValue,
                },
                name: "workfloww",
                reactionid: selectedReaction?.id || "",
                reactionparam: {}
            },
            {
                withCredentials: true,
            }
        );

        localStorage.removeItem("workflowData");
        setSelectedAction(null);
        setSelectedReaction(null);
        navigate("/");
    };


    return (
        <div>
            <div className="flex items-center justify-between px-10 py-5">
                <div>
                    <CancelButton
                        text="Cancel"
                        textColor="text-[#222222]"
                        borderColor="border-[#222222]"
                        bgColor="bg-[#F9FAFB]"
                        onClick={goHome}
                    />
                </div>

                <h1 className="text-[3.2rem] font-[900] text-center flex-1">
                    Create
                </h1>

                <div className="w-[140px]"></div>
            </div>

            <div className="flex items-center justify-center flex-col px-10 pt-[100px] py-5">
                <div className="pt-10">
                    {selectedAction ? (
                        <ButtonIfCustom
                            text={selectedAction.name}
                            color={selectedAction.color}
                            iconPath={selectedAction.icon}
                            onEdit={handleOnEditAction}
                            onDelete={handleDeleteAction}
                        />
                    ) : (
                        <ButtonIfDefault
                            onClick={() =>
                                navigate("/select-service", {
                                    state: { type: "action", selectedReaction },
                                })
                            }
                            text="If This"
                        />
                    )}
                </div>

                <div className="h-[70px] w-[5px] bg-[#EEEEEE]" />

                <div className="mb-[50px]">
                    {selectedReaction ? (
                        <ButtonThenCustom
                            text={selectedReaction.name}
                            color={selectedReaction.color}
                            iconPath={selectedReaction.icon}
                            onEdit={handleOnEditReaction}
                            onDelete={handleDeleteReaction}
                        />
                    ) : (
                        <ButtonThenDefault
                            onClick={() =>
                                navigate("/select-service", {
                                    state: { type: "reaction", selectedAction },
                                })
                            }
                            text="Then That"
                            disabled={!selectedAction}
                        />
                    )}
                </div>

                <div
                    className="w-full max-w-[350px] flex flex-col items-center justify-center"
                >
                    {selectedAction && selectedReaction ? (
                        <TriggerButton
                            text="Create Workflow"
                            textColor="text-white"
                            bgColor="bg-[#222222]"
                            hoverColor="hover:bg-[#333333]"
                            onClick={handleContinue}
                        />
                    ) : (
                        <></>
                    )}
                </div>
            </div>
        </div>
    );
};

export default WorkflowBuilder