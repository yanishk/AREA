import ActionButton from "../buttons/ActionButton";

interface SelectServiceStepProps {
    services: any;
    goBack: () => void;
    handleSelectService: (name: string, color: string, icon: string, isauthneeded: boolean) => void;
}

const SelectServiceStep: React.FC<SelectServiceStepProps> = ({
    services,
    goBack,
    handleSelectService
}) => {
    return (
        <div className="min-h-screen bg-gray-100">
            <div className="flex w-[110%] items-center justify-between lg:px-10 px-2 py-5">
                <div>
                    <ActionButton
                        text="Cancel"
                        textColor="text-[#222222]"
                        borderColor="border-[#222222]"
                        bgColor="bg-[#F9FAFB]"
                        onClick={() => goBack()}
                    />
                </div>

                <h1 className="text-[1.9rem] sm:text-[3.2rem] font-extrabold text-center w-full">
                    Choose a service
                </h1>

                <div className="w-[100px]"></div>
            </div>

            <div className="mt-10 lg:px-36 md:px-16 px-10">
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
                                    `${import.meta.env.VITE_HOST_URL}${service.iconPath}`,
                                    service.isauthneeded
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

export default SelectServiceStep;
