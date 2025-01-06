import { redirect, useNavigate } from "react-router-dom";
import {useEffect, useState } from "react";
import axios from "axios";
import HomeStep from "../components/workflows/HomeStep";
import SummaryStep from "../components/workflows/SummaryStep";
import SelectActionReactionStep from "../components/workflows/SelectActionReactionStep";
import ConnectServiceStep from "../components/workflows/ConnectServiceStep";
import SelectServiceStep from "../components/workflows/ServiceActionReactionStep";
import TriggerActionReactionStep from "../components/workflows/TriggerActionReactionStep";
import Steps from "../utils/CreateStepsType";
import { getFieldValues } from "../utils/getFieldValues";
import { Parameters } from "../utils/ActionReactionParameters";
import { ActionType } from "../utils/ActionType";
import { ReactionType } from "../utils/ReactionType";

type availableActionsType = {
    name: string;
    description: string;
    id: string;
    parameters: Parameters[];
};

type availableReactionsType = {
    name: string;
    description: string;
    id: string;
    parameters: Parameters[];
};

type selectedServiceType = {
    name: string;
    color: string;
    iconPath: string;
    isauthneeded: boolean;
};

const WorkflowBuilder = () => {
    const navigate = useNavigate();

    const [selectedActionService, setSelectedActionService] = useState<selectedServiceType | null>(null);
    const [selectedReactionService, setSelectedReactionService] = useState<selectedServiceType | null>(null);
    const [availableActions, setAvailableActions] = useState<availableActionsType[]>([]);
    const [availableReactions, setAvailableReactions] = useState<availableReactionsType[]>([]);
    const [selectedAction, setSelectedAction] = useState<ActionType | null>(null);
    const [selectedReaction, setSelectedReaction] = useState<ReactionType | null>(null);
    const [actionFieldValues, setActionFieldValues] = useState<Record<string, string>>({});
    const [reactionFieldValues, setReactionFieldValues] = useState<Record<string, string>>({});
    const [actionServices, setActionServices] = useState([]);
    const [reactionServices, setReactionServices] = useState([]);
    const [currentStep, setCurrentStep] = useState<Steps>(Steps.HOME);

    const goToStep = (step: Steps) => {
        setCurrentStep(step);
    };

    const handleOnEditAction = () => {
        goToStep(Steps.TRIGGER_ACTION)
    }

    const handleDeleteAction = () => {
        setSelectedActionService(null)
        setSelectedAction(null);
    }

    const handleOnEditReaction = () => {
        goToStep(Steps.TRIGGER_REACTION)
    }

    const handleDeleteReaction = () => {
        setSelectedReactionService(null)
        setSelectedReaction(null);
    }

    const goBack = () => {
        if (currentStep == Steps.HOME) {
            navigate("/")
            return;
        }
        if (currentStep == Steps.SELECT_SERVICE_ACTION) {
            goToStep(Steps.HOME)
            setSelectedActionService(null)
            setSelectedAction(null);
            return;
        }
        if (currentStep == Steps.SELECT_SERVICE_REACTION) {
            goToStep(Steps.HOME)
            setSelectedReactionService(null)
            setSelectedReaction(null);
            return;
        }
        if (currentStep == Steps.SUMMARY) {
            goToStep(Steps.HOME)
            return;
        }
        if (currentStep == Steps.CONNECT_ACTION_SERVICE) {
            goToStep(Steps.SELECT_ACTION)
            return;
        }
        if (currentStep == Steps.CONNECT_REACTION_SERVICE) {
            goToStep(Steps.SELECT_REACTION)
            return;
        }
        setCurrentStep(currentStep - 1);
    }

    const handleContinue = () => {
        goToStep(Steps.SUMMARY);
    }

    const getServicesWithActions = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}services/actions`, {
                    withCredentials: true,
            });

            if (!result) {
                console.error("Error fetching services");
                return;
            }

            const propertiesServices = result.data.map((service: any) => ({
                name: service.name,
                color: service.color,
                iconPath: service.logo,
                isauthneeded: service.isauthneeded,
            }));
            setActionServices(propertiesServices);
        } catch (err) {
            console.error("Error fetching services:", err);
        }
    };

    const getAvailableActions = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}${selectedActionService?.name}/actions`, {
                    withCredentials: true,
                });

            if (!result) {
                console.error("Unexpected API response", result);
                return;
            }

            const dataKey = "actions";

            if (Array.isArray(result.data[dataKey])) {
                const actions = result.data[dataKey].map((element: any) => ({
                    name: element.name,
                    description: element.description,
                    id: element.id,
                    parameters: element.parameters,
                }));
                setAvailableActions(actions);
            } else {
                console.error(`Error for ${dataKey}:`, result.data[dataKey]);
            }
        }
        catch (err) {
            console.error("Error fetching services:", err);
        }
    };

    const getServicesWithReactions = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}services/reactions`, {
                    withCredentials: true,
            });

            if (!result) {
                console.error("Error fetching services:");
                return;
            }

            const propertiesServices = result.data.map((service: any) => ({
                name: service.name,
                color: service.color,
                iconPath: service.logo,
                isauthneeded: service.isauthneeded,
            }));
            setReactionServices(propertiesServices);
        } catch (err) {
            console.error("Error fetching services:", err);
        }
    };

    const getAvailableReactions = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}${selectedReactionService?.name}/reactions`, {
                    withCredentials: true,
                });

            if (!result) {
                console.error("Unexpected API response");
                return;
            }

            const dataKey = "reactions";

            if (Array.isArray(result.data[dataKey])) {
                const reactions = result.data[dataKey].map((element: any) => ({
                    name: element.name,
                    description: element.description,
                    id: element.id,
                    parameters: element.parameters,
                }));
                setAvailableReactions(reactions);
            } else {
                console.error(`Error for ${dataKey}:`, result.data[dataKey]);
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
        type: "action" | "reaction",
    ): boolean => {
        if (type === "action") {
            if (selectedActionService?.isauthneeded === true) {
                return true;
            }
        }

        if (type === "reaction") {
            if (selectedReactionService?.isauthneeded === true) {
                return true;
            }
        }

        return false;
    };

    const handleSelectActionService = (name: string, color: string, iconPath: string, isauthneeded: boolean) => {
        setSelectedActionService({ name, color, iconPath, isauthneeded });
        goToStep(Steps.SELECT_ACTION);
    };

    const handleSelectAction = (name: string, id: string, parameters: Parameters[]) => {
        setSelectedAction({ name, id, parameters });

        if (selectedActionService?.name && isRequiredConnection("action")) {
            goToStep(Steps.CONNECT_ACTION_SERVICE);
        } else {
            goToStep(Steps.TRIGGER_ACTION);
        }
    };

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
    }

    const handleSelectReactionService = (name: string, color: string, iconPath: string, isauthneeded: boolean) => {
        setSelectedReactionService({ name, color, iconPath, isauthneeded });
        goToStep(Steps.SELECT_REACTION);
    };

    const handleSelectReaction = (name: string, id: string, parameters: Parameters[]) => {
        setSelectedReaction({ name, id, parameters});

        if (selectedReactionService?.name && isRequiredConnection("reaction")) {
            goToStep(Steps.CONNECT_REACTION_SERVICE);
        } else {
            goToStep(Steps.TRIGGER_REACTION);
        }
    };

    const handleConnectService = async (selectedService: any, type: "action" | "reaction") => {
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
                selectedService,
                type
            };

            sessionStorage.setItem("state-data", JSON.stringify(stateData));

            window.location.href = authUrl;
        } catch (error) {
            console.error("Error during authentication:", error);
        }
    }

    const updateStateFromSession = (state: any) => {
        if (!state) return;

        if (state.selectedActionService) {
            setSelectedActionService(state.selectedActionService);
        }
        if (state.selectedReactionService) {
            setSelectedReactionService(state.selectedReactionService);
        }
        if (state.selectedAction) {
            setSelectedAction(state.selectedAction);
        }
        if (state.selectedReaction) {
            setSelectedReaction(state.selectedReaction);
        }
        if (state.actionFieldValues) {
            setActionFieldValues(state.actionFieldValues);
        }
        if (state.reactionFieldValues) {
            setReactionFieldValues(state.reactionFieldValues);
        }
    };

    const handleServiceCallback = async (codeAuth: string, serviceName: string) => {
        if (!codeAuth || !serviceName) return;

        try {
            await axios.post(
                `${import.meta.env.VITE_API_URL}service-callback/${serviceName}/web/?code=${codeAuth}`,
                {
                    withCredentials: true
                }
            );
        } catch (error) {
            console.error("Error with service callback:", error);
        }
    };

    const clearUrlParams = () => {
        const url = new URL(window.location.href);
        url.search = "";
        window.history.replaceState({}, document.title, url.toString());
    };

    useEffect(() => {
        const getSavedData = async () => {
            const params = new URLSearchParams(window.location.search);
            const codeAuth = params.get("code");

            const state = sessionStorage.getItem("state-data");

            if (!state) {
                console.error("Undetected state into session storage");
                return;
            }

            const jsonState = JSON.parse(state);

            if (jsonState.type === "action" && jsonState.selectedAction?.parameters != null) {
                goToStep(Steps.TRIGGER_ACTION);
            }
            if (jsonState.type === "reaction" && jsonState.selectedReaction?.parameters != null) {
                goToStep(Steps.TRIGGER_REACTION);
            }

            const serviceAuth = jsonState.selectedService;

            updateStateFromSession(jsonState);

            if (!codeAuth) {
                return;
            }

            await handleServiceCallback(codeAuth, serviceAuth.name);

            sessionStorage.removeItem("state-data");

            clearUrlParams();
        };

        getSavedData();
    }, []);

    useEffect(() => {
        const updateFieldValues = (parameters: Parameters[], setFieldValues: React.Dispatch<React.SetStateAction<Record<string, string>>>) => {
            const defaultValues: Record<string, string> = {};

            parameters.forEach((param) => {
                const options = getFieldValues(param);

                if (options && options.length > 0) {
                    defaultValues[param.name] = options[0];
                }
            });

            setFieldValues((prev) => ({ ...defaultValues, ...prev }));
        };

        if (selectedAction?.parameters) {
            updateFieldValues(selectedAction.parameters, setActionFieldValues);
        }

        if (selectedReaction?.parameters) {
            updateFieldValues(selectedReaction.parameters, setReactionFieldValues);
        }

    }, [selectedAction, selectedReaction]);

    const getIndexFromValues = (value: string, valuesArray: string[]): number | null => {
        const index = valuesArray.indexOf(value);

        if (index === -1) {
            return null;
        }

        if (valuesArray[0] === "0") {
            return index;
        }

        return index + 1;
    };

    const getActionReactionParams = (
        valuesSelectedByUser: Record<string, string>,
        parameters: Parameters[]
    ): Record<string, any> => {
        return Object.keys(valuesSelectedByUser).reduce((params, key) => {
            const param = parameters.find((param) => param.name === key);

            if (param) {
                if (param.type === "int") {
                    const valuesArray = getFieldValues(param);

                    if (valuesArray) {
                        const index = getIndexFromValues(valuesSelectedByUser[key], valuesArray);

                        if (index !== null) {
                            params[key] = index;
                        }
                    }

                } else {
                    params[key] = valuesSelectedByUser[key];
                }
            }

            return params;
        }, {} as Record<string, any>);
    };

    const handleCreateWorkflow = async (workflowTitle: string) => {
        const actionParams = selectedAction
            ? getActionReactionParams(actionFieldValues, selectedAction.parameters)
            : null;

        const reactionParams = selectedReaction
            ? getActionReactionParams(reactionFieldValues, selectedReaction.parameters)
            : null;

        await axios.post(
            `${import.meta.env.VITE_API_URL}workflows`,
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

        setSelectedActionService(null);
        setSelectedAction(null);
        setSelectedReaction(null);
        navigate("/");
    };

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
                        handleOnEditReaction={handleOnEditReaction}
                        handleOnDeleteReaction={handleDeleteReaction}
                        handleContinue={handleContinue}
                    />
                )

            case Steps.SELECT_SERVICE_ACTION:
                return (
                    <SelectServiceStep
                        services={actionServices}
                        handleSelectService={handleSelectActionService}
                        goBack={() => goBack()}
                    />
                );

            case Steps.SELECT_ACTION:
                return (
                    <SelectActionReactionStep
                        selectedService={selectedActionService}
                        availableItems={availableActions}
                        handleSelect={handleSelectAction}
                        goBack={() => goToStep(Steps.SELECT_SERVICE_ACTION)}
                    />
                )

            case Steps.TRIGGER_ACTION:
                return (
                    <TriggerActionReactionStep
                        selectedService={selectedActionService}
                        selectedItem={selectedAction}
                        itemFieldValues={actionFieldValues}
                        handleFieldChange={handleFieldChange}
                        goToStep={() => goToStep(Steps.HOME)}
                        goBack={() => goBack()}
                        type="action"
                    />
                );

            case Steps.SELECT_SERVICE_REACTION:
                return (
                    <SelectServiceStep
                        services={reactionServices}
                        handleSelectService={handleSelectReactionService}
                        goBack={() => goBack()}
                    />
                )

            case Steps.SELECT_REACTION:
                return (
                    <SelectActionReactionStep
                        selectedService={selectedReactionService}
                        availableItems={availableReactions}
                        handleSelect={handleSelectReaction}
                        goBack={() => goToStep(Steps.SELECT_SERVICE_REACTION)}
                    />
                )

            case Steps.TRIGGER_REACTION:
                return (
                    <TriggerActionReactionStep
                        selectedService={selectedReactionService}
                        selectedItem={selectedReaction}
                        itemFieldValues={reactionFieldValues}
                        handleFieldChange={handleFieldChange}
                        goToStep={() => goToStep(Steps.HOME)}
                        goBack={() => goBack()}
                        type="reaction"
                    />
                )

            case Steps.CONNECT_ACTION_SERVICE:
                    return (
                        <ConnectServiceStep
                            selectedService={selectedActionService}
                            goBack={goBack}
                            handleConnect={() => handleConnectService(selectedActionService, "action")}
                        />
                    )

            case Steps.CONNECT_REACTION_SERVICE:
                return (
                    <ConnectServiceStep
                        selectedService={selectedReactionService}
                        goBack={() => goBack()}
                        handleConnect={() => handleConnectService(selectedReactionService, "reaction")}
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
