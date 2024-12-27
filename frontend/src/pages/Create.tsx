import { useNavigate, useLocation, useSearchParams } from "react-router-dom";
import ButtonIfCustom from "../components/buttons/ButtonIfCustom";
import { useEffect, useState } from "react";
import ButtonIfDefault from "../components/buttons/ButtonIfDefault";
import CancelButton from "../components/buttons/CancelButton";
import ButtonThenCustom from "../components/buttons/ButtonThenCustom";
import ButtonThenDefault from "../components/buttons/ButtonThenDefault";
import TriggerButton from "../components/buttons/TriggerButton";
import axios from "axios";

const WorkflowBuilder = () => {
    const navigate = useNavigate();

    const [selectedAction, setSelectedAction] = useState(null);
    const [selectedReaction, setSelectedReaction] = useState(null);

    const handleOnEditAction = () => {
        return;
    }

    const handleDeleteAction = () => {
        return;
    }

    const goHome = () => {
        return;
    }

    const handleContinue = () => {
        return;
    }

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
                        <ButtonIfDefault
                            onClick={() =>
                                goHome()
                            }
                            text="If This"
                        />
                    ) : (
                        <ButtonIfDefault
                            onClick={() =>
                                goHome()
                            }
                            text="If This"
                        />
                    )}
                </div>

                <div className="h-[70px] w-[5px] bg-[#EEEEEE]" />

                <div className="mb-[50px]">
                    {selectedReaction ? (
                        <ButtonThenDefault
                        onClick={() =>
                            goHome()
                        }
                        text="Then That"
                        disabled={!selectedAction}
                    />
                    ) : (
                        <ButtonThenDefault
                            onClick={() =>
                                goHome()
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
