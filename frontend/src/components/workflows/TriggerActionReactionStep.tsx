import Steps from "../../CreateSteps";
import CancelButton from "../buttons/CancelButton";
import TriggerButton from "../buttons/TriggerButton";
import DynamicFields from "../inputs/DynamicFields";

interface TriggerActionReactionProps {
    selectedService: { name: string; color: string; iconPath: string } | null;
    selectedItem: {name: string, id: string} | null;
    fields: any[];
    itemFieldValues: Record<string, string>;
    handleFieldChange: (id: string, value: string, type: "action" | "reaction") => void;
    goToStep: (step: number) => void;
    goBack: () => void;
    type: "action" | "reaction";
}

const TriggerActionReactionStep: React.FC<TriggerActionReactionProps> = ({
    selectedService,
    selectedItem,
    fields,
    itemFieldValues,
    handleFieldChange,
    goToStep,
    goBack,
    type,
}) => {
    return (
        <div style={{ backgroundColor: selectedService?.color }} className="min-h-screen">
            <div
                className="flex items-center justify-between px-10 py-5 border-b-[1px] border-[#dadada]"
            >
                <div>
                    <CancelButton
                        text="Back"
                        textColor="text-white"
                        borderColor="border-white"
                        bgColor={`bg-[${selectedService?.color}]`}
                        onClick={() => goBack()}
                    />
                </div>

                <h1 className="text-white text-[3.2rem] font-extrabold text-center flex-1">
                    Complete trigger fields
                </h1>

                <div className="w-[100px]"></div>
            </div>

                <div
                    style={{ backgroundColor: selectedService?.color }}
                    className="flex flex-col items-center justify-center"
                >
                    <img
                        alt="icon"
                        src={selectedService?.iconPath}
                        className="w-[122px] h-[122px] mt-14"
                    />
                    <h2
                        className="text-white text-[52px] font-bold mt-6 mb-8"
                    >
                        {selectedItem?.name}
                    </h2>
                </div>

                <div
                    className="flex flex-col items-center"
                >
                    <DynamicFields
                        fields={fields}
                        values={itemFieldValues}
                        onChange={(id, value) => handleFieldChange(id, value, type)}
                    />
                    <TriggerButton
                        text="Create Trigger"
                        textColor="text-black"
                        bgColor="bg-white"
                        hoverColor="hover:bg-[#EEEEEE]"
                        onClick={() => goToStep(Steps.HOME)}
                    />
                </div>

        </div>
    );
};

export default TriggerActionReactionStep;
