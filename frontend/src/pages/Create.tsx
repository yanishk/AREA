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
import HomeStep from "../components/workflows/HomeStep";
import SelectServiceActionStep from "../components/workflows/ServiceActionStep";

enum Steps {
    HOME,
    SELECT_SERVICE_ACTION,
    SELECT_ACTION,
    TRIGGER_ACTION,
    SELECT_SERVICE_REACTION,
    CONNECT_SERVICE,
    SELECT_REACTION,
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
                    <HomeStep
                        selectedAction={selectedAction}
                        selectedReaction={selectedReaction}
                        selectedService={selectedService}
                        goToStep={goToStep}
                        goHome={goHome}
                        handleOnEditAction={handleOnEditAction}
                        handleOnDeleteAction={handleDeleteAction}
                        handleContinue={handleContinue}
                    />
                )

            case Steps.SELECT_SERVICE_ACTION:
                return (
                    <SelectServiceActionStep
                        services={services}
                        handleSelectService={handleSelectService}
                        goToStep={goToStep}
                    />
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
