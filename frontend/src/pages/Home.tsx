import { useEffect, useState } from "react";
import ErrorMessage from "../components/errors/ErrorMessage"
import Navbar from "../components/navbar/Navbar"
import {useLocation, useNavigate, useSearchParams } from "react-router-dom";
import axios from "axios";
import { useAuth } from "./AuthContext";

const HomePage = () => {
    const location = useLocation();
    const [error, setError] = useState("");
    const [errorTrigger, setErrorTrigger] = useState(0);
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const { login } = useAuth();

    useEffect(() => {
        if (location.state?.error) {
            setError(location.state.error);
        }
        if (location.state?.errorTrigger) {
            setErrorTrigger(location.state.errorTrigger);
        }
    }, [location.state]);

    const getAuthCode = async () => {
        const codeAuth = searchParams.get("code");
        const serviceName = sessionStorage.getItem("oauth2-login")

        if (!codeAuth || !serviceName) {
            return;
        }

        await axios.post(
            `${import.meta.env.VITE_API_URL}login-callback/${serviceName}/web/?code=${codeAuth}`,
            {
            },
            {
                withCredentials: true,
            }
        );
        login();
        navigate(location.pathname, { replace: true });
    }

    useEffect(() => {
        getAuthCode();
    }, [useSearchParams]);

    const handleErrorClose = () => {
        setError("");
        setErrorTrigger(0);
        navigate(location.pathname, { replace: true, state: {} });
    };

    return (
        <>
            <Navbar bgColor="bg-white"/>
            <ErrorMessage message={error} trigger={errorTrigger} onClose={handleErrorClose} />
        </>
    )
}

export default HomePage