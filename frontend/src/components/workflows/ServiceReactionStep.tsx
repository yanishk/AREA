import CancelButton from "../buttons/CancelButton";

interface ServiceReactionStepProps {
    services: any;
    goBack: () => void;
    handleSelectService: (name: string, color: string, icon: string) => void;
}

const SelectServiceReactionStep: React.FC<ServiceReactionStepProps> = ({
    services,
    goBack,
    handleSelectService
}) => {
    return (
        <div className="min-h-screen bg-gray-100">
            <div className="flex items-center justify-between px-10 py-5">
                <div>
                    <CancelButton
                        text="Cancel"
                        textColor="text-[#222222]"
                        borderColor="border-[#222222]"
                        bgColor="bg-[#F9FAFB]"
                        onClick={() => goBack()}
                    />
                </div>

                <h1 className="text-[3.2rem] font-extrabold text-center flex-1">
                    Choose a service
                </h1>

                <div className="w-[100px]"></div>
            </div>

            <div className="mt-10 px-[140px]">
                <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5">
                    {services.map((service: any) => (
                        <button
                            key={service.name}
                            style={{ backgroundColor: service.color }}
                            className="h-[240px] flex flex-col items-center cursor-pointer"
                            onClick={() =>
                                handleSelectService(
                                    service.name,
                                    service.color,
                                    `${import.meta.env.VITE_HOST_URL}${service.iconPath}`
                                )
                            }
                        >
                            <img alt="icon" src={`${import.meta.env.VITE_HOST_URL}${service.iconPath}`} className="size-28 mt-9"/>
                            <h2 className="text-white text-lg text-[20.5px] font-bold mt-6">{service.name}</h2>
                        </button>
                    ))}
                </div>
            </div>
        </div>
    );
};

export default SelectServiceReactionStep;
