import { useEffect, useState } from "react";
import { Parameters } from "../../utils/ActionReactionParameters";
import Steps from "../../utils/CreateStepsType";
import { fieldType } from "../../utils/fieldType";
import TriggerButton from "../buttons/TriggerButton";
import DynamicFields from "../inputs/DynamicFields";
import ActionButton from "../buttons/ActionButton";
import { getAsanaAssignees, getAsanaProjects, getAsanaTags, getAsanaWorkspaces, getWorkspaceIdByName, setSelectedWorkspace } from "../../utils/AsanaUtils";
import SelectField from "../inputs/SelectField";
import getFieldValues from "../../utils/getFieldValues";

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
    const [isFieldLoading, setIsFieldLoading] = useState(false);
    const [advancedField, setAdvancedField] = useState(false);
    const [workspaces, setWorkspaces] = useState<string[]>([]);
    const [projects, setProjects] = useState<string[]>([]);
    const [assignees, setAssignees] = useState<string[]>([]);
    const [tags, setTags] = useState<string[]>([]);
    const [selectedWorkspaceName, setSelectedWorkspaceName] = useState<string | null>(null);

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
        if (selectedService?.name != "Asana") {
            return;
        }
        const fetchWorkspaces = async () => {
            const workspaces = await getAsanaWorkspaces();

            setWorkspaces(workspaces.map((ws) => ws.name));
            if (workspaces.length > 0) {
                setSelectedWorkspace(workspaces[0].gid, workspaces[0].name);
                await handleWorkspaceChange(workspaces[0].name);
            }
        };

        fetchWorkspaces();
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
            setAdvancedField(true);
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

            setIsWorkspaceLoaded(true);
            setIsFieldLoading(false);

            if (projects.length === 0) {
                setProjects(["No projects"]);
            } else {
                setProjects(projects.map((project) => project.name));
            }

            if (assignees.length === 0) {
                setAssignees(["No assignees"]);
            } else {
                setAssignees(assignees.map((assignee) => assignee.name));
            }

            if (tags.length === 0) {
                setTags(["No tags"]);
            } else {
                setTags(tags.map((tag) => tag.name));
            }
        }
    };

    return (
        <div style={{ backgroundColor: selectedService?.color }} className="min-h-screen">
            <div
                className="flex items-center justify-between px-10 py-5 border-b-[1px] border-[#dadada]"
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

                <h1 className="text-white text-[3.2rem] font-extrabold text-center flex-1">
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
                        className="text-white text-[52px] font-bold mt-6 mb-8"
                    >
                        {selectedItem?.name}
                    </h2>
                </div>

                <div
                    className="flex flex-col items-center"
                >
                    {isLoading ? (
                        <></>
                    ) : (
                        <>
                            {advancedField ? (
                                <div className="w-full max-w-[520px] mx-auto">
                                    <SelectField
                                        label="Workspace"
                                        value={selectedWorkspaceName || ""}
                                        options={workspaces}
                                        onChange={(value) => handleWorkspaceChange(value)}
                                    />

                                    {isWorkspaceLoaded ? (
                                        <>
                                            <SelectField
                                                label="Project"
                                                value={itemFieldValues["project"] || ""}
                                                options={projects}
                                                onChange={(value) => handleFieldChange("project", value, type)}
                                            />

                                            <SelectField
                                                label="Assignee"
                                                value={itemFieldValues["assignee"] || ""}
                                                options={assignees}
                                                onChange={(value) => handleFieldChange("assignee", value, type)}
                                            />

                                            <SelectField
                                                label="Tag"
                                                value={itemFieldValues["tag"] || ""}
                                                options={tags}
                                                onChange={(value) => handleFieldChange("tag", value, type)}
                                            />
                                        </>
                                    ) : (
                                        <></>
                                    )}
                                </div>
                            ) : (
                                <DynamicFields
                                    fields={fields}
                                    values={itemFieldValues}
                                    onChange={(id, value) => handleFieldChange(id, value, type)}
                                />
                            )}
                            {isFieldLoading ? (
                                <></>
                            ) : (
                                <TriggerButton
                                    text="Create Trigger"
                                    textColor="text-black"
                                    bgColor="bg-white"
                                    hoverColor="hover:bg-[#EEEEEE]"
                                    onClick={() => goToStep(Steps.HOME)}
                                />
                            )}
                        </>
                    )}
                </div>

        </div>
    );
};

export default TriggerActionReactionStep;

