import { useNavigate, useParams } from "react-router-dom";
import { useEffect, useRef, useState } from "react";
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
    const [isWorkflowActivated, setIsWorkflowActivated] = useState(false);
    const [buttonText, setButtonText] = useState("");
    const [isEditing, setIsEditing] = useState(false);
    const [workflowName, setWorkflowName] = useState("");
    const inputRef = useRef<HTMLTextAreaElement | null>(null);
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
                actionname: actionDetails?.name,
                reactionlogo: reactionDetails?.logo,
            };

            setWorkflow(newWorkflow)
            setIsWorkflowActivated(newWorkflow.isactivated);
            setWorkflowName(newWorkflow.name);
            setButtonText(newWorkflow.isactivated ? "Connected" : "Connect");
            console.log("new workflow ici : ", newWorkflow);
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

    const handleActivation = async () => {
        if (!workflow) return;

        const updatedWorkflow = {
            ...workflow,
            isactivated: !isWorkflowActivated,
        };

        try {
            await axios.put(
                `${import.meta.env.VITE_API_URL}workflows/${workflowId}`,
                updatedWorkflow,
                {
                    withCredentials: true
                }
            );

            setIsWorkflowActivated(!isWorkflowActivated);
            setButtonText(!isWorkflowActivated ? "Connected" : "Connect");
        } catch (error) {
            console.error("Error updating workflow activation:", error);
        }
    };

    const handleSaveName = async () => {
        if (!workflow) return;

        if (!workflowName.trim()) {
            return;
        }
        console.log("update name ici ", workflowName);

        const updatedWorkflow = {
            ...workflow,
            name: workflowName,
        };


        try {
            await axios.put(
                `${import.meta.env.VITE_API_URL}workflows/${workflowId}`,
                updatedWorkflow,
                {
                    withCredentials: true,
                }
            );

            setWorkflow(updatedWorkflow);
            setIsEditing(false);
        } catch (error) {
            console.error("Error updating workflow name:", error);
        }
    };

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

    const handleClickOutside = (event: MouseEvent) => {
        if (inputRef.current && !inputRef.current.contains(event.target as Node)) {
            handleSaveName();
        }
    };

    useEffect(() => {
        if (isEditing) {
            document.addEventListener("mousedown", handleClickOutside);
        } else {
            document.removeEventListener("mousedown", handleClickOutside);
        }
        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, [isEditing]);

    if (loading) {
        return;
    }

    return (
        <>
            <div
                style={{backgroundColor: workflow.actioncolor}}
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
                        {isEditing ? (
                            <>
                                <textarea
                                    ref={inputRef}
                                    value={workflowName}
                                    onChange={(e) => setWorkflowName(e.target.value)}
                                    onBlur={handleSaveName}
                                    className="text-white text-[36px] h-[340px] w-full font-[800] mt-[10px] bg-transparent border-4 p-4 rounded-lg border-white outline-none break-words whitespace-normal"
                                    style={{overflow: "hidden"}}
                                />
                                <div
                                    className="flex justify-end text-center"
                                >
                                    <p
                                        className="text-white font-bold text-[16px] underline cursor-pointer"
                                        onClick={() => setIsEditing(true)}
                                    >
                                        Save
                                    </p>
                                </div>
                            </>
                        ) : (
                            <>
                                <h1
                                    className="text-white text-[36px] font-[800] mt-2 break-words whitespace-normal cursor-pointer"
                                    onClick={() => setIsEditing(true)}
                                >
                                    {workflowName}
                                </h1>
                                <div
                                    className="flex justify-end text-center"
                                >
                                    <p
                                        className="text-white font-bold text-[16px] underline cursor-pointer"
                                        onClick={() => setIsEditing(true)}
                                    >
                                        Edit title
                                    </p>
                                </div>
                        </>
                        )}
                    </div>

                    <p
                        className="text-white text-[15.7px] mt-[35px] pb-[45px]"
                    >
                        by{' '}
                        <span className="font-bold">
                            {username}
                        </span>
                    </p>
                    <div className="flex flex-row">
                        <img
                            className="w-[14px] h-[18px]"
                            src="/user-icon.png"
                            alt="Avatar"
                        />
                        <span
                            className="text-white ml-2 text-[20px] font-bold -translate-y-1.5 mb-[10px]"
                        >
                            1
                        </span>
                    </div>
                </MaxWidthWrapper>
            </div>

            <MaxWidthWrapper
                maxWidth="450px"
            >
                <div
                    style={{
                        backgroundColor: "#222222",
                        boxShadow: '0 1.5px 6px 0 rgba(0,0,0,.24)'
                    }}
                    className="h-[100px] w-full mt-[55px] rounded-full relative flex items-center justify-between cursor-pointer"
                    onClick={handleActivation}
                >
                    {isWorkflowActivated ? (
                        <>
                            <p
                                className="text-white text-center flex-1 text-[29px] font-bold ml-[40px]"
                            >
                                {buttonText}
                            </p>
                        </>
                    ) : (
                        <>
                            <p
                                className="text-white text-center flex-1 text-[29px] font-bold ml-[150px]"
                            >
                                {buttonText}
                            </p>
                        </>
                    )}

                    <div
                        style={{
                            backgroundColor: workflow.actioncolor,
                            transform: isWorkflowActivated
                            ? "translateX(0)"
                            : "translateX(-348px)",
                        }}
                        className="h-[90px] w-[90px] rounded-full mr-[6px] transition-transform duration-300 ease-in-out">
                    </div>
                </div>
            </MaxWidthWrapper>

            <div className="flex justify-center mt-[55px] mb-[45px]">
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
