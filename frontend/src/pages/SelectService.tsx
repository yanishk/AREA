import { useLocation, useNavigate } from "react-router-dom";
import CancelButton from "../components/buttons/CancelButton";
import axios from "axios";
import { useEffect, useState } from "react";

const SelectServicePage = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const [services, setServices] = useState([]);

    const goBack = () => {
        navigate(-1);
    };

    const type = location.state?.type;
    const selectedAction = location.state?.selectedAction;
    const selectedReaction = location.state?.selectedReaction;

    const getServices = async () => {
        try {
            let result;

            if (type === "action") {
                result = await axios.get(`${import.meta.env.VITE_API_URL}actions-services`, {
                    withCredentials: true,
                });
            } else if (type === "reaction") {
                result = await axios.get(`${import.meta.env.VITE_API_URL}reactions-services`, {
                    withCredentials: true,
                });
            }

            if (result) {
                const propertiesServices = result.data.map((service: any) => ({
                    name: service.Name,
                    color: service.Color,
                    iconPath: service.Logo,
                }));
                setServices(propertiesServices);
            }
        } catch (err) {
            console.error("Error fetching services:", err);
        }
    };

    useEffect(() => {
        getServices();
    }, []);

    const handleSelect = (name: string, color: string, iconPath: string) => {
        navigate(`/select-service/${encodeURIComponent(name)}`, {
            state: {
                selectedReaction,
                selectedAction,
                type,
                name,
                color,
                iconPath,
            }
        });
    };

    return (
        <div className="min-h-screen bg-gray-100">
            <div className="flex items-center justify-between px-10 py-5">
                <div>
                    <CancelButton
                        text="Cancel"
                        textColor="text-[#222222]"
                        borderColor="border-[#222222]"
                        bgColor="bg-[#F9FAFB]"
                        onClick={goBack}
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
                        <div
                            key={service.name}
                            style={{ backgroundColor: service.color }}
                            className="h-[240px] flex flex-col items-center cursor-pointer"
                            onClick={() => handleSelect(
                                service.name,
                                service.color,
                                `${import.meta.env.VITE_HOST_URL}${service.iconPath}`)
                            }
                        >
                            <img alt="icon" src={`${import.meta.env.VITE_HOST_URL}${service.iconPath}`} className="size-28 mt-9" />
                            <h2 className="text-white text-lg text-[20.5px] font-bold mt-6">{service.name}</h2>
                        </div>
                    ))}
                </div>
            </div>

        </div>
    );
};

export default SelectServicePage;
