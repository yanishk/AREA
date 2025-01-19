import { useNavigate } from "react-router-dom";
import {useEffect, useState } from "react";
import axios from "axios";
import HomeStep from "../components/workflows/HomeStep";
import SummaryStep from "../components/workflows/SummaryStep";
import SelectActionReactionStep from "../components/workflows/SelectActionReactionStep";
import ConnectServiceStep from "../components/workflows/ConnectServiceStep";
import SelectServiceStep from "../components/workflows/ServiceActionReactionStep";
import TriggerActionReactionStep from "../components/workflows/TriggerActionReactionStep";
import Steps from "../utils/CreateStepsType";
import { Parameters } from "../utils/ActionReactionParameters";
import { ActionType } from "../utils/ActionType";
import { ReactionType } from "../utils/ReactionType";
import getFieldValues from "../utils/GetFieldValues";
import { setSelectedWorkspace } from "../utils/AsanaUtils";
import { getDiscordChannels, getServerIdByName, setSelectedServer } from "../utils/DiscordUtils";

type availableActionsType = {
    name: string;
    description: string;
    id: string;
    nbparam: number;
    parameters: Parameters[];
};

type availableReactionsType = {
    name: string;
    description: string;
    id: string;
    nbparam: number;
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
    const [isServiceConnected, setIsServiceConnected] = useState(false);

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
            handleDeleteAction();
            handleDeleteReaction();
            sessionStorage.removeItem("state-data");
            navigate(-1);
            return;
        }
        if (currentStep == Steps.SELECT_SERVICE_ACTION) {
            goToStep(Steps.HOME)
            handleDeleteAction();
            return;
        }
        if (currentStep == Steps.SELECT_SERVICE_REACTION) {
            goToStep(Steps.HOME)
            handleDeleteReaction();
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
                    nbparam: element.nbparam,
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
                    nbparam: element.nbparam,
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

    const isServiceAlreadyConnected = async (selectedService: any): Promise<boolean> => {
        try {
            const result = await axios.get(
                `${import.meta.env.VITE_API_URL}service-authentication-status?service=${selectedService.name}`,
                {
                    withCredentials: true,
                }
            );

            return result.data.authenticated;
        } catch (error) {
            console.error("Error with service authentication status:", error);
            return false;
        }
    };

    const isRequiredConnection = async (
        type: "action" | "reaction",
    ): Promise<boolean> => {
        if (type === "action") {
            if (selectedActionService?.isauthneeded === true && !(await isServiceAlreadyConnected(selectedActionService))) {
                return true;
            }
        }

        if (type === "reaction") {
            if (selectedReactionService?.isauthneeded === true && !(await isServiceAlreadyConnected(selectedReactionService))) {
                return true;
            }
        }

        return false;
    };

    const handleSelectActionService = (name: string, color: string, iconPath: string, isauthneeded: boolean) => {
        setSelectedActionService({ name, color, iconPath, isauthneeded });
        goToStep(Steps.SELECT_ACTION);
    };

    const handleSelectAction = async (name: string, id: string, nbparam: number, parameters: Parameters[]) => {
        setSelectedAction({ name, id, nbparam, parameters });

        if (selectedActionService?.name && (await isRequiredConnection("action"))) {
            goToStep(Steps.CONNECT_ACTION_SERVICE);
            return;
        }

        if (nbparam === 0) {
            goToStep(Steps.HOME);
            return;
        }

        goToStep(Steps.TRIGGER_ACTION);
    };

    const handleFieldChange = async (id: string, value: string, type: "action" | "reaction") => {
        if (type === "action") {
            setActionFieldValues((prev) => ({
                ...prev,
                [id]: value,
            }));
        }
        if (type === "reaction") {
            if (id === "server") {
                setSelectedServer(value, id);

                const serverId = await getServerIdByName(value);
                if (serverId) {
                    const channels = await getDiscordChannels(serverId);

                    if (channels.length > 0) {
                        setReactionFieldValues((prev) => ({
                            ...prev,
                            channel: channels[0].name,
                        }));
                    }
                }
            }

            if (id === "workflow") {
                setSelectedWorkspace(value, id);
            }

            setReactionFieldValues((prev) => ({
                ...prev,
                [id]: value,
            }));
        }
    };

    const handleSelectReactionService = (name: string, color: string, iconPath: string, isauthneeded: boolean) => {
        setSelectedReactionService({ name, color, iconPath, isauthneeded });
        goToStep(Steps.SELECT_REACTION);
    };

    const handleSelectReaction = async (name: string, id: string, nbparam: number, parameters: Parameters[]) => {
        setSelectedReaction({ name, id, nbparam, parameters});

        if (selectedReactionService?.name && await isRequiredConnection("reaction")) {
            goToStep(Steps.CONNECT_REACTION_SERVICE);
            return;
        }

        if (nbparam === 0) {
            goToStep(Steps.HOME);
            return;
        }

        goToStep(Steps.TRIGGER_REACTION);
    };

    const handleConnectService = async (selectedService: any, type: "action" | "reaction") => {
        try {
            const result = await axios.get(
                `${import.meta.env.VITE_API_URL}authentication?service=${selectedService.name}&callbacktype=service&apptype=web`,
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

    let isCallbackSent = false;

    const handleServiceCallback = async (codeAuth: string, serviceName: string) => {
        if (!codeAuth || !serviceName) return;

        if (isCallbackSent) {
            console.warn("Callback already sent");
            return;
        }

        isCallbackSent = true;

        try {
            const result = await axios.post(
                `${import.meta.env.VITE_API_URL}service-callback/?code=${codeAuth}`,
                {
                    apptype: "web",
                    service: serviceName
                },
                {
                    withCredentials: true
                }
            );

            if (result.data?.success) {
                setIsServiceConnected(true);
            }

        } catch (error) {
            console.error("Error with service callback:", error);
        }
    };

    const clearUrlParams = () => {
        const url = new URL(window.location.href);
        url.search = "";
        window.history.replaceState({}, document.title, url.toString());
    };

    const isServiceConnectedStatus = async (Servicename: string): Promise<boolean> => {
        try {
            const result = await axios.get(
                `${import.meta.env.VITE_API_URL}service-authentication-status?service=${Servicename}`,
                {
                    withCredentials: true,
                }
            );

            return result.data.authenticated;
        } catch (error) {
            console.error("Error with service authentication status:", error);
            return false;
        }
    };


    useEffect(() => {
        const getSavedData = async () => {
            const params = new URLSearchParams(window.location.search);
            const codeAuth = params.get("code");

            const state = sessionStorage.getItem("state-data");

            if (!state) {
                console.warn("Undetected state into session storage");
                return;
            }

            const jsonState = JSON.parse(state);

            const serviceAuth = jsonState.selectedService;

            updateStateFromSession(jsonState);

            if (!codeAuth) {
                return;
            }

            await handleServiceCallback(codeAuth, serviceAuth.name);

            if (jsonState.type === "action" && jsonState.selectedAction?.parameters != null) {
                const isconnected = await isServiceConnectedStatus(jsonState.selectedActionService?.name)

                if (isconnected) {
                    goToStep(Steps.TRIGGER_ACTION);
                }
                return;
            }
            if (jsonState.type === "reaction" && jsonState.selectedReaction?.parameters != null) {
                const isconnected = await isServiceConnectedStatus(jsonState.selectedReactionService?.name)

                if (isconnected) {
                    goToStep(Steps.TRIGGER_REACTION);
                }
                return;
            }

            sessionStorage.removeItem("state-data");

            clearUrlParams();
        };

        getSavedData();
    }, []);

    const updateFieldValues = async (parameters: Parameters[], setFieldValues: React.Dispatch<React.SetStateAction<Record<string, string>>>) => {
        const defaultValues: Record<string, string> = {};

        for (const param of parameters) {
            const options = await getFieldValues(param);

            if (options && options.length > 0) {
                defaultValues[param.name] = options[0];
            }
        }

        setFieldValues((prev) => ({ ...defaultValues, ...prev }));
    };

    useEffect(() => {
        const fetchFieldValues = async () => {
            if (selectedAction?.parameters) {
                if (await isRequiredConnection("action") && !isServiceConnected) {
                    return;
                }
                await updateFieldValues(selectedAction.parameters, setActionFieldValues);
                setIsServiceConnected(false);
            }
            if (selectedReaction?.parameters) {
                if (await isRequiredConnection("reaction") && !isServiceConnected) {
                    return;
                }
                await updateFieldValues(selectedReaction.parameters, setReactionFieldValues);
                setIsServiceConnected(false);
            }
        };

        const test = async () => {
            if (isServiceConnected) {
                fetchFieldValues();
            }
        }

        test();

    }, [selectedAction, selectedReaction, isServiceConnected]);

    const getActionReactionParams = async (
        valuesSelectedByUser: Record<string, string>,
        parameters: Parameters[]
    ): Promise<Record<string, any>> => {
        const params: Record<string, any> = {};

        if (!parameters) {
            return params;
        }

        for (const param of parameters) {
            params[param.name] = valuesSelectedByUser[param.name];
        }

        return params;
    };

    const getGitlabProjectId = async (actionParams: any) => {
        try {
            const result = await axios.get(
                `${import.meta.env.VITE_API_URL}gitlab/user/projects`,
                    {
                        withCredentials: true,
                    }
                );

                const gitlabProjects = result.data.projects;

                Object.keys(actionParams).forEach((key) => {
                    const projectName = actionParams[key];

                    const project = gitlabProjects.find((item: { name: string }) => item.name === projectName);

                    if (!project) {
                        console.error("No gitlab projects found");
                        return;
                    }
                    actionParams[key] = project.id.toString();
                }
            );
        } catch (error) {
            console.error("Error with fetch gitlab projects", error);
        }
    }

    const handleDiscordParameters = async (reactionParams: any) => {
        try {
            const result = await axios.get(
                `${import.meta.env.VITE_API_URL}discord/user/servers`,
                    {
                        withCredentials: true,
                    }
                );

                const discordServers = result.data.servers;

                Object.keys(reactionParams).forEach((key) => {
                    const projectName = reactionParams[key];

                    const server = discordServers.find((item: { name: string }) => item.name === projectName);

                    if (!server) {
                        console.error("No discord servers found");
                        return;
                    }

                    reactionParams[key] = server.id;
                }
            );
        } catch (error) {
            console.error("Error with fetch discord servers", error);
        }
    }

    const handleCreateWorkflow = async (workflowTitle: string) => {
        const actionParams = selectedAction
            ? await getActionReactionParams(actionFieldValues, selectedAction.parameters)
            : null;

        const reactionParams = selectedReaction
            ? await getActionReactionParams(reactionFieldValues, selectedReaction.parameters)
            : null;

        if (selectedActionService?.name === "Gitlab" && actionParams) {
            await getGitlabProjectId(actionParams);
        }

        if (selectedReactionService?.name === "Discord" && reactionParams) {
            await handleDiscordParameters(reactionParams);
        }

        if (selectedAction?.parameters) {
            selectedAction.parameters.forEach((param) => {
                if (param.type != "int" || actionParams === null) {
                    return null;
                }

                const selectedParam = actionParams?.[param.name];
                let index = param.values.indexOf(selectedParam);

                if (param.values[0] === "0") {
                    actionParams[param.name] = index;
                }
                if (index != -1) {
                    actionParams[param.name] = index + 1;
                } else {
                    actionParams[param.name] = parseInt(actionParams[param.name], 10);
                }

                return actionParams[param.name]
            });
        }

        if (selectedReaction?.parameters) {
            selectedReaction.parameters.forEach((param) => {
                if (param.type != "int" || reactionParams === null)
                    return null;
                if (param?.name != "server" && param?.name != "channel") {
                    reactionParams[param.name] = parseInt(reactionParams[param.name], 10);
                }
            });
        }

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

        navigate("/my_workflows");

        sessionStorage.removeItem("state-data");
        setSelectedActionService(null);
        setSelectedAction(null);
        setSelectedReaction(null);
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
                        type="action"
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
                        type="reaction"
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
