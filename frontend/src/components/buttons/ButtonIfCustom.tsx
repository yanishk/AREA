import React from "react";

interface ButtonIfCustomProps {
    text: string;
    color: string;
    iconPath: string;
    onEdit: () => void;
    onEditFalse?: boolean;
    onDelete: () => void;
}

const ButtonIfCustom: React.FC<ButtonIfCustomProps> = ({ text, color, iconPath, onEdit, onDelete, onEditFalse }) => {
    const buttonColor = color || "#222222";

    return (
        <div className="relative inline-block">
            <button
                style={{ backgroundColor: buttonColor }}
                className="w-[730px] h-[125px] flex justify-start items-center text-white rounded-lg font-[900]"
            >
                <p className="ml-[35px] pr-10 text-[5.8rem] text-left">If</p>
                <img alt="icon" src={iconPath} className="size-11" />
                <p className="text-[1.4rem] font-bold ml-4">{text}</p>
            </button>

            {!onEditFalse && (
                <button
                    onClick={onEdit}
                    className="absolute top-[15px] right-[70px] text-[12px] font-bold text-white underline"
                >
                    Edit
                </button>
            )}

            <button
                onClick={onDelete}
                className="absolute top-[15px] right-[15px] text-[12px] font-bold text-white underline"
            >
                Delete
            </button>
        </div>
    );
};

export default ButtonIfCustom;
