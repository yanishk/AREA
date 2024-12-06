import { useLocation, useNavigate } from "react-router-dom";

const SelectReactionPage = () => {
    const navigate = useNavigate();
    const location = useLocation();

    const selectedAction = location.state?.selectedAction;

    const handleReactionClick = (reaction: { id: string; name: string; color: string }) => {
        navigate("/create", {
            state: {
                selectedAction,
                selectedReaction: reaction,
            }
        });
    };

    return (
        <div>
            <h1>Select a Reaction</h1>
            <button onClick={() =>
                handleReactionClick({
                    id: "2",
                    name: "send Mail",
                    color: "#FFE560"
                })
            }>
                Send mail
            </button>
        </div>
    );
};

export default SelectReactionPage