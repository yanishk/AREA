import axios from "axios";

const isAuthenticated = async () => {
    try {
        const response = await axios.get(`${import.meta.env.VITE_API_URL}user`, {
            withCredentials: true,
        });

        if (response.status === 200) {
            return true;
        }

        return false;
    } catch (error) {
        console.error("Failed to get user's informations : ", error);
        return false;
    }
};

export default isAuthenticated