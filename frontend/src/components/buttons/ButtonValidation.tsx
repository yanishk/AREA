import React from "react";

interface ButtonValidationProps {
    onClick: () => void;
    text: string;
}

const ButtonValidation: React.FC<ButtonValidationProps> = ({ onClick, text }) => {
    return (
        <button
            onClick={onClick}
            className="w-full h-28 bg-[#222222] text-white text-[2.2rem] rounded-full hover:bg-[#333333] mt-5 font-[900]"
        >
            {text}
        </button>
    );
};

export default ButtonValidation;
