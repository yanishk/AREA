import axios from "axios";

export const getUserName = async () => {
    try {
        const result = await axios.get(`${import.meta.env.VITE_API_URL}user`, {
            withCredentials: true,
        });

        const email = result.data.user.email;

        const username = email.split('@')[0].replace(/\./g, '');

        return username;
    } catch (error) {
        console.error("Failed to get email : ", error);
    }
};