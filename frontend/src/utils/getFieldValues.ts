import axios from "axios";
import { Parameters } from "./ActionReactionParameters";
import { getAsanaWorkspaces } from "./AsanaUtils";

export const gitlabNameToIdMap: Record<string, number> = {};

interface GitlabParams {
    id: number;
    name: string;
}

interface GithubParams {
    full_name: string;
}

const getGitlabFieldValues = async (route: string): Promise<{ data: { projects: GitlabParams[] } }> => {
    try {
        const newRoute = route.slice(1);
        const userProject = await axios.get(
            `${import.meta.env.VITE_API_URL}${newRoute}`,
            {
                withCredentials: true
            }
        );

        return userProject;
    } catch (error) {
        console.error("Error with gitlab fields values:", error);
        return { data: { projects: [] } };
    }
}

const getGithubFieldValues = async (route: string): Promise<{ data: { repositories: GithubParams[] } }> => {
    try {
        const newRoute = route.slice(1);
        const userRepo = await axios.get(
            `${import.meta.env.VITE_API_URL}${newRoute}`,
            {
                withCredentials: true
            }
        );

        return userRepo;
    } catch (error) {
        console.error("Error with github fields values:", error);
        return { data: { repositories: [] } };
    }
}

const handleAsanaService = async (param: Parameters): Promise<string[] | undefined> => {
    if (param.route === "/asana/user/workspaces") {
        const workspaces = await getAsanaWorkspaces();
        return workspaces.map((workspace) => workspace.name);
    }
}

const getFieldValues = async (param: Parameters): Promise<string[] | undefined> => {

    if (param.route) {
        const serviceName = param.route.split("/")[1];

        if (serviceName === "gitlab") {
            const result = await getGitlabFieldValues(param.route);

            if (result.data && Array.isArray(result.data.projects)) {
                result.data.projects.forEach((item) => {
                    gitlabNameToIdMap[item.name] = item.id;
                });

                const values = result.data.projects.map((item) => item.name);
                return values;
            }
        }

        if (serviceName === "github") {
            const result = await getGithubFieldValues(param.route);

            if (result.data && Array.isArray(result.data.repositories)) {
                const values = result.data.repositories.map((item) => item.full_name);
                return values;
            }
        }

        if (serviceName === "asana") {
            return handleAsanaService(param);
        }
    }
    if (param.values[0] === null) {
        return []
    }
    if (param.isexhaustive === true && param.values.length > 1) {
        return param.values;
    }

    if (param.isexhaustive === false && param.values.length === 1) {
        const value = param.values[0];
        const [start, end] = value.split("-").map(Number);
        const newArray = Array.from({ length: end - start + 1 }, (_, i) => (start + i).toString());
        return newArray;
    }

    return undefined;
};

export default getFieldValues
