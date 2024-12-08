import React from "react";

interface ButtonIfDefaultProps {
    onClick: () => void;
    text: string;
}

const ButtonIfDefault: React.FC<ButtonIfDefaultProps> = ({ text, onClick }) => {
    return (
        <button
            onClick={onClick}
            className="w-[600px] bg-[#222222] flex justify-center items-center text-white text-[5.2rem] rounded-lg font-[900]"
        >
            {text}
        </button>
    );
};

export default ButtonIfDefault;
