import { useEffect, useState } from "react";
import { Parameters } from "../../utils/ActionReactionParameters";
import Steps from "../../utils/CreateStepsType";
import { fieldType } from "../../utils/FieldType";
import TriggerButton from "../buttons/TriggerButton";
import DynamicFields from "../inputs/DynamicFields";
import ActionButton from "../buttons/ActionButton";
import { getAsanaAssignees, getAsanaProjects, getAsanaTags, getAsanaWorkspaces, getWorkspaceIdByName, setSelectedWorkspace } from "../../utils/AsanaUtils";
import SelectField from "../inputs/SelectField";
import SelectFieldId from "../inputs/SelectFieldId";
import getFieldValues from "../../utils/GetFieldValues";
import { getDiscordChannels, getDiscordServers, getServerIdByName, setSelectedServer } from "../../utils/DiscordUtils";

interface TriggerActionReactionProps {
    selectedService: { name: string; color: string; iconPath: string } | null;
    selectedItem: {name: string, id: string, parameters: Parameters[]} | null;
    itemFieldValues: Record<string, string>;
    handleFieldChange: (id: string, value: string, type: "action" | "reaction") => void;
    goToStep: (step: number) => void;
    goBack: () => void;
    type: "action" | "reaction";
}

const TriggerActionReactionStep: React.FC<TriggerActionReactionProps> = ({
    selectedService,
    selectedItem,
    itemFieldValues,
    handleFieldChange,
    goToStep,
    goBack,
    type,
}) => {

    const [isLoading, setIsLoading] = useState(true);
    const [isWorkspaceLoaded, setIsWorkspaceLoaded] = useState(true);
    const [isServerLoaded, setIsServerLoaded] = useState(true);
    const [isFieldLoading, setIsFieldLoading] = useState(false);
    const [asanaField, setAsanaField] = useState(false);
    const [discordField, setDiscordField] = useState(false);
    const [servers, setServers] = useState<string[]>([]);
    const [channels, setChannels] = useState<string[]>([]);
    const [channelsId, setChannelsId] = useState<string[]>([]);
    const [workspaces, setWorkspaces] = useState<string[]>([]);
    const [projects, setProjects] = useState<string[]>([]);
    const [projectsId, setProjectsId] = useState<string[]>([]);
    const [assignees, setAssignees] = useState<string[]>([]);
    const [assigneesId, setAssigneesId] = useState<string[]>([]);
    const [tags, setTags] = useState<string[]>([]);
    const [tagsId, setTagsId] = useState<string[]>([]);
    const [selectedWorkspaceName, setSelectedWorkspaceName] = useState<string | null>(null);
    const [selectedServerName, setSelectedServerName] = useState<string | null>(null);

    const getFieldType = (param: Parameters) => {
        if (param.isexhaustive || (!param.isexhaustive && param.route) || param.route || param.values[0]?.length >= 1) {
            return fieldType.SELECT;
        }

        return fieldType.INPUT;
    };

    const [fields, setFields] = useState<{
        id: string;
        label: string;
        type: fieldType;
        options?: string[]
        }[]
    > ([]);

    useEffect(() => {
        if (selectedService?.name === "Asana") {
            const fetchWorkspaces = async () => {
                const workspaces = await getAsanaWorkspaces();

                setWorkspaces(workspaces.map((ws) => ws.name));
                if (workspaces.length > 0) {
                    setSelectedWorkspace(workspaces[0].gid, workspaces[0].name);
                    await handleWorkspaceChange(workspaces[0].name);
                }
            };

            fetchWorkspaces();
        }

        if (selectedService?.name === "Discord") {
            const fetchServers = async () => {
                const servers = await getDiscordServers();

                setServers(servers.map((server) => server.name));
                if (servers.length > 0) {
                    setSelectedServer(servers[0].id, servers[0].name);
                    await handleServerChange(servers[0].name);
                }
            };

            fetchServers();
        }

    }, []);

    useEffect(() => {
        const fetchFields = async () => {
            if (!selectedItem) return;

            setIsLoading(true);

            const fetchedFields = await Promise.all(
                selectedItem.parameters.map(async (param) => {
                    const options =
                        getFieldType(param) === fieldType.SELECT
                            ? await getFieldValues(param)
                            : undefined;

                            if (!itemFieldValues[param.name] && options && options.length > 0) {
                                handleFieldChange(param.name, options[0], type);
                            }

                    return {
                        id: param.name,
                        label: param.name,
                        type: getFieldType(param),
                        options: options || [],
                    };
                })
            );

            setFields(fetchedFields);
            setIsLoading(false);
        };

        fetchFields();
    }, [selectedItem]);

    useEffect(() => {
        if (selectedItem?.parameters.some((param) => param.route === "/asana/workspace/projects")) {
            setAsanaField(true);
        }
        if (selectedItem?.parameters.some((param) => param.route === "/discord/user/servers")) {
            setDiscordField(true);
        }
    });

    const handleWorkspaceChange = async (workspaceName: string) => {
        setIsWorkspaceLoaded(false);
        setIsFieldLoading(true);
        setSelectedWorkspaceName(workspaceName);
        const workspaceId = await getWorkspaceIdByName(workspaceName);

        if (workspaceId) {
            setSelectedWorkspace(workspaceId, workspaceName);

            const [projects, assignees, tags] = await Promise.all([
                getAsanaProjects(workspaceId),
                getAsanaAssignees(workspaceId),
                getAsanaTags(workspaceId),
            ]);

            const [projectsId, assigneesId, tagsId] = await Promise.all([
                getAsanaProjects(workspaceId),
                getAsanaAssignees(workspaceId),
                getAsanaTags(workspaceId),
            ]);

            setIsWorkspaceLoaded(true);
            setIsFieldLoading(false);

            if (projects.length === 0) {
                setProjects(["No projects"]);
                setProjectsId(["No projects"]);
            } else {
                setProjects(projects.map((project) => project.name));
                setProjectsId(projectsId.map((project) => project.gid));
            }

            if (assignees.length === 0) {
                setAssignees(["No assignees"]);
                setAssigneesId(["No assignees"]);
            } else {
                setAssignees(assignees.map((assignee) => assignee.name));
                setAssigneesId(assigneesId.map((assignee) => assignee.gid));
            }

            if (tags.length === 0) {
                setTags(["No tags"]);
                setTagsId(["No tags"]);
            } else {
                setTags(tags.map((tag) => tag.name));
                setTagsId(tagsId.map((tag) => tag.gid));
            }
            handleFieldChange("project", projectsId[0].gid, type);
            handleFieldChange("assignee", assigneesId[0].gid, type);
            handleFieldChange("tag", tagsId[0].gid, type);
            handleFieldChange("workspace", workspaceId, type);
        }
    };

    const handleServerChange = async (serverName: string) => {
        setIsServerLoaded(false);
        setIsFieldLoading(true);
        setSelectedServerName(serverName);
        const serverId = await getServerIdByName(serverName);

        if (serverId) {
            setSelectedServer(serverId, serverName);

            const channels = await getDiscordChannels(serverId);
            const channelsId = await getDiscordChannels(serverId);

            setIsServerLoaded(true);
            setIsFieldLoading(false);

            if (channels.length === 0) {
                setChannels(["No channels"]);
                setChannelsId(["No channels"]);
            } else {
                setChannels(channels.map((channel) => channel.name));
                setChannelsId(channelsId.map((channel) => channel.id));
            }
        }
    };

    const renderAsanaFields = () => {
        if (!selectedItem) {
            return null;
        }

        return (
            <div className="w-full max-w-[520px] mx-auto">
                {isWorkspaceLoaded ? (
                    selectedItem.parameters.map((param) => {
                        if (param.name === "workspace") {
                            return (
                                <SelectField
                                    key={param.name}
                                    label="Workspace"
                                    value={selectedWorkspaceName || ""}
                                    options={workspaces}
                                    onChange={(value) => handleWorkspaceChange(value)}
                                />
                            );
                        }

                        if (param.name === "project") {
                            const paramMap = projects.reduce<Record<string, string>>((acc, project, index) => {
                              acc[project] = projectsId[index];
                              return acc;
                            }, {});
                            return (
                                <SelectFieldId
                                    key={param.name}
                                    label="Project"
                                    value={itemFieldValues["project"] || ""}
                                    paramMap={paramMap}
                                    onChange={(value) => handleFieldChange("project", value, type)}
                                />
                            );
                        }

                        if (param.name === "assignee") {
                            const paramMap = assignees.reduce<Record<string, string>>((acc, assignee, index) => {
                              acc[assignee] = assigneesId[index];
                              return acc;
                            }, {});
                            return (
                                <SelectFieldId
                                    key={param.name}
                                    label="Assignee"
                                    value={itemFieldValues["assignee"] || ""}
                                    paramMap={paramMap}
                                    onChange={(value) => handleFieldChange("assignee", value, type)}
                                />
                            );
                        }

                        if (param.name === "tag") {
                            const paramMap = tags.reduce<Record<string, string>>((acc, tag, index) => {
                              acc[tag] = tagsId[index];
                              return acc;
                            }, {});
                            return (
                                <SelectFieldId
                                    key={param.name}
                                    label="Tag"
                                    value={itemFieldValues["tag"] || ""}
                                    paramMap={paramMap}
                                    onChange={(value) => handleFieldChange("tag", value, type)}
                                />
                            );
                        }


                        if (param.name === "due" || param.name === "notes" || param.name === "name") {
                            return (
                                <div className="pt-4 pb-2" key={param.name}>
                                    <h3 className="text-white font-[900] text-[1.6rem] pb-2" key={param.name}>{param.name}</h3>
                                    <input
                                        key={param.name}
                                        className="w-full text-black h-[4.8rem] px-4 border-[3px] rounded-[10px] border-[#EEEEEE] font-[900] text-[1.6rem] leading-7"
                                        type="text"
                                        value={itemFieldValues[param.name] || ""}
                                        onChange={(event) =>
                                            handleFieldChange(param.name, event.target.value, type)
                                        }
                                    />
                                </div>
                            );
                        }

                        return null;
                    })
                ) : (
                    <></>
                )}
            </div>
        );
    };

    const renderDiscordFields = () => {
    if (!selectedItem) {
        return null;
    }

    return (
        <div className="w-full max-w-[520px] mx-auto">
            {isServerLoaded ? (
                selectedItem.parameters.map((param) => {
                    if (param.name === "server") {
                        return (
                            <SelectField
                                key={param.name}
                                label="Server"
                                value={selectedServerName || ""}
                                options={servers}
                                onChange={(value) => handleServerChange(value)}
                            />
                        );
                    }

                    if (param.name === "channel") {
                        const paramMap = channels.reduce<Record<string, string>>((acc, channel, index) => {
                          acc[channel] = channelsId[index];
                          return acc;
                        }, {});
                        return (
                            <SelectFieldId
                                key={param.name}
                                label="Channel"
                                value={itemFieldValues["channel"] || ""}
                                paramMap={paramMap}
                                onChange={(value) => handleFieldChange("channel", value, type)}
                            />
                        );
                    }

                    if (param.name === "message" || param.name === "title") {
                        return (
                            <input
                                key={param.name}
                                className="w-full text-black h-[4.8rem] px-4 border-[3px] rounded-[10px] border-[#EEEEEE] font-[900] text-[1.6rem] leading-7"
                                type="text"
                                value={itemFieldValues[param.name] || ""}
                                onChange={(event) =>
                                    handleFieldChange(param.name, event.target.value, type)
                                }
                            />
                        );
                    }

                    return null;
                })
            ) : (
                <></>
            )}
        </div>
    );
};


    const renderFields = () => {
        if (asanaField) {
            return renderAsanaFields();
        }

        if (discordField) {
            return renderDiscordFields();
        }

        return (
            <DynamicFields
                fields={fields}
                values={itemFieldValues}
                onChange={(id, value) => handleFieldChange(id, value, type)}
            />
        );
    };

    return (
        <div style={{ backgroundColor: selectedService?.color }} className="min-h-screen">
            <div
                className="flex items-center w-[110%] lg:px-10 px-2 py-5 border-b-[1px] border-[#dadada]"
            >
                <div>
                    <ActionButton
                        text="Back"
                        textColor="text-white"
                        borderColor="border-white"
                        bgColor={`bg-[${selectedService?.color}]`}
                        onClick={() => goBack()}
                    />
                </div>

                <h1 className="text-white text-[1.3rem] sm:text-[3.2rem] font-extrabold text-center w-full">
                    {type === "action" ? "Complete action fields" : "Complete reaction fields"}
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
                        className="text-white text-[2.2rem] sm:text-[3.2rem] font-bold mt-6 mb-8 text-center lg:px-0 px-4"
                    >
                        {selectedItem?.name}
                    </h2>
                </div>

                <div className="flex flex-col items-center lg:px-0 px-4">
                    {isLoading ? (
                        <></>
                    ) : (
                        <>
                            {renderFields()}

                            {isFieldLoading ? (
                                <></>
                            ) : (
                                <TriggerButton
                                    text={type === "action" ? "Create Action" : "Create Reaction"}
                                    textColor="text-black"
                                    bgColor="bg-white"
                                    hoverColor="hover:bg-[#EEEEEE]"
                                    onClick={() => goToStep(Steps.HOME)}
                                />
                            )}
                        </>
                    )}
                </div>
                <div className="pt-5"></div>

        </div>
    );
};

export default TriggerActionReactionStep;

