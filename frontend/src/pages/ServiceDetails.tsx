import { useEffect, useState } from "react";
import Navbar from "../components/navbar/Navbar"
import axios from "axios";
import { useNavigate, useParams } from "react-router-dom";

const ServiceDetailsPage = () => {
    const { serviceName } = useParams();
    const navigate = useNavigate();
    const [service, setService] = useState<any | null>(null);

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

            const matchedService = result.data.find((service: any) => service.name.toLowerCase() === serviceName?.toLowerCase());

            if (!matchedService) {
                navigate("/explore", {
                    replace: true,
                    state: {
                        error: `Service doesn't exist or this account doesn't have access to it`,
                        errorTrigger: 1,
                    },
                });
            } else {
                setService(matchedService);

                window.history.replaceState(
                    null,
                    '',
                    `/${serviceName?.toLowerCase()}`
                );
            }
        }

        fetchServices();
    }, []);

    return (
        <>
            <Navbar bgColor="bg-white"/>
            <div
                    className="flex items-center justify-center mt-10"
                >
                    <div
                        className="flex flex-col items-center text-center w-full"
                    >
                        <h1
                            className="text-5xl text-[#222222] font-[900]"
                        >
                            Service
                        </h1>
                    </div>
            </div>

            {service && (
                <div>
                    <h2 className="text-3xl font-bold">{service.name}</h2>
                    <img
                        src={`${import.meta.env.VITE_HOST_URL}${service.logo}`}
                        alt={service.name}
                        className="bg-red-500 w-32 h-32 mx-auto"
                    />
                </div>
            )}

        </>
    )
}

export default ServiceDetailsPage