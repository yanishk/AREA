import React from "react";

interface ButtonThenCustomProps {
    text: string;
    color: string;
    onEdit: () => void;
    onDelete: () => void;
}

const ButtonThenCustom: React.FC<ButtonThenCustomProps> = ({ text, color, onEdit, onDelete }) => {
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
                onClick={onEdit}
                className="absolute top-[10px] right-[70px] text-sm font-extrabold text-white underline"
            >
                Edit
            </button>

            <button
                onClick={onDelete}
                className="absolute top-[10px] right-[10px] text-sm font-extrabold text-white underline"
            >
                Delete
            </button>
        </div>
    );
};

export default ButtonThenCustom;
