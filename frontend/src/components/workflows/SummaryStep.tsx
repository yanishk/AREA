import { useEffect, useState } from "react";
import ActionButton from "../buttons/ActionButton";
import TriggerButton from "../buttons/TriggerButton";
import BasicInput from "../inputs/BasicInput";
import MaxWidthWrapper from "../MaxWidthWrapper";
import { getUserName } from "../../utils/GetUserName";

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

    let defaultWorkflowTitle

    if (selectedReactionService?.name === "Asana" || selectedReactionService?.name === "Discord") {
        defaultWorkflowTitle = `If ${actionDescription}, then ${selectedReaction?.name}`;
    } else {
        defaultWorkflowTitle = `If ${actionDescription}, then ${reactionDescription}`;
    }

    const [workflowTitle, setWorkflowTitle] = useState<string>(defaultWorkflowTitle);

    const handleTitleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setWorkflowTitle(event.target.value);
    };

    const [userName, setUserName] = useState("You");

    useEffect(() => {
        const fetchUserName = async () => {
            try {
                const username = await getUserName();
                setUserName(username);
            } catch (error) {
                console.error("Failed to get username: ", error);
            }
        }
        fetchUserName();
    }, []);

    return (
        <div
            className="min-h-screen bg-gray-100"
        >
            <div
                style={{ backgroundColor: selectedActionService?.color }}
                className="flex items-center justify-between lg:px-10 px-2 py-5 border-b-[1px] border-[#dadada]"
            >
                <div>
                    <ActionButton
                        text="Back"
                        textColor="text-white"
                        borderColor="border-white"
                        bgColor={`bg-[${selectedActionService?.color}]`}
                        onClick={() => goBack()}
                    />
                </div>

                <h1
                    className="text-white text-[2.2rem] sm:text-[3.2rem] font-extrabold text-center flex-1"
                >
                    Review and finish
                </h1>

                <div className="w-[100px]"></div>
            </div>

            <div
                style={{ backgroundColor: selectedActionService?.color }}
                className="flex flex-col items-center justify-center lg:px-0 px-3"
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
                        by{' '}
                        <span className="font-bold">
                            {userName}
                        </span>
                    </p>
                </MaxWidthWrapper>
            </div>

            <div
                className="mt-5 flex flex-col !items-center lg:px-0 px-3"
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
