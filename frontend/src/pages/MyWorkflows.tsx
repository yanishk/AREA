import axios from "axios";
import MaxWidthWrapper from "../components/MaxWidthWrapper"
import Navbar from "../components/navbar/Navbar"
import { useEffect, useState } from "react";
import { getUserName } from "../utils/GetUserName";
import { useLocation, useNavigate } from "react-router-dom";
import getServiceDetails from "../utils/GetServiceDetails";
import ButtonValidation from "../components/buttons/ButtonValidation";
import MessageBox from "../components/notification/MessageBox";

const MyWorkflow = () => {

    const [userWorkflows, setUserWorkflows] = useState<any[]>([]);
    const [userName, setUserName] = useState("You");
    const [isLoading, setIsLoading] = useState(true);
    const [notification, setNotification] = useState("");
    const [notificationTrigger, setNotificationTrigger] = useState(0);
    const navigate = useNavigate();
    const location = useLocation();

    const getUserWorkflows = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}workflows`, {
                withCredentials: true,
            });

            if (!result) {
                console.error("Error fetching workflows");
                return;
            }

            const workflows = result.data.workflows;

            if (!result.data.workflows) {
                setUserWorkflows([]);
                setIsLoading(false);
                return;
            }

            const workflowsDetails = await Promise.all(
                workflows.map(async (workflow: any) => {
                    const actionDetails = await getServiceDetails(workflow.actionid, "action");
                    const reactionDetails = await getServiceDetails(workflow.reactionid, "reaction");

                    return {
                        ...workflow,
                        actioncolor: actionDetails?.color,
                        actionlogo: actionDetails?.logo,
                        reactionlogo: reactionDetails?.logo,
                    };
                })
            );

            setUserWorkflows(workflowsDetails);
            setIsLoading(false);
        } catch (err) {
            console.error("Error fetching workflows:", err);
        }
    };

    useEffect(() => {
        if (location.state?.notification) {
            setNotification(location.state.notification);
        }
        if (location.state?.notificationTrigger) {
            setNotificationTrigger(location.state.notificationTrigger);
        }
    }, [location.state]);

    useEffect(() => {
        const fetchUserName = async () => {
            try {
                const username = await getUserName();
                setUserName(username);
            } catch (error) {
                console.error("Failed to get username: ", error);
            }
        }
        getUserWorkflows();
        fetchUserName();
    }, []);

    const sliceText = (text: string, maxLen: number) => {
        if (text.length > maxLen) {
            return text.slice(0, maxLen) + "...";
        }
        return text;
    }

    const handleNotificationClose = () => {
        setNotification("");
        setNotificationTrigger(0);
        navigate(location.pathname, { replace: true, state: {} });
    };

    const renderContent = () => {
        if (isLoading) {
            return null;
        }

        if (userWorkflows.length === 0) {
            return (
                <div
                    className="flex flex-col items-center justify-center min-h-[75vh] text-center"
                >
                    <h1
                        className="text-[#222222] text-[2.96666667rem] font-[900]"
                    >
                        Start connecting your world
                    </h1>
                    <p
                        className="text-[#222222] mt-[30px] text-[20px] mb-[50px] font-[500]"
                    >
                        Save time and money by making the internet work for you!
                    </p>
                    <MaxWidthWrapper
                        maxWidth="420px"
                    >
                        <ButtonValidation
                            text="Create your workflow"
                            onClick={() => navigate("/create")}
                        />
                    </MaxWidthWrapper>
                </div>
            )
        }

        return (
            <>
                <div
                    className="flex items-center justify-center mt-10"
                >
                    <div
                        className="flex flex-col items-center text-center w-full"
                    >
                        <h1
                            className="text-5xl text-[#222222] font-[900]"
                        >
                            My Workflows
                        </h1>

                    </div>
                </div>
                <div className="mt-16">
                    <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-3 ml-3">
                        {userWorkflows?.map((workflow) => (
                            <button
                                key={workflow.id}
                                style={{
                                    backgroundColor: workflow.actioncolor,
                                    boxShadow: '0 1.5px 6px 0 rgba(0,0,0,.24)'
                                }}
                                className={`h-[420px] flex flex-col items-start pt-2 cursor-pointer rounded-[10px] mr-3 mb-5 hover:opacity-90 relative ${
                                    workflow.isactivated ? "opacity-100" : "opacity-50 hover:opacity-30"
                                }`}
                                onClick={() => navigate(`/workflows/${workflow.id}`)}
                            >
                                <div className="flex flex-row pl-[25px] pt-[15px]">
                                    <img
                                        alt="icon"
                                        src={`${import.meta.env.VITE_HOST_URL}${workflow.actionlogo}`}
                                        className="size-[26px]"
                                    />
                                    <img
                                        alt="icon"
                                        src={`${import.meta.env.VITE_HOST_URL}${workflow.reactionlogo}`}
                                        className="size-[26px] ml-[10px]"
                                    />
                                </div>
                                <div className="text-left w-full px-[25px]">
                                    <h2
                                        className="text-white text-[28px] font-[800] leading-[1.125] mt-6 m-0 p-0 break-words whitespace-normal"
                                    >
                                        {sliceText(workflow.name, 70)}
                                    </h2>

                                    <div className="absolute bottom-[20px]">
                                        <p
                                            className="text-white text-[15.7px] mt-[30px] max-w-[100%]"
                                        >
                                            by{' '}
                                            <span className="font-bold">
                                                {userName}
                                            </span>
                                        </p>

                                        <div
                                            style={{backgroundColor: "white"}}
                                            className="h-[38px] w-[160px] mt-[22px] rounded-full relative flex items-center justify-between px-1"
                                        >
                                            {workflow.isactivated ? (
                                                <p
                                                    className="text-black text-center flex-1 text-[14px] font-bold"
                                                >
                                                    Connected
                                                </p>
                                            ) : (
                                                <p
                                                    className="text-black text-center flex-1 text-[14px] font-bold ml-[50px]"
                                                >
                                                    Connect
                                                </p>
                                            )}

                                            <div
                                                style={{
                                                    backgroundColor: workflow.actioncolor,
                                                    transform: workflow.isactivated
                                                        ? "translateX(0)"
                                                        : "translateX(-120px)",
                                                }}
                                                className="h-[33px] w-[33px] rounded-full">
                                            </div>
                                        </div>
                                        {workflow.isactivated ? (
                                            <div className="flex flex-row mt-[22px]">
                                                <img
                                                    className="w-[14px] h-[18px]"
                                                    src="/user-icon.png"
                                                    alt="Avatar"
                                                />
                                                <span
                                                    className="text-white ml-2 text-[20px] font-bold -translate-y-1.5"
                                                >
                                                    1
                                                </span>
                                            </div>
                                        ) : (
                                            <div className="mb-[51px]"></div>
                                        )}
                                    </div>
                                </div>
                            </button>
                        ))}
                    </div>
                </div>
            </>
        )
    }

    return (
        <>
            <Navbar />

            {notification && (
                <MessageBox
                    message={notification}
                    trigger={notificationTrigger}
                    onClose={handleNotificationClose}
                    type="notification"
                    timeout={3000}
                />
            )}

            <MaxWidthWrapper
                maxWidth="1120px"
            >
                {
                    renderContent()
                }
            </MaxWidthWrapper>
        </>
    );
}

export default MyWorkflow