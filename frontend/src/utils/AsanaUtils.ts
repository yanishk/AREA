import axios from "axios";

export interface Workspace {
    gid: string;
    name: string;
}

export interface AsanaItem {
    gid: string;
    name: string;
}

let Workspaces: Workspace[] = [];
let selectedWorkspaceId: string | null = null;
let selectedWorkspaceName: string | null = null;

const setDefaultWorkspace = (workspaces: Workspace[]) => {
    if (workspaces.length > 0) {
        selectedWorkspaceId = workspaces[0].gid;
        selectedWorkspaceName = workspaces[0].name;
    }
};

export const getAsanaWorkspaces = async (): Promise<Workspace[]> => {
    try {
        const response = await axios.get(`${import.meta.env.VITE_API_URL}asana/user/workspaces`, {
            withCredentials: true,
        });

        if (response.data && response.data.workspaces && Array.isArray(response.data.workspaces.data)) {
            Workspaces = response.data.workspaces.data.map((workspace: { gid: string; name: string }) => ({
                gid: workspace.gid,
                name: workspace.name,
            }));

            setDefaultWorkspace(Workspaces);

            return Workspaces;
        }

        console.warn("No workspaces founded");
        return [];
    } catch (error) {
        console.error("Error fetching Asana workspaces : ", error);
        return [];
    }
};

export const getWorkspaceIdByName = async (workspaceName: string): Promise<string | null> => {
    const workspaces = await getAsanaWorkspaces();
    const workspace = workspaces.find((ws) => ws.name === workspaceName);
    return workspace ? workspace.gid : null;
};

export const getSelectedWorkspaceId = (): string | null => {
    return selectedWorkspaceId;
};

export const getSelectedWorkspaceName = (): string | null => {
    return selectedWorkspaceName;
};

export const setSelectedWorkspace = (workspaceId: string, workspaceName: string) => {
    selectedWorkspaceId = workspaceId;
    selectedWorkspaceName = workspaceName;
};

export const getAsanaProjects = async (workspaceId: string): Promise<AsanaItem[]> => {
    try {
        const response = await axios.get(
            `${import.meta.env.VITE_API_URL}asana/workspace/projects?id=${workspaceId}`,
            {
                withCredentials: true,
            }
        );

        if (response.data && response.data.projects && Array.isArray(response.data.projects.data)) {
            return response.data.projects.data.map((project: { gid: string; name: string }) => ({
                gid: project.gid,
                name: project.name,
            }));
        }

        console.warn("No projects founded : ", workspaceId);
        return [];
    } catch (error) {
        console.error("Error fetching Asana projects : ", error);
        return [];
    }
};

export const getAsanaAssignees = async (workspaceId: string): Promise<AsanaItem[]> => {
    try {
        const response = await axios.get(
            `${import.meta.env.VITE_API_URL}asana/workspace/assignees?id=${workspaceId}`,
            {
                withCredentials: true,
            }
        );

        if (response.data && response.data.assignees && Array.isArray(response.data.assignees.data)) {
            return response.data.assignees.data.map((assignee: { gid: string; name: string }) => ({
                gid: assignee.gid,
                name: assignee.name,
            }));
        }

        console.warn("No assignees founded : ", workspaceId);
        return [];
    } catch (error) {
        console.error("Error fetching Asana assignees : ", error);
        return [];
    }
};

export const getAsanaTags = async (workspaceId: string): Promise<AsanaItem[]> => {
    try {
        const response = await axios.get(
            `${import.meta.env.VITE_API_URL}asana/workspace/tags?id=${workspaceId}`,
            {
                withCredentials: true,
            }
        );

        if (response.data && response.data.tags && Array.isArray(response.data.tags.data)) {
            return response.data.tags.data.map((tag: { gid: string; name: string }) => ({
                gid: tag.gid,
                name: tag.name,
            }));
        }

        console.warn("No tags founded : ", workspaceId);
        return [];
    } catch (error) {
        console.error("Error fetching Asana tags :", error);
        return [];
    }
};
