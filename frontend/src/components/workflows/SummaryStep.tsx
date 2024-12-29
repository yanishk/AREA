import { useState } from "react";
import CancelButton from "../buttons/CancelButton";
import TriggerButton from "../buttons/TriggerButton";
import InputField from "../inputs/InputsField";

interface SummaryStepProps {
    selectedActionService: {name: string; color: string; iconPath: string } | null;
    selectedReactionService: {name: string; color: string; iconPath: string } | null;
    selectedAction: {id: string; name: string;} | null;
    selectedReaction: {id: string; name: string;} | null;
    actionFields: Record<string, string>;
    reactionFields: Record<string, string>;
    goBack: () => void;
    handleContinue: () => void;
}

const SummaryStep: React.FC<SummaryStepProps> = ({
    selectedActionService,
    selectedReactionService,
    selectedAction,
    selectedReaction,
    actionFields,
    reactionFields,
    goBack,
    handleContinue
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
                <h2
                    className="text-white text-[52px] font-bold mt-6 mb-14"
                >
                    Workflow Title
                </h2>
                <p
                    className="text-white"
                >
                    <InputField
                        type="text"
                        value={workflowTitle}
                        onChange={handleTitleChange}
                        disableFocus={true}
                    />
                </p>
            </div>

            <div
                className="mt-16 flex flex-col !items-center "
            >
                <div className="w-full max-w-[300px]">
                    <TriggerButton
                        text="Finish"
                        textColor="text-white"
                        bgColor="bg-[#222222]"
                        hoverColor="hover:bg-[#333333]"
                        onClick={handleContinue}
                    />
                </div>
            </div>
        </div>
    );
};

export default SummaryStep;
