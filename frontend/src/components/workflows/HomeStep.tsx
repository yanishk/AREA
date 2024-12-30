import { servicesData } from "../../ServicesData";
import ButtonIfCustom from "../buttons/ButtonIfCustom";
import ButtonIfDefault from "../buttons/ButtonIfDefault";
import ButtonThenCustom from "../buttons/ButtonThenCustom";
import ButtonThenDefault from "../buttons/ButtonThenDefault";
import CancelButton from "../buttons/CancelButton";
import TriggerButton from "../buttons/TriggerButton";

interface Action {
    name: string;
    id: string;
}

interface Reaction {
    name: string;
    id: string;
}
interface HomeStepProps {
    selectedAction: Action | null;
    selectedReaction: Reaction | null;
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

    const isTriggerable = (serviceName: string, actionOrReactionName: string, type: "action" | "reaction"): boolean => {
        const service = servicesData[serviceName];

        if (!service) {
            return false;
        }

        if (type === "action" && service.actions) {
            const action = service.actions[actionOrReactionName];

            if (!action) {
                return false;
            }

            return action.some((actionData) => actionData.trigger === "yes");
        }

        if (type === "reaction" && service.reactions) {
            const reaction = service.reactions[actionOrReactionName];

            if (!reaction) {
                return false;
            }

            return reaction.some((reactionData) => reactionData.trigger === "yes");
        }

        return false;
    };

    const isActionTriggerable =
        selectedServiceAction && selectedAction
            ? isTriggerable(selectedServiceAction.name, selectedAction.name, "action")
            : false;

    const isReactionTriggerable =
        selectedServiceReaction && selectedReaction
            ? isTriggerable(selectedServiceReaction.name, selectedReaction.name, "reaction")
            : false;

    return (
        <div>
        <div
            className="flex items-center justify-between px-10 py-5"
        >
            <div>
                <CancelButton
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
            <div className="pt-10">
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
                            goToStep(1)
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
                        color={selectedServiceReaction?.color || ""}
                        iconPath={selectedServiceReaction?.iconPath || ""}
                        onEdit={handleOnEditReaction}
                        isEdit={isReactionTriggerable}
                        onDelete={handleOnDeleteReaction}
                    />
                ) : (
                    <ButtonThenDefault
                        onClick={() =>
                            goToStep(5)
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
