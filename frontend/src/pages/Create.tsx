import { useNavigate, useLocation, useSearchParams } from "react-router-dom";
import ButtonIfCustom from "../components/buttons/ButtonIfCustom";
import { ElementType, useEffect, useState } from "react";
import ButtonIfDefault from "../components/buttons/ButtonIfDefault";
import CancelButton from "../components/buttons/CancelButton";
import ButtonThenCustom from "../components/buttons/ButtonThenCustom";
import ButtonThenDefault from "../components/buttons/ButtonThenDefault";
import TriggerButton from "../components/buttons/TriggerButton";
import axios from "axios";
import DynamicFields from "../components/inputs/DynamicFields";
import { ServiceDataType, servicesData } from "../ServicesData";

enum Steps {
    HOME,
    SELECT_SERVICE_ACTION,
    SELECT_ACTION,
    TRIGGER_ACTION,
    SELECT_SERVICE_REACTION,
    CONNECT_SERVICE,
    TRIGGER_REACTION,
    SUMMARY,
}

type availableActionsType = {
    name: string;
    description: string;
    id: string;
};

type selectedServiceType = {
    name: string;
    color: string;
    iconPath: string;
};

type selectedActionType = {
    name: string;
    id: string;
};

const WorkflowBuilder = () => {
    const navigate = useNavigate();

    const [selectedService, setSelectedService] = useState<selectedServiceType | null>(null);
    const [availableActions, setAvailableActions] = useState<availableActionsType[]>([]);
    const [selectedAction, setSelectedAction] = useState<selectedActionType | null>(null);
    const [selectedReaction, setSelectedReaction] = useState(null);
    const [services, setServices] = useState([]);

    const [currentStep, setCurrentStep] = useState<Steps>(Steps.HOME);

    const goToStep = (step: Steps) => {
        setCurrentStep(step);
    };

    const handleOnEditAction = () => {
        return;
    }

    const handleDeleteAction = () => {
        return;
    }

    const goHome = () => {
        return;
    }

    const handleContinue = () => {
        return;
    }

    const getServicesWithActions = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}actions-services`, {
                    withCredentials: true,
            });

            if (result) {
                const propertiesServices = result.data.map((service: any) => ({
                    name: service.name,
                    color: service.color,
                    iconPath: service.logo,
                }));
                setServices(propertiesServices);
            }
        } catch (err) {
            console.error("Error fetching services:", err);
        }
    };

    const getAvailableActions = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}${selectedService?.name}/actions`, {
                    withCredentials: true,
                });

            if (result) {
                const dataKey = "actions";

                if (Array.isArray(result.data[dataKey])) {
                    const actions = result.data[dataKey].map((element: any) => ({
                        name: element.name,
                        description: element.description,
                        id: element.id,
                    }));
                    setAvailableActions(actions);
                } else {
                    console.error(`Error for ${dataKey}:`, result.data[dataKey]);
                }
            } else {
                console.error("Unexpected API response", result);
            }
        }
        catch (err) {
            console.error("Error fetching services:", err);
        }
    };

    useEffect(() => {
        if (currentStep === Steps.SELECT_SERVICE_ACTION) {
            getServicesWithActions();
        }
        if (currentStep === Steps.SELECT_ACTION) {
            getAvailableActions();
        }
    }, [currentStep]);

    const handleSelectService = (name: string, color: string, iconPath: string) => {
        setSelectedService({ name, color, iconPath });
        goToStep(Steps.SELECT_ACTION);
    };

    const handleSelectAction = (name: string, id: string) => {
        setSelectedAction({name, id});
        goToStep(Steps.TRIGGER_ACTION);
    }

    const serviceName = selectedService?.name;
    const actionName = selectedAction?.name;

    const fields = serviceName && actionName &&
    servicesData[serviceName]?.actions?.[actionName]
    ? servicesData[serviceName].actions[actionName]
    : [];

    const [actionFieldValues, setActionFieldValues] = useState<Record<string, string>>({});

    const handleFieldChange = (id: string, value: string) => {
        setActionFieldValues((prev) => ({
            ...prev,
            [id]: value,
        }));
    };

    useEffect(() => {
        const defaultValues = fields.reduce((acc, field) => {
            acc[field.id] = field.options?.[0] ?? "";
            return acc;
        }, {} as Record<string, string>);
        setActionFieldValues(defaultValues);
    }, [fields]);

    const renderStep = () => {
        switch (currentStep) {
            case Steps.HOME:
                return (
                    <div>
                        <div
                            className="flex items-center justify-between px-10 py-5"
                        >
                            <div>
                                <CancelButton
                                    text="Cancel"
                                    textColor="text-[#222222]"
                                    borderColor="border-[#222222]"
                                    bgColor="bg-[#F9FAFB]"
                                    onClick={goHome}
                                />
                            </div>

                            <h1
                                className="text-[3.2rem] font-[900] text-center flex-1"
                            >
                                Create
                            </h1>

                            <div className="w-[140px]"></div>
                        </div>

                        <div
                            className="flex items-center justify-center flex-col px-10 pt-[100px] py-5"
                        >
                            <div className="pt-10">
                                {selectedAction ? (
                                    <ButtonIfCustom
                                        text={selectedAction.name}
                                        color={selectedService?.color || ""}
                                        iconPath={selectedService?.iconPath || ""}
                                        onEdit={handleOnEditAction}
                                        onDelete={handleDeleteAction}
                                    />
                                ) : (
                                    <ButtonIfDefault
                                        onClick={() =>
                                            goToStep(Steps.SELECT_SERVICE_ACTION)
                                        }
                                        text="If This"
                                    />
                                )}
                            </div>

                            <div className="h-[70px] w-[5px] bg-[#EEEEEE]" />

                            <div className="mb-[50px]">
                                {selectedReaction ? (
                                    <ButtonThenDefault
                                    onClick={() =>
                                        goHome()
                                    }
                                    text="Then That"
                                    disabled={!selectedAction}
                                />
                                ) : (
                                    <ButtonThenDefault
                                        onClick={() =>
                                            goHome()
                                        }
                                        text="Then That"
                                        disabled={!selectedAction}
                                    />
                                )}
                            </div>

                            <div
                                className="w-full max-w-[350px] flex flex-col items-center justify-center"
                            >
                                {selectedAction && selectedReaction ? (
                                    <TriggerButton
                                        text="Create Workflow"
                                        textColor="text-white"
                                        bgColor="bg-[#222222]"
                                        hoverColor="hover:bg-[#333333]"
                                        onClick={handleContinue}
                                    />
                                ) : (
                                    <></>
                                )}
                            </div>
                        </div>
                    </div>
                )

            case Steps.SELECT_SERVICE_ACTION:
                return (
                    <div className="min-h-screen bg-gray-100">
                        <div className="flex items-center justify-between px-10 py-5">
                            <div>
                                <CancelButton
                                    text="Cancel"
                                    textColor="text-[#222222]"
                                    borderColor="border-[#222222]"
                                    bgColor="bg-[#F9FAFB]"
                                    onClick={() => goToStep(Steps.HOME)}
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
                                        key={service.id}
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

                case Steps.SELECT_ACTION:
                    return (
                        <div className="min-h-screen bg-gray-100">
                            {selectedService && (
                                <>
                                    <div
                                        style={{ backgroundColor: selectedService.color }}
                                        className="flex items-center justify-between px-10 py-5 border-b-[1px] border-[#dadada]"
                                    >
                                        <div>
                                            <CancelButton
                                                text="Back"
                                                textColor="text-white"
                                                borderColor="border-white"
                                                bgColor={`bg-[${selectedService.color}]`}
                                                onClick={() => goToStep(Steps.SELECT_SERVICE_ACTION)}
                                            />
                                        </div>

                                        <h1 className="text-white text-[3.2rem] font-extrabold text-center flex-1">
                                            Choose a trigger
                                        </h1>

                                        <div className="w-[100px]"></div>
                                    </div>

                                    <div style={{ backgroundColor: selectedService.color }} className="flex flex-col items-center justify-center">
                                        <img alt="icon" src={selectedService.iconPath} className="w-[122px] h-[122px] mt-14" />
                                        <h2 className="text-white text-[52px] font-bold mt-6 mb-14">{selectedService.name}</h2>
                                    </div>

                                    <div className="mt-16 ml-[35px]">
                                        <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4">
                                            {availableActions.map((element) => (
                                                <button
                                                    key={element.name}
                                                    style={{ backgroundColor: selectedService.color}}
                                                    className="h-[240px] flex flex-col items-start pt-2 cursor-pointer rounded-[10px]"
                                                    // onClick={() => goToStep(Steps.TRIGGER_ACTION)}
                                                    onClick={() => handleSelectAction(
                                                        element.name,
                                                        element.id
                                                    )
                                                }
                                                >
                                                    <div className="text-left w-full pl-[30px]">
                                                        <h2 className="text-white text-lg text-[1.6em] font-[800] leading-[1.125] mt-6 m-0 p-0">{element.name}</h2>
                                                        <h2 className="text-white text-lg text-[1.3em] font-[400] mt-6 m-0">{element.description}</h2>
                                                    </div>
                                                </button>
                                            ))}
                                        </div>
                                    </div>
                                </>
                            )}
                        </div>
                    )

                    case Steps.TRIGGER_ACTION:
                        return (
                            <div style={{ backgroundColor: selectedService?.color }} className="min-h-screen">
                                <div
                                    className="flex items-center justify-between px-10 py-5 border-b-[1px] border-[#dadada]"
                                >
                                    <div>
                                        <CancelButton
                                            text="Back"
                                            textColor="text-white"
                                            borderColor="border-white"
                                            bgColor={`bg-[${selectedService?.color}]`}
                                            onClick={() => goToStep(Steps.SELECT_ACTION)}
                                        />
                                    </div>

                                    <h1 className="text-white text-[3.2rem] font-extrabold text-center flex-1">
                                        Complete trigger fields
                                    </h1>

                                    <div className="w-[100px]"></div>
                                </div>

                                    <div
                                        style={{ backgroundColor: selectedService?.color }}
                                        className="flex flex-col items-center justify-center"
                                    >
                                        <img
                                            alt="icon"
                                            src={selectedService?.iconPath}
                                            className="w-[122px] h-[122px] mt-14"
                                        />
                                        <h2
                                            className="text-white text-[52px] font-bold mt-6 mb-8"
                                        >
                                            {selectedAction?.name}
                                        </h2>
                                    </div>

                                    <div
                                        className="flex flex-col items-center"
                                    >
                                        <DynamicFields
                                            fields={fields}
                                            values={actionFieldValues}
                                            onChange={handleFieldChange}
                                        />
                                        <TriggerButton
                                            text="Create Trigger"
                                            textColor="text-black"
                                            bgColor="bg-white"
                                            hoverColor="hover:bg-[#EEEEEE]"
                                            onClick={() => goToStep(Steps.HOME)}
                                        />
                                    </div>

                            </div>
                        );
            }
    }

    return (
        <>
            {renderStep()}
        </>
    );
};

export default WorkflowBuilder
