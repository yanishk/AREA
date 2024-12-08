import { useLocation, useNavigate } from "react-router-dom";
import CancelButton from "../components/buttons/CancelButton";
import axios from "axios";
import { useEffect, useState } from "react";

const ServicePage = () => {
    type ElementType = {
        name: string;
        description: string;
    };

    const navigate = useNavigate();
    const location = useLocation();
    const [elements, setElements] = useState<ElementType[]>([]);

    const goBack = () => {
        navigate(-1);
    };

    const type = location.state?.type;
    const serviceName = location.state?.name;
    const serviceColor = location.state?.color;
    const serviceIcon = location.state?.iconPath;

    const getElements = async () => {
        try {
            let result;

            if (type === "action") {
                result = await axios.get(`${import.meta.env.VITE_API_URL}${serviceName}/actions`);
                console.log("API Response Data: ", result.data);
            } else if (type === "reaction") {
                result = await axios.get(`${import.meta.env.VITE_API_URL}${serviceName}/reactions`);
            }

            console.log("API Response: ", result);

            if (result && result.data) {
                const dataKey = type === "action" ? "actions" : "reactions";

                if (Array.isArray(result.data[dataKey])) {
                    const propertiesElements = result.data[dataKey].map((element: any) => ({
                        name: element.Name,
                        description: element.Description,
                    }));

                    console.log("Mapped Elements: ", propertiesElements);
                    setElements(propertiesElements);
                } else {
                    console.error(`Unexpected structure for ${dataKey}:`, result.data[dataKey]);
                }
            } else {
                console.error("Unexpected API response structure:", result);
            }
        } catch (err) {
            console.error("Error fetching services:", err);
        }
    };

    useEffect(() => {
        getElements();
    }, []);

    return (
        <div className="min-h-screen bg-gray-100">
            <div style={{ backgroundColor: serviceColor }} className="flex items-center justify-between px-10 py-5 border-b-[1px] border-[#dadada]">
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
                    Choose a trigger
                </h1>

                <div className="w-[100px]"></div>
            </div>

            <div style={{ backgroundColor: serviceColor }} className="flex flex-col items-center justify-center">
                <img src={serviceIcon} className="w-[122px] h-[122px] mt-14" />
                <h2 className="text-white text-[52px] font-bold mt-6 mb-14">{serviceName}</h2>
            </div>

            <div className="mt-16 ml-[35px]">
                <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
                    {elements.map((element) => (
                        <div
                            key={element.name}
                            style={{ backgroundColor: serviceColor }}
                            className="h-[240px] flex flex-col items-start pl-[30px] pt-2 cursor-pointer rounded-[10px]"
                        >
                            <h2 className="text-white text-lg text-[1.6em] font-[800] leading-[1.125] mt-6">{element.name}</h2>
                            <h2 className="text-white text-lg text-[1.3em] font-[400] mt-6">{element.description}</h2>
                        </div>
                    ))}
                </div>
            </div>

        </div>
    );
};

export default ServicePage;
