import axios from "axios";
import MaxWidthWrapper from "../components/MaxWidthWrapper"
import Navbar from "../components/navbar/Navbar"
import { useEffect, useState } from "react";
import { getUserName } from "../utils/GetUserName";
import { useNavigate } from "react-router-dom";
import getServiceDetails from "../utils/GetServiceDetails";

const MyWorkflow = () => {

    const [userWorkflows, setUserWorkflows] = useState<any[]>([]);
    const [userName, setUserName] = useState("You");
    const navigate = useNavigate();

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
            console.log("ici : ", workflowsDetails);
        } catch (err) {
            console.error("Error fetching workflows:", err);
        }
    };

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

    return (
        <>
            <Navbar />

            <MaxWidthWrapper
                maxWidth="1120px"
            >
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
                <div className="mt-16 ml-[35px]">
                    <div className="grid sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-3">
                        {userWorkflows?.map((workflow) => (
                            <button
                                key={workflow.id}
                                style={{
                                    backgroundColor: workflow.actioncolor,
                                    boxShadow: '0 1.5px 6px 0 rgba(0,0,0,.24)'
                                }}
                                className="h-[420px] flex flex-col items-start pt-2 cursor-pointer rounded-[10px] mr-6 mb-5 hover:opacity-90"
                                onClick={() => navigate(`/workflows/${workflow.id}`)}
                            >
                                <div className="flex flex-row pl-[25px] pt-[15px]">
                                    <img
                                        src={`${import.meta.env.VITE_HOST_URL}${workflow.actionlogo}`}
                                        className="size-[26px]"
                                    />
                                    <img
                                        src={`${import.meta.env.VITE_HOST_URL}${workflow.reactionlogo}`}
                                        className="size-[26px] ml-[10px]"
                                    />
                                </div>
                                <div className="text-left w-full px-[25px]">
                                    <h2
                                        className="text-white text-[28px] font-[800] leading-[1.125] mt-6 m-0 p-0 break-words whitespace-normal"
                                    >
                                        {workflow.name}
                                    </h2>
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
                                        <p
                                            className="text-black text-center flex-1 text-[14px] font-bold"
                                        >
                                            Connected
                                        </p>
                                        <div
                                            style={{backgroundColor: workflow.actioncolor}}
                                            className="h-[33px] w-[33px] rounded-full">
                                        </div>
                                    </div>
                                </div>

                            </button>

                        ))}
                    </div>
                </div>
            </MaxWidthWrapper>
        </>
    )
}

export default MyWorkflow