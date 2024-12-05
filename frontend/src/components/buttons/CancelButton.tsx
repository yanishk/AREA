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
                className="w-full h-28 bg-[#F9FAFB] text-[#222222] border-4 border-[#222222] text-[2.2rem] rounded-full mt-5 font-[900]"
            >
                {text}
            </button>
        </>
    );
};

export default CancelButton;
