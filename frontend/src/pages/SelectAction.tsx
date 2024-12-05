import { useNavigate } from "react-router-dom";

const SelectActionPage = () => {
  const navigate = useNavigate();

  const handleActionClick = (action: { id: string; name: string; color: string }) => {
    navigate("/create", { state: { selectedAction: action } });
  };

  return (
    <div>
      <h1>Select an Action</h1>
      <button onClick={() => handleActionClick({ id: "1", name: "Every day at", color: "black" })}>
        Every day at
      </button>
    </div>
  );
};

export default SelectActionPage