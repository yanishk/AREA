import React from "react";

interface TriggerButtonProps {
    onClick: () => void;
    text: string;
}

const TriggerButton: React.FC<TriggerButtonProps> = ({ onClick, text }) => {
    return (
        <button
            onClick={onClick}
            className="w-full max-w-[420px] h-28 bg-white text-black text-[2.2rem] rounded-full hover:bg-[#EEEEEE] mt-5 font-[900]"
        >
            {text}
        </button>
    );
};

export default TriggerButton;
