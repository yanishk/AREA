import { useNavigate, useSearchParams } from "react-router-dom";
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
import ConnectServiceStep from "../components/workflows/ConnectActionServiceStep";
import SummaryStep from "../components/workflows/SummaryStep";

enum Steps {
    HOME,
    SELECT_SERVICE_ACTION,
    SELECT_ACTION,
    CONNECT_ACTION_SERVICE,
    TRIGGER_ACTION,
    SELECT_SERVICE_REACTION,
    SELECT_REACTION,
    CONNECT_REACTION_SERVICE,
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
    const [searchParams] = useSearchParams();

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

    const goBack = () => {
        if (currentStep == Steps.HOME) {
            navigate(-1)
        }
        if (currentStep == Steps.SELECT_SERVICE_ACTION) {
            goToStep(Steps.HOME)
            setSelectedActionService(null)
            setSelectedAction(null);
        }
        if (currentStep == Steps.SELECT_ACTION) {
            goToStep(Steps.SELECT_SERVICE_ACTION)
        }
        if (currentStep == Steps.CONNECT_ACTION_SERVICE) {
            goToStep(Steps.SELECT_SERVICE_ACTION)
        }
        if (currentStep == Steps.TRIGGER_ACTION) {
            goToStep(Steps.SELECT_ACTION)
        }
        if (currentStep == Steps.SELECT_SERVICE_REACTION) {
            goToStep(Steps.HOME)
            setSelectedReaction(null)
            setSelectedReaction(null);
        }
        if (currentStep == Steps.SELECT_REACTION) {
            goToStep(Steps.SELECT_SERVICE_REACTION)
        }
        if (currentStep == Steps.CONNECT_REACTION_SERVICE) {
            goToStep(Steps.SELECT_REACTION)
        }
        if (currentStep == Steps.TRIGGER_REACTION) {
            goToStep(Steps.CONNECT_REACTION_SERVICE)
        }
        if (currentStep == Steps.SUMMARY) {
            goToStep(Steps.HOME)
        }
        return;
    }

    const handleContinue = () => {
        goToStep(Steps.SUMMARY);
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

    const isRequiredConnection = (
        serviceName: string,
        type: "action" | "reaction",
        actionOrReactionName?: string
    ): boolean => {
        const service = servicesData[serviceName];
        if (!service || !actionOrReactionName) {
            return false;
        }

        if (type === "action") {
            const action = service.actions?.[actionOrReactionName];

            if (!action) {
                return false;
            }

            return action.some((actionData) => actionData.auth === "yes");
        }

        if (type === "reaction") {
            const reaction = service.reactions?.[actionOrReactionName];

            if (!reaction) {
                return false;
            }

            return reaction.some((reactionData) => reactionData.auth === "yes");
        }

        return false;
    };

    const handleSelectActionService = (name: string, color: string, iconPath: string) => {
        setSelectedActionService({ name, color, iconPath });
        goToStep(Steps.SELECT_ACTION);
    };

    const handleSelectAction = (name: string, id: string) => {
        setSelectedAction({ name, id });

        if (selectedActionService?.name && isRequiredConnection(selectedActionService.name, "action", name)) {
            goToStep(Steps.CONNECT_ACTION_SERVICE);
        } else {
            goToStep(Steps.TRIGGER_ACTION);
        }
    };

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

    const handleSelectReaction = (name: string, id: string) => {
        setSelectedReaction({ name, id });

        if (selectedReactionService?.name && isRequiredConnection(selectedReactionService.name, "reaction", name)) {
            goToStep(Steps.CONNECT_REACTION_SERVICE);
        } else {
            goToStep(Steps.TRIGGER_REACTION);
        }
    };

    const handleConnectService = async (selectedService: any) => {
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
    }

    useEffect(() => {
        const getDataFromUrl = async () => {
            const params = new URLSearchParams(window.location.search);

            const codeAuth = searchParams.get("code");
            const state = params.get("state");
            let decodedState;

            if (state) {
                try {
                    decodedState = JSON.parse(decodeURIComponent(state));

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

            if (codeAuth) {
                try {
                    await axios.post(
                        `${import.meta.env.VITE_API_URL}service-callback/${decodedState.serviceType === "action" ? selectedActionService?.name : selectedReactionService?.name}/web/?code=${codeAuth}`,{},
                        {
                            withCredentials: true,
                        }
                    );
                } catch (error) {
                    console.error("Error during service callback:", error);
                }
            }

            const url = new URL(window.location.href);
            url.search = "";
            window.history.replaceState({}, document.title, url.toString());
        }

        getDataFromUrl();
    }, []);

    const handleCreateWorkflow = async (workflowTitle: string) =>{
            const actionParams = Object.keys(actionFieldValues).length > 0
            ? Object.keys(actionFieldValues).reduce((params, key) => {
                params[key] = actionFieldValues[key];
                return params;
            }, {} as Record<string, any>)
            : null;

            const reactionParams = Object.keys(reactionFieldValues).length > 0
                ? Object.keys(reactionFieldValues).reduce((params, key) => {
                    params[key] = reactionFieldValues[key];
                    return params;
            }, {} as Record<string, any>)
            : null;

            await axios.post(
                `${import.meta.env.VITE_API_URL}workflow/create`,
                {
                    actionid: selectedAction?.id,
                    actionparam: actionParams,
                    name: workflowTitle || "My Workflow",
                    reactionid: selectedReaction?.id,
                    reactionparam: reactionParams
                },
                {
                    withCredentials: true,
                }
            );

            setSelectedActionService(null)
            setSelectedAction(null);
            setSelectedReaction(null)
            setSelectedReaction(null);
            navigate("/");
    }

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
                        goHome={() => goBack()}
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
                        goBack={() => goBack()}
                    />
                );

            case Steps.SELECT_ACTION:
                return (
                    <SelectActionStep
                        selectedService={selectedActionService}
                        availableActions={availableActions}
                        handleSelectAction={handleSelectAction}
                        goBack={() => goBack()}
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
                        goToStep={() => goToStep(Steps.HOME)}
                        goBack={() => goBack()}
                    />
                );

            case Steps.SELECT_SERVICE_REACTION:
                return (
                    <SelectServiceReactionStep
                        services={reactionServices}
                        handleSelectService={handleSelectReactionService}
                        goBack={() => goBack()}
                    />
                )

            case Steps.SELECT_REACTION:
                return (
                    <SelectReactionStep
                        selectedService={selectedReactionService}
                        availableReactions={availableReactions}
                        handleSelectReaction={handleSelectReaction}
                        goBack={() => goBack()}
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
                        goBack={() => goBack()}
                    />
                )

            case Steps.CONNECT_ACTION_SERVICE:
                    return (
                        <ConnectServiceStep
                            selectedService={selectedActionService}
                            goBack={goBack}
                            handleConnect={() => handleConnectService(selectedActionService)}
                        />
                    )

            case Steps.CONNECT_REACTION_SERVICE:
                return (
                    <ConnectServiceStep
                        selectedService={selectedReactionService}
                        goBack={() => goBack()}
                        handleConnect={() => handleConnectService(selectedReactionService)}
                    />
                )

            case Steps.SUMMARY:
                return (
                    <SummaryStep
                        selectedActionService={selectedActionService}
                        selectedReactionService={selectedReactionService}
                        selectedAction={selectedAction}
                        selectedReaction={selectedReaction}
                        actionFields={actionFieldValues}
                        reactionFields={reactionFieldValues}
                        goBack={() => goBack()}
                        handleCreateWorkflow={handleCreateWorkflow}
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
