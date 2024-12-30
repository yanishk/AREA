import { useEffect, useState } from "react";
import CancelButton from "../buttons/CancelButton";
import TriggerButton from "../buttons/TriggerButton";
import BasicInput from "../inputs/BasicInput";
import MaxWidthWrapper from "../MaxWidthWrapper";
import axios from "axios";

interface SummaryStepProps {
    selectedActionService: {name: string; color: string; iconPath: string } | null;
    selectedReactionService: {name: string; color: string; iconPath: string } | null;
    selectedAction: {id: string; name: string;} | null;
    selectedReaction: {id: string; name: string;} | null;
    actionFields: Record<string, string>;
    reactionFields: Record<string, string>;
    goBack: () => void;
    handleCreateWorkflow: (workflowTitle: string) => void;
}

const SummaryStep: React.FC<SummaryStepProps> = ({
    selectedActionService,
    selectedReactionService,
    selectedAction,
    selectedReaction,
    actionFields,
    reactionFields,
    goBack,
    handleCreateWorkflow
}) => {
    const actionDescription = selectedAction
        ? `${selectedAction.name} ${Object.values(actionFields).join(" ")}`
        : "";

    const reactionDescription = selectedReaction
        ? `${selectedReaction.name} ${Object.values(reactionFields).join(" ")}`
        : "";

    const defaultWorkflowTitle = `If ${actionDescription}, then ${reactionDescription}`;

    const [workflowTitle, setWorkflowTitle] = useState<string>(defaultWorkflowTitle);

    const handleTitleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setWorkflowTitle(event.target.value);
    };

    const [userName, setUserName] = useState("You");

    const getUserName = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}user/email`, {
                withCredentials: true,
            });

            const email = result.data.email;

            const username = email.split('@')[0].replace(/\./g, '');

            setUserName(username);
        } catch (error) {
            console.error("Failed to get email : ", error);
        }
    };

    useEffect(() => {
        getUserName();
    }, []);

    return (
        <div
            className="min-h-screen bg-gray-100"
        >
            <div
                style={{ backgroundColor: selectedActionService?.color }}
                className="flex items-center justify-between px-10 py-5 border-b-[1px] border-[#dadada]"
            >
                <div>
                    <CancelButton
                        text="Back"
                        textColor="text-white"
                        borderColor="border-white"
                        bgColor={`bg-[${selectedActionService?.color}]`}
                        onClick={() => goBack()}
                    />
                </div>

                <h1
                    className="text-white text-[3.2rem] font-extrabold text-center flex-1"
                >
                    Review and finish
                </h1>

                <div className="w-[100px]"></div>
            </div>

            <div
                style={{ backgroundColor: selectedActionService?.color }}
                className="flex flex-col items-center justify-center"
            >
                <MaxWidthWrapper
                    maxWidth="680px"
                >
                    <div
                        className="flex flex-row mt-14 items-start justify-start"
                    >
                        <img
                            alt="icon"
                            src={selectedActionService?.iconPath}
                            className="w-[50px] h-[50px]"
                        />
                        <img
                            alt="icon"
                            src={selectedReactionService?.iconPath}
                            className="w-[50px] h-[50px] ml-4"
                        />
                    </div>
                    <p
                        className="text-white text-[20px] font-normal mt-6 mb-2"
                    >
                        Workflow Title
                    </p>
                    <BasicInput
                        type="text"
                        value={workflowTitle}
                        onChange={handleTitleChange}
                    />
                    <p
                        className="text-white text-[16.7px] mt-[2px] mb-[110px]"
                    >
                        by
                        <span className="font-bold ml-1">
                            {userName}
                        </span>
                    </p>
                </MaxWidthWrapper>
            </div>

            <div
                className="mt-5 flex flex-col !items-center "
            >
                <div className="w-full max-w-[300px]">
                    <TriggerButton
                        text="Finish"
                        textColor="text-white"
                        bgColor="bg-[#222222]"
                        hoverColor="hover:bg-[#333333]"
                        onClick={() => handleCreateWorkflow(workflowTitle)}
                    />
                </div>
            </div>
        </div>
    );
};

export default SummaryStep;
