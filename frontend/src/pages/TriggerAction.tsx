import { useLocation, useNavigate } from "react-router-dom";
import CancelButton from "../components/buttons/CancelButton";
import { servicesData } from "../ServicesData";
import { useEffect, useState } from "react";
import DynamicFields from "../components/inputs/DynamicFields";
import TriggerButton from "../components/buttons/TriggerButton";

const TriggerActionPage = () => {
    const navigate = useNavigate();
    const location = useLocation();

    const serviceName = location.state?.serviceName;
    const serviceColor = location.state?.serviceColor;
    const serviceIcon = location.state?.serviceIcon;
    const actionName = location.state?.actionName;
    const actionId = location.state?.actionId;
    const selectedReaction = location.state?.selectedReaction;

    const fields = serviceName && actionName &&
        servicesData[serviceName]?.actions?.[actionName]
        ? servicesData[serviceName].actions[actionName]
        : [];

    const [fieldValues, setFieldValues] = useState<Record<string, string>>({});

    const handleFieldChange = (id: string, value: string) => {
        setFieldValues((prev) => ({
            ...prev,
            [id]: value,
        }));
    };

    useEffect(() => {
        const defaultValues = fields.reduce((acc, field) => {
            acc[field.id] = field.options?.[0] ?? "";
            return acc;
        }, {} as Record<string, string>);
        setFieldValues(defaultValues);
    }, [fields]);

    const handleActionClick = () => {
        navigate("/create", {
            state: {
                selectedReaction,
                selectedAction: {
                    name: actionName,
                    id: actionId,
                    color: serviceColor,
                    icon: serviceIcon,
                    fields: fieldValues,
                },
            },
        });
    };

    const goBack = () => {
        navigate(-1);
    };

    return (
        <div style={{ backgroundColor: serviceColor }} className="min-h-screen">
            <div
                className="flex items-center justify-between px-10 py-5 border-b-[1px] border-[#dadada]"
            >
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
                    Complete trigger fields
                </h1>

                <div className="w-[100px]"></div>
            </div>

                <div
                    style={{ backgroundColor: serviceColor }}
                    className="flex flex-col items-center justify-center"
                >
                    <img
                        alt="icon"
                        src={serviceIcon}
                        className="w-[122px] h-[122px] mt-14"
                    />
                    <h2
                        className="text-white text-[52px] font-bold mt-6 mb-8"
                    >
                        {actionName}
                    </h2>
                </div>

                <div
                    className="flex flex-col items-center"
                >
                    <DynamicFields
                        fields={fields}
                        values={fieldValues}
                        onChange={handleFieldChange}
                    />
                    <TriggerButton
                        text="Create Trigger"
                        textColor="text-black"
                        bgColor="bg-white"
                        hoverColor="hover:bg-[#EEEEEE]"
                        onClick={handleActionClick}
                    />
                </div>

        </div>
    );
};

export default TriggerActionPage