import React from "react";

interface CancelButtonProps {
    onClick: () => void;
    text: string;
}

const CancelButton: React.FC<CancelButtonProps> = ({ onClick, text }) => {
    return (
        <>
            <button
                onClick={onClick}
                className="w-full bg-[#F9FAFB] text-[#222222] border-4 border-[#222222] text-[18px] rounded-full font-[900] py-3"
            >
                {text}
            </button>
        </>
    );
};

export default CancelButton;
