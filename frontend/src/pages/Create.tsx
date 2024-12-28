import { useNavigate } from "react-router-dom";
import {useEffect, useState } from "react";
import axios from "axios";
import {servicesData } from "../ServicesData";
import HomeStep from "../components/workflows/HomeStep";
import SelectServiceActionStep from "../components/workflows/ServiceActionStep";
import SelectActionStep from "../components/workflows/SelectActionStep";
import TriggerActionStep from "../components/workflows/TriggerActionStep";
import SelectServiceReactionStep from "../components/workflows/ServiceReactionStep";

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

    const [selectedActionService, setSelectedActionService] = useState<selectedServiceType | null>(null);
    const [selectedReactionService, setSelectedReactionService] = useState<selectedServiceType | null>(null);
    const [availableActions, setAvailableActions] = useState<availableActionsType[]>([]);
    const [selectedAction, setSelectedAction] = useState<selectedActionType | null>(null);
    const [selectedReaction, setSelectedReaction] = useState(null);
    const [actionServices, setActionServices] = useState([]);
    const [reactionServices, setReactionServices] = useState([]);

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
                setActionServices(propertiesServices);
            }
        } catch (err) {
            console.error("Error fetching services:", err);
        }
    };

    const getAvailableActions = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}${selectedActionService?.name}/actions`, {
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

    const getServicesWithReactions = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}reactions-services`, {
                    withCredentials: true,
            });

            if (result) {
                const propertiesServices = result.data.map((service: any) => ({
                    name: service.name,
                    color: service.color,
                    iconPath: service.logo,
                }));
                setReactionServices(propertiesServices);
            }
        } catch (err) {
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
        if (currentStep === Steps.SELECT_SERVICE_REACTION) {
            getServicesWithReactions();
        }
    }, [currentStep]);

    const handleSelectActionService = (name: string, color: string, iconPath: string) => {
        setSelectedActionService({ name, color, iconPath });
        goToStep(Steps.SELECT_ACTION);
    };

    const handleSelectAction = (name: string, id: string) => {
        setSelectedAction({name, id});
        goToStep(Steps.TRIGGER_ACTION);
    }

    const actionServiceName = selectedActionService?.name;
    const actionName = selectedAction?.name;

    const fields = actionServiceName && actionName &&
    servicesData[actionServiceName]?.actions?.[actionName]
    ? servicesData[actionServiceName].actions[actionName]
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

    const handleSelectReactionService = (name: string, color: string, iconPath: string) => {
        setSelectedReactionService({ name, color, iconPath });
        goToStep(Steps.SELECT_REACTION);
    };

    const renderStep = () => {
        switch (currentStep) {
            case Steps.HOME:
                return (
                    <HomeStep
                        selectedAction={selectedAction}
                        selectedReaction={selectedReaction}
                        selectedService={selectedActionService}
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
                        services={actionServices}
                        handleSelectService={handleSelectActionService}
                        goToStep={goToStep}
                    />
                );

            case Steps.SELECT_ACTION:
                return (
                    <SelectActionStep
                        selectedService={selectedActionService}
                        availableActions={availableActions}
                        handleSelectAction={handleSelectAction}
                        goToStep={goToStep}
                    />
                )

            case Steps.TRIGGER_ACTION:
                return (
                    <TriggerActionStep
                        selectedService={selectedActionService}
                        selectedAction={selectedAction}
                        fields={fields}
                        actionFieldValues={actionFieldValues}
                        handleFieldChange={handleFieldChange}
                        goToStep={goToStep}
                    />
                );

            case Steps.SELECT_SERVICE_REACTION:
                return (
                    <SelectServiceReactionStep
                        services={reactionServices}
                        handleSelectService={handleSelectReactionService}
                        goToStep={goToStep}
                    />
                )
        }
    }

    return (
        <>
            {renderStep()}
        </>
    );
};

export default WorkflowBuilder
