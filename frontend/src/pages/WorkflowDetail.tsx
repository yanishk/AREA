import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import axios from "axios";
import Navbar from "../components/navbar/Navbar";
import getServiceDetails from "../utils/GetServiceDetails";

const WorkflowDetail = () => {
    const { workflowId } = useParams();
    const [workflow, setWorkflow] = useState<any | null>(null);
    const [loading, setLoading] = useState(true);

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

    useEffect(() => {
        getWorkflowDetails();
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
                <div>
                <h1>{workflow.name}</h1>
                <p>{workflow.id}</p>
                <p>{workflow.actionid}</p>
                <p>{workflow.actioncolor}</p>
                </div>
            </div>
        </>
    );
};

export default WorkflowDetail;
