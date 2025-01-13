import { useEffect, useState } from "react";
import Navbar from "../components/navbar/Navbar"
import axios from "axios";
import { useNavigate, useParams } from "react-router-dom";
import ActionButton from "../components/buttons/ActionButton";
import MaxWidthWrapper from "../components/MaxWidthWrapper";

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
        console.log("serviceNameici: ", serviceName);
        fetchServices();
    }, []);

    return (
        <>
        <div
            style={{background: service?.color}}
        >
            <Navbar
                isWhiteMod={true}
            />

            <div
                className="flex items-center justify-between px-10 py-5"
            >
                <div>
                    <ActionButton
                        text="Back"
                        textColor="text-white"
                        borderColor="border-white"
                        bgColor={`bg-[${service?.color}]`}
                        onClick={() => navigate(-1)}
                    />
                </div>
            </div>
            <div
                    className="flex items-center justify-center mt-[28px]"
                >
                    <MaxWidthWrapper
                        maxWidth="750px"
                    >
                        <div
                            className="flex flex-col items-center text-center w-full"
                        >
                            <img
                                src={`${import.meta.env.VITE_HOST_URL}${service?.logo}`}
                                alt={service?.name}
                                className="w-[120px] h-[120px]"
                            />

                            <h1
                                className="text-white text-[3.6rem] font-[900] mt-[55px]"
                            >
                                {service?.name} {' '}
                                <span>
                                    integrations
                                </span>
                            </h1>

                            <p className="text-white text-[1.1rem] tracking-wider mt-[10px] mb-[35px]">
                                Spotify is a digital music service that gives you access to millions of songs. Applets can help you save your Discover Weekly and Release Radar playlists, share your favorite tunes, and much more.
                            </p>

                            <div className="flex flex-row items-center justify-center w-full mb-[60px]">
                                <div>
                                    <ActionButton
                                        text={"Create"}
                                        textColor="text-black"
                                        borderColor="border-white"
                                        bgColor={`bg-[white]`}
                                        onClick={() => navigate("/create")}
                                    />
                                </div>
                                <div className="ml-[40px]">
                                    {service?.name != "SMS" && service?.name != "Email" && service?.name != "Time & Date" && (
                                        <ActionButton
                                            text={`Visit ${serviceName}`}
                                            textColor="text-white"
                                            borderColor="border-white"
                                            bgColor={`bg-[${service?.color}]`}
                                            onClick={() =>
                                                window.open(`https://${serviceName}.com`, "_blank")
                                            }
                                        />
                                    )}
                                </div>
                            </div>
                        </div>
                    </MaxWidthWrapper>
            </div>
        </div>

        <div className="flex items-center justify-center mt-[50px]">
            <p className="text-black text-[20px] pb-[5px] font-bold border-b-4 border-black tracking-wide">
                Details
            </p>
        </div>
        </>
    )
}

export default ServiceDetailsPage