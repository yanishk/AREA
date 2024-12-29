import { useNavigate } from "react-router-dom";
import {useEffect, useMemo, useState } from "react";
import axios from "axios";
import {servicesData } from "../ServicesData";
import HomeStep from "../components/workflows/HomeStep";
import SelectServiceActionStep from "../components/workflows/ServiceActionStep";
import SelectActionStep from "../components/workflows/SelectActionStep";
import TriggerActionStep from "../components/workflows/TriggerActionStep";
import SelectServiceReactionStep from "../components/workflows/ServiceReactionStep";
import SelectReactionStep from "../components/workflows/SelectReactionStep";
import TriggerReactionStep from "../components/workflows/TriggerReactionStep";
import ConnectServiceStep from "../components/workflows/ConnectServiceStep";

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

type availableReactionsType = {
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

type selectedReactionType = {
    name: string;
    id: string;
};

const WorkflowBuilder = () => {
    const navigate = useNavigate();

    const [selectedActionService, setSelectedActionService] = useState<selectedServiceType | null>(null);
    const [selectedReactionService, setSelectedReactionService] = useState<selectedServiceType | null>(null);
    const [availableActions, setAvailableActions] = useState<availableActionsType[]>([]);
    const [availableReactions, setAvailableReactions] = useState<availableReactionsType[]>([]);
    const [selectedAction, setSelectedAction] = useState<selectedActionType | null>(null);
    const [selectedReaction, setSelectedReaction] = useState<selectedReactionType | null>(null);
    const [actionFieldValues, setActionFieldValues] = useState<Record<string, string>>({});
    const [reactionFieldValues, setReactionFieldValues] = useState<Record<string, string>>({});
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

    const getAvailableReactions = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}${selectedReactionService?.name}/reactions`, {
                    withCredentials: true,
                });

            if (result) {
                const dataKey = "reactions";

                if (Array.isArray(result.data[dataKey])) {
                    const reactions = result.data[dataKey].map((element: any) => ({
                        name: element.name,
                        description: element.description,
                        id: element.id,
                    }));
                    setAvailableReactions(reactions);
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
        if (currentStep === Steps.SELECT_SERVICE_REACTION) {
            getServicesWithReactions();
        }
        if (currentStep === Steps.SELECT_REACTION) {
            getAvailableReactions();
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

    const reactionServiceName = selectedReactionService?.name;
    const reactionName = selectedReaction?.name;

    const actionFields = useMemo(() => {
        return actionServiceName && actionName &&
            servicesData[actionServiceName]?.actions?.[actionName]
            ? servicesData[actionServiceName].actions[actionName]
            : [];
    }, [actionServiceName, actionName]);

    const reactionFields = useMemo(() => {
        return reactionServiceName && reactionName &&
            servicesData[reactionServiceName]?.actions?.[reactionName]
            ? servicesData[reactionServiceName].actions[reactionName]
            : [];
    }, [reactionServiceName, reactionName]);

    const handleFieldChange = (id: string, value: string, type: "action" | "reaction") => {
        if (type === "action") {
            setActionFieldValues((prev) => ({
                ...prev,
                [id]: value,
            }));
        }
        if (type === "reaction") {
            setReactionFieldValues((prev) => ({
                ...prev,
                [id]: value,
            }));
        }
        else {
            console.error(`Invalid type parameter: ${type}`);
        }
    };

    useEffect(() => {
        if (currentStep === Steps.TRIGGER_ACTION) {
            const defaultValues = actionFields.reduce((acc, field) => {
                acc[field.id] = field.options?.[0] ?? "";
                return acc;
            }, {} as Record<string, string>);
            setActionFieldValues(defaultValues);
        } else if (currentStep === Steps.TRIGGER_REACTION) {
            const defaultValues = reactionFields.reduce((acc, field) => {
                acc[field.id] = field.options?.[0] ?? "";
                return acc;
            }, {} as Record<string, string>);
            setReactionFieldValues(defaultValues);
        }
    }, [actionFields, reactionFields, currentStep]);

    const handleSelectReactionService = (name: string, color: string, iconPath: string) => {
        setSelectedReactionService({ name, color, iconPath });
        goToStep(Steps.SELECT_REACTION);
    };

    const getConnectReactionStatus = () =>
        selectedReactionService?.name && selectedReaction?.name
            ? servicesData[selectedReactionService?.name]?.reactions?.[selectedReaction?.name]?.[0]?.serviceConnection
            : undefined;

    const getTiggerReactionStatus = () =>
        selectedReactionService?.name && selectedReaction?.name
            ? servicesData[selectedReactionService?.name]?.reactions?.[selectedReaction?.name]?.[0]?.triggerReaction
            : undefined;

    const handleSelectReaction = (name: string, id: string) => {
        setSelectedReaction({name, id});

        if (getConnectReactionStatus() === "yes") {
            goToStep(Steps.CONNECT_SERVICE);
        }
    }

    const handleConnect = async (selectedService: any) => {
        try {
            const result = await axios.get(
                `${import.meta.env.VITE_API_URL}authentication/${selectedService.name}/service/web`,
                {
                    withCredentials: true,
                }
            );

            const authUrl = result.data["auth-url"];
            if (!authUrl) {
                console.error("Error: Authentication URL not received");
                return;
            }

            const stateData = {
                selectedActionService,
                selectedReactionService,
                selectedAction,
                selectedReaction,
                actionFieldValues,
                reactionFieldValues,
            };
            const state = encodeURIComponent(JSON.stringify(stateData));

            const redirectUrl = `${authUrl}&state=${state}`;

            window.location.href = redirectUrl;
        } catch (error) {
            console.error("Error during authentication:", error);
        }
        // if (getTiggerReactionStatus() === "yes") {
        //     goToStep(Steps.TRIGGER_REACTION);
        // } else {
        //     goToStep(Steps.HOME);
        // }
    }

    useEffect(() => {
        const params = new URLSearchParams(window.location.search);

        const state = params.get("state");
        if (state) {
            try {
                const decodedState = JSON.parse(decodeURIComponent(state));

                if (decodedState.selectedActionService) {
                    setSelectedActionService(decodedState.selectedActionService);
                }
                if (decodedState.selectedReactionService) {
                    setSelectedReactionService(decodedState.selectedReactionService);
                }
                if (decodedState.selectedAction) {
                    setSelectedAction(decodedState.selectedAction);
                }
                if (decodedState.selectedReaction) {
                    setSelectedReaction(decodedState.selectedReaction);
                }
                if (decodedState.actionFieldValues) {
                    setActionFieldValues(decodedState.actionFieldValues);
                }
                if (decodedState.reactionFieldValues) {
                    setReactionFieldValues(decodedState.reactionFieldValues);
                }
            } catch (error) {
                console.error("Error parsing state:", error);
            }
        }

        // const code = params.get("code");

        const url = new URL(window.location.href);
        url.search = "";
        window.history.replaceState({}, document.title, url.toString());
    }, []);

    const renderStep = () => {
        switch (currentStep) {
            case Steps.HOME:
                return (
                    <HomeStep
                        selectedAction={selectedAction}
                        selectedReaction={selectedReaction}
                        selectedServiceAction={selectedActionService}
                        selectedServiceReaction={selectedReactionService}
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
                        fields={actionFields}
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

            case Steps.SELECT_REACTION:
                return (
                    <SelectReactionStep
                        selectedService={selectedReactionService}
                        availableReactions={availableReactions}
                        handleSelectReaction={handleSelectReaction}
                        goToStep={goToStep}
                    />
                )

            case Steps.TRIGGER_REACTION:
                return (
                    <TriggerReactionStep
                        selectedService={selectedReactionService}
                        selectedReaction={selectedReaction}
                        fields={reactionFields}
                        reactionFieldValues={reactionFieldValues}
                        handleFieldChange={handleFieldChange}
                        goToStep={goToStep}
                    />
                )

            case Steps.CONNECT_SERVICE:
                return (
                    <ConnectServiceStep
                        selectedService={selectedReactionService}
                        goBack={goHome}
                        handleConnect={() => handleConnect(selectedReactionService)}
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
