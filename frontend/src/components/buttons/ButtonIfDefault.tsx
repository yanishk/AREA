import React from "react";

interface ButtonIfDefaultProps {
    onClick: () => void;
    text: string;
}

const ButtonIfDefault: React.FC<ButtonIfDefaultProps> = ({ text, onClick }) => {
    return (
        <button
            onClick={onClick}
            className="px-64 bg-[#222222] flex justify-center items-center text-white text-[5.2rem] rounded-lg font-[900]"
        >
            {text}
        </button>
    );
};

export default ButtonIfDefault;
