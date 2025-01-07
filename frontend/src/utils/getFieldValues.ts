import axios from "axios";
import { Parameters } from "./ActionReactionParameters";

interface GitlabParams {
    id: number;
    name: string;
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

const getFieldValues = async (param: Parameters): Promise<string[] | undefined> => {
    if (param.route) {
        const result = await getGitlabFieldValues(param.route);

        if (result.data && Array.isArray(result.data.projects)) {
            const values = result.data.projects.map((item) => item.name);
            return values;
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
