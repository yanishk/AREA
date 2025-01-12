import { useEffect, useState } from "react";
import Navbar from "../components/navbar/Navbar"
import axios from "axios";
import MaxWidthWrapper from "../components/MaxWidthWrapper";
import { useNavigate, useParams } from "react-router-dom";

const ServiceDetailsPage = () => {
    const { serviceName } = useParams();
    const navigate = useNavigate();
    const [service, setService] = useState<any[]>([]);

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

            // setServices(result.data);

            const isServiceExist = result.data.some((service: any) => service.name.toLowerCase() === serviceName?.toLowerCase());

            if (!isServiceExist) {
                navigate("/explore", {
                    replace: true,
                    state: {
                        error: `Service doesn't exist or this account doesn't have access to it`,
                        errorTrigger: 1,
                    },
                });
            } else {
                window.history.replaceState(
                    null,
                    '',
                    `/${serviceName?.toLowerCase()}`
                );
            }

            console.log("service ici : ", service);
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

            {/* <MaxWidthWrapper maxWidth="1140px">
                <div className="mt-16 ml-[35px]">
                    <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-3">
                        {services?.filter(service => service.name !== "Google").map((service) => (
                            <button
                                key={service.id}
                                style={{
                                    backgroundColor: service.color,
                                    boxShadow: '0 1.5px 6px 0 rgba(0,0,0,.24)'
                                }}
                                className={`h-[380px] flex flex-col items-center justify-center cursor-pointer rounded-[10px] mr-6 mb-5 hover:opacity-90 relative`}
                            >
                                <img src={`${import.meta.env.VITE_HOST_URL}${service.logo}`} className="w-[120px] h-[120px]" />
                                <p className="mt-[35px] text-white text-[1.215em] font-extrabold">{service.name}</p>
                            </button>
                        ))}
                    </div>
                </div>
            </MaxWidthWrapper> */}
        </>
    )
}

export default ServiceDetailsPage