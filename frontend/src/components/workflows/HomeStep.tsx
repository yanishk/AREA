import ButtonIfCustom from "../buttons/ButtonIfCustom";
import ButtonIfDefault from "../buttons/ButtonIfDefault";
import ButtonThenCustom from "../buttons/ButtonThenCustom";
import ButtonThenDefault from "../buttons/ButtonThenDefault";
import CancelButton from "../buttons/CancelButton";
import TriggerButton from "../buttons/TriggerButton";

interface HomeStepProps {
    selectedAction: any | null;
    selectedReaction: any | null;
    selectedServiceAction: { color: string; iconPath: string } | null;
    selectedServiceReaction: { color: string; iconPath: string } | null;
    goToStep: (step: number) => void;
    goHome: () => void;
    handleOnEditAction: () => void;
    handleOnDeleteAction: () => void;
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
    handleContinue,
}) => {
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
                        onEdit={handleOnEditAction}
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
                        onEdit={handleOnEditAction}
                        onDelete={handleOnDeleteAction}
                    />
                ) : (
                    <ButtonThenDefault
                        onClick={() =>
                            goToStep(4)
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

export default HomeStep;
