import axios from "axios";

export const handleOAuthLogin = async (serviceName: string) => {
    try {
        const result = await axios.get(
            `${import.meta.env.VITE_API_URL}authentication?service=${serviceName}&callbacktype=login&apptype=web`,
            {
                withCredentials: true,
            }
        );

        const authUrl = result.data["auth-url"];
        if (!authUrl) {
            console.error("Authentication URL not received");
            return;
        }

        sessionStorage.setItem("oauth2-login", serviceName)
        window.location.href = authUrl;
    } catch (error) {
        console.error(`Error during ${serviceName} authentication:`, error);
    }
};
