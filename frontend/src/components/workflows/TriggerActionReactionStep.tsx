import { useEffect, useState } from "react";
import { Parameters } from "../../utils/ActionReactionParameters";
import Steps from "../../utils/CreateStepsType";
import { fieldType } from "../../utils/fieldType";
import getFieldValues from "../../utils/getFieldValues";
import CancelButton from "../buttons/CancelButton";
import TriggerButton from "../buttons/TriggerButton";
import DynamicFields from "../inputs/DynamicFields";

interface TriggerActionReactionProps {
    selectedService: { name: string; color: string; iconPath: string } | null;
    selectedItem: {name: string, id: string, parameters: Parameters[]} | null;
    itemFieldValues: Record<string, string>;
    handleFieldChange: (id: string, value: string, type: "action" | "reaction") => void;
    goToStep: (step: number) => void;
    goBack: () => void;
    type: "action" | "reaction";
}

const TriggerActionReactionStep: React.FC<TriggerActionReactionProps> = ({
    selectedService,
    selectedItem,
    itemFieldValues,
    handleFieldChange,
    goToStep,
    goBack,
    type,
}) => {

    const getFieldType = (param: Parameters) => {
        if (!param.isexhaustive && param.route) {
            return fieldType.SELECT;
        }
        if (param.isexhaustive) {
            return fieldType.SELECT;
        }
        if (!param.isexhaustive && param.values.length === 1) {
            return fieldType.SELECT;
        }

        if (!param.isexhaustive && param.route === null) {
            return fieldType.INPUT;
        }

        return fieldType.INPUT;
    };

    const [fields, setFields] = useState<{
        id: string;
        label: string;
        type: fieldType;
        options?: string[]
        }[]
    > ([]);

    useEffect(() => {
        const fetchFields = async () => {
            if (!selectedItem) return;

            const fetchedFields = await Promise.all(
                selectedItem.parameters.map(async (param) => {
                    const options =
                        getFieldType(param) === fieldType.SELECT
                            ? await getFieldValues(param)
                            : undefined;

                    return {
                        id: param.name,
                        label: param.name,
                        type: getFieldType(param),
                        options: options || [],
                    };
                })
            );

            setFields(fetchedFields);
        };

        fetchFields();
    }, [selectedItem]);

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

