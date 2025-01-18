import { useEffect, useState } from "react";
import Navbar from "../components/navbar/Navbar"
import axios from "axios";
import { useNavigate, useParams } from "react-router-dom";
import ActionButton from "../components/buttons/ActionButton";
import MaxWidthWrapper from "../components/MaxWidthWrapper";

type serviceType = {
    name: string;
    color: string;
    logo: string;
    description: string;
}

type serviceItemType = {
    id: string;
    name: string,
    description: string,
}

const ServiceDetailsPage = () => {
    const { serviceName } = useParams();
    const navigate = useNavigate();
    const [service, setService] = useState<serviceType | null>(null);
    const [serviceActions, setServiceActions] = useState<serviceItemType[] | null>(null);
    const [serviceReactions, setServiceReactions] = useState<serviceItemType[] | null>(null);

    const getServiceActions = async () => {
        const result = await axios.get(`${import.meta.env.VITE_API_URL}${service?.name}/actions`, {
            withCredentials: true,
        });


        if (!result?.data?.actions) {
            return;
        }

        setServiceActions(result.data.actions);
    }

    const getServiceReactions = async () => {
        const result = await axios.get(`${import.meta.env.VITE_API_URL}${service?.name}/reactions`, {
            withCredentials: true,
        });


        if (!result?.data?.reactions) {
            return;
        }

        setServiceReactions(result.data.reactions);
    }

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
                        error: "Service doesn't exist or this account doesn't have access to it",
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

    useEffect(() => {
        if (service) {
            getServiceActions();
            getServiceReactions();
        }
    }, [service])

    return (
        <>
            <div
                style={{background: service?.color}}
            >
                <Navbar
                    isWhiteMode={true}
                />

                <div
                    className="flex items-center justify-between px-8 lg:py-5 pb-5"
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

                                <p className="text-white text-[1.1rem] tracking-wider mt-[10px] mb-[35px] lg:px-0 px-5">
                                    {service?.description}
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
                                    {service?.name != "SMS" && service?.name != "Email" && service?.name != "Time & Date" && (
                                        <div className="ml-[40px]">
                                                <ActionButton
                                                    text={`Visit ${service?.name}`}
                                                    textColor="text-white"
                                                    borderColor="border-white"
                                                    bgColor={`bg-[${service?.color}]`}
                                                    onClick={() =>
                                                        window.open(`https://${serviceName}.com`, "_blank")
                                                    }
                                                />
                                        </div>
                                    )}
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

            <MaxWidthWrapper
                maxWidth="1060px"
            >

                {serviceActions && serviceActions.length > 0 ? (
                    <div
                        className="mt-[60px]"
                    >
                        <h3
                            className="text-[#222222] text-[1.5rem] font-[800] mb-[25px] px-3 lg:px-2"
                        >
                            Actions
                        </h3>
                        <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-3 gap-5 ml-2 mr-2 mb-[20px]">
                                {serviceActions.map((action) => (
                                    <div
                                        key={action.id}
                                        style={{background: service?.color}}
                                        className="h-[220px] flex flex-col items-left justify-center px-10 rounded-[10px] relative overflow-hidden"
                                    >
                                        <p className="text-white text-[1.4em] font-[800] mb-2 whitespace-normal">
                                            {action.name}
                                        </p>
                                        <p className="text-white text-[18px] overflow-hidden text-ellipsis whitespace-normal line-clamp-3">
                                            {action.description}
                                        </p>
                                    </div>
                                ))}
                        </div>
                    </div>
                ) : (
                    <></>
                )}

                {serviceReactions && serviceReactions.length > 0 ? (
                    <div
                        className="mt-[60px]"
                    >
                        <h3
                            className="text-[#222222] text-[1.5rem] font-[800] mb-[25px] px-3 lg:px-2"
                        >
                            Reactions
                        </h3>
                        <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-3 gap-5 ml-2 mr-2 mb-[20px]">
                                {serviceReactions.map((reaction) => (
                                    <div
                                        key={reaction.id}
                                        style={{background: service?.color}}
                                        className="h-[220px] flex flex-col items-left justify-center px-10 rounded-[10px] relative overflow-hidden"
                                    >
                                        <p className="text-white text-[1.4em] font-[800] mb-2 whitespace-normal">
                                            {reaction.name}
                                        </p>
                                        <p className="text-white text-[18px] overflow-hidden text-ellipsis whitespace-normal line-clamp-3">
                                            {reaction.description}
                                        </p>
                                    </div>
                                ))}
                        </div>
                    </div>
                ) : (
                    <></>
                )}
            </MaxWidthWrapper>
        </>
    )
}

export default ServiceDetailsPage