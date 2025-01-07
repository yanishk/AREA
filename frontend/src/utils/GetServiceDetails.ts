import axios from "axios";

const getServiceDetails = async (id: string, type: string) => {
    try {
        const request = type === "action"
            ? `${import.meta.env.VITE_API_URL}services/action/${id}`
            : `${import.meta.env.VITE_API_URL}services/reaction/${id}`;

        const result = await axios.get(request, {
            withCredentials: true,
        });

        return result.data.service;
    } catch (err) {
        console.error(`Failed to get details`, err);
        return null;
    }
};

export default getServiceDetails