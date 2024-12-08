import { useLocation, useNavigate } from "react-router-dom";
import CancelButton from "../components/buttons/CancelButton";
import TriggerButton from "../components/buttons/TriggerButton";

const ConnectService = () => {
    const navigate = useNavigate();
    const location = useLocation();

    const goBack = () => {
        navigate(-1);
    };
    const serviceName = location.state?.serviceName;
    const reactionName = location.state?.reactionName;
    const serviceColor = location.state?.serviceColor;
    const serviceIcon = location.state?.serviceIcon;

    const handleConnect = () => {
        return;
    }

    return (
        <div style={{ backgroundColor: serviceColor }} className="min-h-screen bg-gray-100">
            <div className="flex items-center justify-between px-10 py-5 border-b-[1px] border-[#dadada]">
                <div>
                <CancelButton
                        text="Back"
                        textColor="text-white"
                        borderColor="border-white"
                        bgColor={`bg-[${serviceColor}]`}
                        onClick={goBack}
                    />
                </div>

                <h1 className="text-white text-[3.2rem] font-extrabold text-center flex-1">
                    Connect Service
                </h1>

                <div className="w-[100px]"></div>
            </div>

            <div style={{ backgroundColor: serviceColor }} className="flex flex-col items-center justify-center">
                <img src={serviceIcon} className="w-[122px] h-[122px] mt-14" />
                <h2 className="text-white text-[52px] font-bold mt-6 mb-14">{serviceName}</h2>
            </div>
            <div className="cursor-pointer flex justify-center items-center">
                <TriggerButton
                    onClick={handleConnect}
                    text="Connect"
                />
            </div>

        </div>
    );
};

export default ConnectService;
