import { useEffect, useState } from "react";
import Navbar from "../components/navbar/Navbar"
import {useLocation, useNavigate, useSearchParams } from "react-router-dom";
import axios from "axios";
import { useAuth } from "./AuthContext";
import MaxWidthWrapper from "../components/MaxWidthWrapper";
import MessageBox from "../components/notification/MessageBox";

const ExplorePage = () => {
    const location = useLocation();
    const [notification, setNotification] = useState("");
    const [notificationTrigger, setNotificationTrigger] = useState(0);
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const [services, setServices] = useState<any[]>([]);
    const { login } = useAuth();

    useEffect(() => {
        if (location.state?.notification) {
            setNotification(location.state.notification);
        }
        if (location.state?.notificationTrigger) {
            setNotificationTrigger(location.state.notificationTrigger);
        }
    }, [location.state]);

    const getAuthCode = async () => {
        const codeAuth = searchParams.get("code");
        const serviceName = sessionStorage.getItem("oauth2-login")

        if (!codeAuth || !serviceName) {
            return;
        }

        await axios.post(
            `${import.meta.env.VITE_API_URL}login-callback/?code=${codeAuth}`,
            {
                service: serviceName,
                apptype: "web"
            },
            {
                withCredentials: true
            }
        );
        login();
        navigate(location.pathname, { replace: true });
    }

    useEffect(() => {
        getAuthCode();
    }, [searchParams]);

    useEffect(() => {
        const fetchServices = async () => {
            const result = await axios.get(
                `${import.meta.env.VITE_API_URL}services`,
                {
                    withCredentials: true
                }
            );

            if (!result) {
                console.error("Error fetching services");
                return;
            }

            setServices(result.data);
        }

        fetchServices();
    }, []);

    const handleErrorClose = () => {
        setNotification("");
        setNotificationTrigger(0);
        navigate(location.pathname, { replace: true, state: {} });
    };

    return (
        <>
            <Navbar bgColor="bg-white"/>

            <MessageBox message={notification} trigger={notificationTrigger} onClose={handleErrorClose} type="validation" timeout={3000} />

            <div
                    className="flex items-center justify-center mt-10"
                >
                    <div
                        className="flex flex-col items-center text-center w-full"
                    >
                        <h1
                            className="text-5xl text-[#222222] font-[900]"
                        >
                            Explore
                        </h1>
                    </div>
            </div>

            <MaxWidthWrapper maxWidth="1140px">
                <div className="mt-16 ml-6">
                    <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-3">
                        {services?.filter(service => service.name !== "Google").map((service) => (
                            <button
                                key={service.id}
                                style={{
                                    backgroundColor: service.color,
                                    boxShadow: '0 1.5px 6px 0 rgba(0,0,0,.24)'
                                }}
                                className={`h-[380px] flex flex-col items-center justify-center cursor-pointer rounded-[10px] mr-6 mb-5 hover:opacity-90 relative`}
                                onClick={() => navigate(`/${service.name}`)}
                            >
                                <img alt="icon" src={`${import.meta.env.VITE_HOST_URL}${service.logo}`} className="w-[120px] h-[120px]" />
                                <p className="mt-[35px] text-white text-[1.215em] font-extrabold">{service.name}</p>
                            </button>
                        ))}
                    </div>
                </div>
            </MaxWidthWrapper>
        </>
    )
}

export default ExplorePage