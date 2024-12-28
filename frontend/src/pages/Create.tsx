import { useNavigate } from "react-router-dom";
import {useEffect, useState } from "react";
import axios from "axios";
import {servicesData } from "../ServicesData";
import HomeStep from "../components/workflows/HomeStep";
import SelectServiceActionStep from "../components/workflows/ServiceActionStep";
import SelectActionStep from "../components/workflows/SelectActionStep";
import TriggerActionStep from "../components/workflows/TriggerActionStep";

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
        navigate(-1);
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
                    <SelectActionStep
                        selectedService={selectedService}
                        availableActions={availableActions}
                        handleSelectAction={handleSelectAction}
                        goToStep={goToStep}
                    />
                )

            case Steps.TRIGGER_ACTION:
                return (
                    <TriggerActionStep
                        selectedService={selectedService}
                        selectedAction={selectedAction}
                        fields={fields}
                        actionFieldValues={actionFieldValues}
                        handleFieldChange={handleFieldChange}
                        goToStep={goToStep}
                    />
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
