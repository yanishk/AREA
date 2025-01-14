import ActionButton from "../buttons/ActionButton";
import TriggerButton from "../buttons/TriggerButton";

interface ConnectServiceProps {
    selectedService: {name: string; color: string; iconPath: string } | null;
    goBack: () => void;
    handleConnect: () => void;
}

const ConnectServiceStep: React.FC<ConnectServiceProps> = ({
    selectedService,
    goBack,
    handleConnect
}) => {
    return (
        <div
            style={{ backgroundColor: selectedService?.color }}
            className="min-h-screen bg-gray-100"
        >
            <div className="flex items-center justify-between px-10 py-5 border-b-[1px] border-[#dadada]">
                <div>
                    <ActionButton
                        text="Back"
                        textColor="text-white"
                        borderColor="border-white"
                        bgColor={`bg-[${selectedService?.color}]`}
                        onClick={goBack}
                    />
                </div>

                <h1 className="text-white text-[3.2rem] font-extrabold text-center flex-1">
                    Connect Service
                </h1>

                <div className="w-[100px]"></div>
            </div>

            <div
                style={{ backgroundColor: selectedService?.color }}
                className="flex flex-col items-center justify-center"
            >
                <img alt="icon" src={selectedService?.iconPath} className="w-[122px] h-[122px] mt-14" />
                <h2 className="text-white text-[52px] font-bold mt-6 mb-14">
                    {}
                </h2>
            </div>

            <div className="cursor-pointer flex justify-center items-center">
                <TriggerButton
                    text={"Connect"}
                    textColor="text-black"
                    bgColor="bg-white"
                    hoverColor="hover:bg-[#EEEEEE]"
                    onClick={handleConnect}
                />
            </div>
        </div>
    );
};

export default ConnectServiceStep;
