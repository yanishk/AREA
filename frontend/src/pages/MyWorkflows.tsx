import axios from "axios";
import MaxWidthWrapper from "../components/MaxWidthWrapper"
import Navbar from "../components/navbar/Navbar"
import { useEffect, useState } from "react";
import { getUserName } from "../utils/GetUserName";
import { useNavigate } from "react-router-dom";

const MyWorkflow = () => {

    const [userWorkflows, setUserWorkflows] = useState<any[]>([]);
    const navigate = useNavigate();

    const getUserWorkflows = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}workflows`, {
                    withCredentials: true,
            });

            if (!result) {
                console.error("Error fetching services");
                return;
            }

            const workflowsProperties = result.data.workflows.map((workflow: any) => ({
                name: workflow.name,
                id: workflow.id,
            }));

            setUserWorkflows(workflowsProperties);
        } catch (err) {
            console.error("Error fetching services:", err);
        }
    };

    const [userName, setUserName] = useState("You");


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
                                // style={{ backgroundColor: "#333333"}}
                                className="h-[420px] flex flex-col items-start pt-2 cursor-pointer rounded-[10px] mr-6 mb-5 bg-[#333333] hover:bg-[#474747]"
                                onClick={() => navigate(`/workflow/${workflow.id}`)}
                            >
                                <div className="text-left w-full px-[30px]">
                                    <h2
                                        className="text-white text-[28px] font-[800] leading-[1.125] mt-6 m-0 p-0"
                                    >
                                        {workflow.name}
                                    </h2>
                                    <p
                                        className="text-white text-[15.7px] mt-[35px] mb-[110px]"
                                    >
                                        by{' '}
                                        <span className="font-bold">
                                            {userName}
                                        </span>
                                    </p>
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