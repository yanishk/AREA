import React from "react";

interface ButtonThenCustomProps {
    text: string;
    color: string;
    onDelete: () => void;
}

const ButtonThenCustom: React.FC<ButtonThenCustomProps> = ({ text, color, onDelete }) => {
    const buttonColor = color || "#222222";

    return (
        <div className="relative inline-block">
            <button
                style={{ backgroundColor: buttonColor }}
                className="px-64 flex justify-center items-center text-white text-[5.2rem] rounded-lg font-[900]"
            >
                {text}
            </button>

            <button
                onClick={onDelete}
                className="absolute top-[-10px] right-[-10px] text-sm text-red-500 underline hover:text-red-700"
            >
                Delete
            </button>
        </div>
    );
};

export default ButtonThenCustom;
