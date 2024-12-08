import React from "react";

interface ButtonIfCustomProps {
    text: string;
    color: string;
    iconPath: string;
    onEdit: () => void;
    onDelete: () => void;
}

const ButtonIfCustom: React.FC<ButtonIfCustomProps> = ({ text, color, iconPath, onEdit, onDelete }) => {
    const buttonColor = color || "#222222";

    return (
        <div className="relative inline-block">
            <button
                style={{ backgroundColor: buttonColor }}
                className="w-[600px] flex justify-center items-center text-white rounded-lg font-[900]"
            >
                <p className="pr-10 text-[5.2rem]">If</p>
                <img src={iconPath} className="size-10" />
                <p className="text-[1.2rem] ml-4">{text}</p>
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

export default ButtonIfCustom;
