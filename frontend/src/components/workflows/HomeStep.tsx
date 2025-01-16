import { ActionType } from "../../utils/ActionType";
import Steps from "../../utils/CreateStepsType";
import { ReactionType } from "../../utils/ReactionType";
import ButtonIfCustom from "../buttons/ButtonIfCustom";
import ButtonIfDefault from "../buttons/ButtonIfDefault";
import ButtonThenCustom from "../buttons/ButtonThenCustom";
import ButtonThenDefault from "../buttons/ButtonThenDefault";
import ActionButton from "../buttons/ActionButton";
import TriggerButton from "../buttons/TriggerButton";

interface HomeStepProps {
    selectedAction: ActionType | null;
    selectedReaction: ReactionType | null;
    selectedServiceAction: { name: string; color: string; iconPath: string } | null;
    selectedServiceReaction: { name: string; color: string; iconPath: string } | null;
    goToStep: (step: number) => void;
    goHome: () => void;
    handleOnEditAction: () => void;
    handleOnDeleteAction: () => void;
    handleOnEditReaction: () => void;
    handleOnDeleteReaction: () => void;
    handleContinue: () => void;
}

const HomeStep: React.FC<HomeStepProps> = ({
    selectedAction,
    selectedReaction,
    selectedServiceAction,
    selectedServiceReaction,
    goToStep,
    goHome,
    handleOnEditAction,
    handleOnDeleteAction,
    handleOnEditReaction,
    handleOnDeleteReaction,
    handleContinue,
}) => {

    const isTriggerable = (type: "action" | "reaction"): boolean => {
        if (type === "action") {
            if (selectedAction?.parameters != null) {
                return true;
            }
        }

        if (type === "reaction") {
            if (selectedReaction?.parameters != null) {
                return true;
            }
        }
        return false;
    };

    const isActionTriggerable =
        selectedServiceAction && selectedAction
            ? isTriggerable("action")
            : false;

    const isReactionTriggerable =
        selectedServiceReaction && selectedReaction
            ? isTriggerable("reaction")
            : false;

    return (
        <div>
        <div
            className="flex items-center justify-between lg:px-10 px-2 py-5"
        >
            <div>
                <ActionButton
                    text="Cancel"
                    textColor="text-[#222222]"
                    borderColor="border-[#222222]"
                    bgColor="bg-[#F9FAFB]"
                    onClick={goHome}
                />
            </div>

            <h1
                className="text-[3.2rem] font-[900] text-center flex-1"
            >
                Create
            </h1>

            <div className="w-[140px]"></div>
        </div>

        <div
            className="flex items-center justify-center flex-col px-10 pt-[100px] py-5"
        >
            <div className="pt-10 max-w-[730px] w-full">
                {selectedAction ? (
                    <ButtonIfCustom
                        text={selectedAction.name}
                        color={selectedServiceAction?.color || ""}
                        iconPath={selectedServiceAction?.iconPath || ""}
                        onEdit={isActionTriggerable ? handleOnEditAction : () => {}}
                        isEdit={isActionTriggerable}
                        onDelete={handleOnDeleteAction}
                    />
                ) : (
                    <ButtonIfDefault
                        onClick={() =>
                            goToStep(Steps.SELECT_SERVICE_ACTION)
                        }
                        text="If This"
                    />
                )}
            </div>

            <div className="h-[70px] w-[5px] bg-[#EEEEEE]" />

            <div className="mb-[50px] max-w-[730px] w-full">
                {selectedReaction ? (
                    <ButtonThenCustom
                        text={selectedReaction.name}
                        color={selectedServiceReaction?.color || ""}
                        iconPath={selectedServiceReaction?.iconPath || ""}
                        onEdit={handleOnEditReaction}
                        isEdit={isReactionTriggerable}
                        onDelete={handleOnDeleteReaction}
                    />
                ) : (
                    <ButtonThenDefault
                        onClick={() =>
                            goToStep(Steps.SELECT_SERVICE_REACTION)
                        }
                        text="Then That"
                        isDisabled={!selectedAction}
                    />
                )}
            </div>

            <div
                className="w-full max-w-[350px] flex flex-col items-center justify-center"
            >
                {selectedAction && selectedReaction ? (
                    <TriggerButton
                        text="Continue"
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

export default HomeStep;
