import React from "react";

interface ButtonThenDefaultProps {
    onClick: () => void;
    text: string;
}

const ButtonThenDefault: React.FC<ButtonThenDefaultProps> = ({ text, onClick }) => {
    return (
        <button
            onClick={onClick}
            className="px-64 bg-[#999999] flex justify-center items-center text-white text-[5.2rem] rounded-lg font-[900]"
        >
            {text}
        </button>
    );
};

export default ButtonThenDefault;
