import axios from "axios";
import { Parameters } from "./ActionReactionParameters";
import { getAsanaWorkspaces } from "./AsanaUtils";
import { getDiscordServers } from "./DiscordUtils";

interface GitlabParams {
    id: number;
    name: string;
}

interface GithubParams {
    full_name: string;
}

const getGitlabFieldValues = async (route: string): Promise<{ data: { projects: GitlabParams[] } }> => {
    try {
        const slicedRoute = route.slice(1);
        const userProject = await axios.get(
            `${import.meta.env.VITE_API_URL}${slicedRoute}`,
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

const isServiceConnected = async (Servicename: string): Promise<boolean> => {
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

const getGithubFieldValues = async (route: string): Promise<{ data: { repositories: GithubParams[] } }> => {
    try {
        const result = await isServiceConnected("Github")
        if (!result) {
            return { data: { repositories: [] } };
        }
        const slicedRoute = route.slice(1);
        const userRepo = await axios.get(
            `${import.meta.env.VITE_API_URL}${slicedRoute}`,
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

const handleService = async (route: string, param: Parameters, serviceName: string): Promise<string[] | undefined> => {
        if (serviceName === "gitlab") {
            const result = await getGitlabFieldValues(route);

            return result.data.projects.map((item) => item.name);
        }

        if (serviceName === "github") {
            const result = await getGithubFieldValues(route);

            if (result.data?.repositories) {
                return result.data.repositories.map((item) => item.full_name);
            }

            return undefined;
        };

        if (serviceName === "asana") {
            if (param.route === "/asana/user/workspaces") {
                const workspaces = await getAsanaWorkspaces();
                return workspaces.map((workspace) => workspace.name);
            }
        }

        if (serviceName === "discord") {
            if (param.route === "/discord/user/servers") {
                const servers = await getDiscordServers();
                return servers.map((server) => server.name);
            }
        }

        return undefined;
    };

const handleExhaustiveValues = (param: Parameters): string[] | undefined => {
    if (param.isexhaustive && param.values.length > 1) {
        return param.values;
    }

    if (!param.isexhaustive && param.values.length === 1) {
        const value = param.values[0];
        const [start, end] = value.split("-").map(Number);
        return Array.from({ length: end - start + 1 }, (_, i) => (start + i).toString());
    }

    return undefined;
};

const getFieldValues = async (param: Parameters): Promise<string[] | undefined> => {
    if (param.route) {
        const serviceName = param.route.split("/")[1];

        switch (serviceName) {
            case "gitlab":
                return handleService(param.route, param, "gitlab");
            case "github":
                return handleService(param.route, param, "github");
            case "asana":
                return handleService(param.route, param, "asana");
            case "discord":
                return handleService(param.route, param, "discord");
        }
    }

    if (param.values[0] === null) {
        return [];
    }

    return handleExhaustiveValues(param);
};

export default getFieldValues
