import { useParams } from "react-router-dom";
import { useEffect, useState } from "react";
import axios from "axios";
import Navbar from "../components/navbar/Navbar";

const WorkflowDetail = () => {
    const { workflowId } = useParams();
    const [workflow, setWorkflow] = useState<any | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const getWorkflowDetails = async () => {
            try {
                const result = await axios.get(`${import.meta.env.VITE_API_URL}workflows`, {
                    withCredentials: true,
                });

                const rightWorkflow = result.data.workflows.find((element: any) => element.id === workflowId);

                setWorkflow(rightWorkflow);
            } catch (error) {
                console.error("Error fetching workflow details:", error);
            }
            setLoading(false);
        };

        if (workflowId) {
            getWorkflowDetails();
        }
    }, [workflowId]);

        if (loading) {
            return;
        }

    return (
        <>
            <div>
                <Navbar />
                <div>
                <h1>{workflow.name}</h1>
                <p>{workflow.id}</p>
                <p>{workflow.actionid}</p>
                </div>
            </div>
        </>
    );
};

export default WorkflowDetail;
