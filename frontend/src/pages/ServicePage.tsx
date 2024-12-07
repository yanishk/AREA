import { useLocation, useNavigate } from "react-router-dom";
import CancelButton from "../components/buttons/CancelButton";

const ServicePage = () => {
    const navigate = useNavigate();
    const location = useLocation();

    const goBack = () => {
        navigate(-1);
    };

    const type = location.state?.type;
    const serviceName = location.state?.name;
    const serviceColor = location.state?.color;
    const serviceIcon = location.state?.iconPath;

    return (
        <div className="min-h-screen bg-gray-100">
            <div style={{ backgroundColor: serviceColor }} className="flex items-center justify-between px-10 py-5 border-b-[1px] border-[#dadada]">
                <div>
                    <CancelButton
                        text="Cancel"
                        onClick={goBack}
                    />
                </div>

                <h1 className="text-white text-[3.2rem] font-extrabold text-center flex-1">
                    Choose a trigger
                </h1>

                <div className="w-[100px]"></div>
            </div>

            <div style={{ backgroundColor: serviceColor }} className="flex flex-col items-center justify-center">
                <img src={serviceIcon} className="w-[122px] h-[122px] mt-14" />
                <h2 className="text-white text-[52px] font-bold mt-6 mb-14">{serviceName}</h2>
            </div>

        </div>
    );
};

export default ServicePage;
