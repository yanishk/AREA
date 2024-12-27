import axios from "axios";

export const handleOAuthLogin = async (serviceName: string) => {
    try {
        const result = await axios.get(
            `${import.meta.env.VITE_API_URL}authentication/${serviceName}/login/web`,
            {
                withCredentials: true,
            }
        );

        const authUrl = result.data["auth-url"];
        if (!authUrl) {
            console.error("Error: Authentication URL not received");
            return;
        }

        const state = JSON.stringify({ service: serviceName });
        const authUrlWithService = `${authUrl}&state=${encodeURIComponent(state)}`;

        window.location.href = authUrlWithService;
    } catch (error) {
        console.error(`Error during ${serviceName} authentication:`, error);
    }
};
