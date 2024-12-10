import React from "react";

interface ButtonThenDefaultProps {
    onClick: () => void;
    text: string;
    disabled?: boolean;
}

const ButtonThenDefault: React.FC<ButtonThenDefaultProps> = ({ text, onClick, disabled }) => {
    return (
        <button
            onClick={disabled ? undefined : onClick}
            className={`w-[600px] bg-[#999999] flex justify-center items-center text-white text-[5.2rem] rounded-lg font-[900] ${
                disabled ? "cursor-not-allowed opacity-50" : ""
            }`}
        >
            {text}
        </button>
    );
};

export default ButtonThenDefault;
