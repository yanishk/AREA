import { useEffect, useState } from "react";
import ErrorMessage from "../components/ErrorMessage"
import Navbar from "../components/navbar/Navbar"
import {useLocation, useNavigate } from "react-router-dom";

const HomePage = () => {
    const location = useLocation();
    const [error, setError] = useState("");
    const [errorTrigger, setErrorTrigger] = useState(0);
    const navigate = useNavigate();

    useEffect(() => {
        if (location.state?.error) {
            setError(location.state.error);
        }
        if (location.state?.errorTrigger) {
            setErrorTrigger(location.state.errorTrigger);
        }
    }, [location.state]);

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