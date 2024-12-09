import { useLocation, useNavigate } from "react-router-dom";
import CancelButton from "../components/buttons/CancelButton";
import TriggerButton from "../components/buttons/TriggerButton";
import { servicesData } from "../ServicesData";
import axios from "axios";

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
    const selectedAction = location.state?.selectedAction;

    const triggerReactionValue =
        servicesData[serviceName]?.reactions?.[reactionName]?.[0]?.triggerReaction;

    const handleConnect = async () => {
        try {
            const result = await axios.get(
                `${import.meta.env.VITE_API_URL}authentication/${serviceName}`
            );

            const authUrl = result.data["auth-url"];
            if (!authUrl) {
                console.error("Error: Authentication URL not received");
                return;
            }

            console.log("Redirecting to authentication URL:", authUrl);

            const authData = {
                selectedAction,
                selectedReaction: {
                    name: reactionName,
                    color: serviceColor,
                    icon: serviceIcon,
                },
            };
            localStorage.setItem("workflowData", JSON.stringify(authData));

            window.location.href = authUrl;
        } catch (error) {
            console.error("Error during authentication:", error);
        }
    };

    return (
        <div
            style={{ backgroundColor: serviceColor }}
            className="min-h-screen bg-gray-100"
        >
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

            <div
                style={{ backgroundColor: serviceColor }}
                className="flex flex-col items-center justify-center"
            >
                <img src={serviceIcon} className="w-[122px] h-[122px] mt-14" />
                <h2 className="text-white text-[52px] font-bold mt-6 mb-14">
                    {serviceName}
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

export default ConnectService;
