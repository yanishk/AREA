import ActionButton from "../buttons/ActionButton";

type Parameters = {
    isexhaustive: boolean;
    name: string;
    type: string;
    route: string | null;
    values: string[];
};

interface SelectActionReactionStepProps {
    selectedService: { name: string; color: string; iconPath: string } | null;
    availableItems: { name: string; description: string; id: string; nbparam: number, parameters: Parameters[] } [] | null;
    handleSelect: (name: string, id: string, nbparam: number, parameters: Parameters[]) => void;
    goBack: () => void;
    type: "action" | "reaction";
}

const SelectActionReactionStep: React.FC<SelectActionReactionStepProps> = ({
    selectedService,
    availableItems,
    handleSelect,
    goBack,
    type,
}) => {
    return (
        <div className="min-h-screen bg-gray-100">
        {selectedService && (
            <>
                <div
                    style={{ backgroundColor: selectedService.color }}
                    className="flex items-center justify-between lg:px-10 px-2 py-5 border-b-[1px] border-[#dadada] w-full"
                >
                    <div>
                        <ActionButton
                            text="Back"
                            textColor="text-white"
                            borderColor="border-white"
                            bgColor={`bg-[${selectedService.color}]`}
                            onClick={() => goBack()}
                        />
                    </div>

                    <h1 className="text-white text-[1.4rem] sm:text-[3.2rem] font-extrabold text-center w-full">
                        {type === "action" ? "Choose an action" : "Choose a reaction"}
                    </h1>

                    <div className="w-[100px]"></div>
                </div>

                <div style={{ backgroundColor: selectedService.color }} className="flex flex-col items-center justify-center">
                    <img alt="icon" src={selectedService.iconPath} className="w-[122px] h-[122px] mt-14" />
                    <h2 className="text-white text-[52px] font-bold mt-6 mb-14">{selectedService.name}</h2>
                </div>

                <div className="mt-16 ml-[35px]">
                    <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
                        {availableItems?.map((element) => (
                            <button
                                key={element.name}
                                style={{ backgroundColor: selectedService.color}}
                                className="h-[240px] flex flex-col items-start pt-2 cursor-pointer rounded-[10px] mr-6 mb-5"
                                onClick={() => handleSelect(
                                    element.name,
                                    element.id,
                                    element.nbparam,
                                    element.parameters
                                )
                            }
                            >
                                <div className="text-left w-full pl-[30px]">
                                    <h2 className="text-white text-lg text-[1.6em] font-[800] leading-[1.125] mt-6 m-0 p-0">{element.name}</h2>
                                    <h2 className="text-white text-lg text-[1.3em] font-[400] mt-6 mr-4">{element.description}</h2>
                                </div>
                            </button>
                        ))}
                    </div>
                </div>
            </>
        )}
    </div>
    );
};

export default SelectActionReactionStep;
