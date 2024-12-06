import { useLocation, useNavigate } from "react-router-dom";

const SelectActionPage = () => {
    const navigate = useNavigate();
    const location = useLocation();

    const selectedReaction = location.state?.selectedReaction;

    const handleActionClick = (action: { id: string; name: string; color: string }) => {
        navigate("/create", {
            state: {
                selectedReaction,
                selectedAction: action,
            }
        });
    };

    return (
        <div>
            <h1>Select an Action</h1>
            <button onClick={() =>
                handleActionClick({
                    id: "1",
                    name: "Every day at",
                    color: "#ABABAB"
                })
            }>
                Every day at
            </button>
        </div>
    );
};

export default SelectActionPage