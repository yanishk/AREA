import { useNavigate, useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import axios from "axios";
import Navbar from "../components/navbar/Navbar";
import getServiceDetails from "../utils/GetServiceDetails";
import CancelButton from "../components/buttons/CancelButton";
import MaxWidthWrapper from "../components/MaxWidthWrapper";
import { getUserName } from "../utils/GetUserName";

const WorkflowDetail = () => {
    const { workflowId } = useParams();
    const [workflow, setWorkflow] = useState<any | null>(null);
    const [loading, setLoading] = useState(true);
    const [username, setUsername] = useState("You");
    const navigate = useNavigate();

    const getWorkflowDetails = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}workflows`, {
                withCredentials: true,
            });

            const rightWorkflow = result.data.workflows.find((element: any) => element.id === workflowId);

            const actionDetails = await getServiceDetails(rightWorkflow.actionid, "action");
            const reactionDetails = await getServiceDetails(rightWorkflow.reactionid, "reaction");

            const newWorkflow = {
                ...rightWorkflow,
                actioncolor: actionDetails?.color,
                actionlogo: actionDetails?.logo,
                reactionlogo: reactionDetails?.logo,
            };

            setWorkflow(newWorkflow)
        } catch (error) {
            console.error("Error fetching workflow details:", error);
        }
        setLoading(false);
    };

    const goBack = () => {
        navigate(-1);
    }

    const handleDelete = async () => {
        try {
            const result = await axios.delete(`${import.meta.env.VITE_API_URL}workflows/${workflowId}`, {
                withCredentials: true,
            });

            if (result.data.success) {
                navigate("/");
            }

        } catch (error) {
            console.error("Error deleting workflow:", error);
        }
    }

    useEffect(() => {
        const fetchUserName = async () => {
            try {
                const username = await getUserName();
                setUsername(username);
            } catch (error) {
                console.error("Failed to get username: ", error);
            }
        }
        getWorkflowDetails();
        fetchUserName();
    }, []);

    if (loading) {
        return;
    }

    return (
        <>
            <div
                style={{ backgroundColor: workflow.actioncolor}}
            >
                <Navbar
                    isWhiteMod={true}
                />
                <div className="flex items-center justify-between px-10 py-5">
                    <div>
                        <CancelButton
                            text="Back"
                            textColor="text-white"
                            borderColor="border-white"
                            bgColor={`bg-[${workflow?.actioncolor}]`}
                            onClick={goBack}
                        />
                    </div>
                </div>
                <MaxWidthWrapper
                    maxWidth="450px"
                >
                    <div
                        className="flex flex-col text-left justify-center"
                    >
                        <div className="flex flex-row pt-[15px]">
                            <img
                                src={`${import.meta.env.VITE_HOST_URL}${workflow.actionlogo}`}
                                className="size-[60px]"
                            />
                            <img
                                src={`${import.meta.env.VITE_HOST_URL}${workflow.reactionlogo}`}
                                className="size-[60px] ml-[10px]"
                            />
                        </div>
                        <h1
                            className="text-white text-[35px] font-[800] mt-2"
                        >
                            {workflow.name}
                        </h1>
                    </div>

                    <p
                        className="text-white text-[15.7px] mt-[35px] pb-24"
                    >
                        by{' '}
                        <span className="font-bold">
                            {username}
                        </span>
                    </p>
                </MaxWidthWrapper>
            </div>

            <div className="flex justify-center mt-[20px]">
                <button
                    className="font-[700] text-[20px] text-[#BC5555] hover:opacity-85 hover:cursor-pointer"
                    onClick={handleDelete}
                >
                    Delete workflow
                </button>
            </div>

        </>
    );
};

export default WorkflowDetail;
