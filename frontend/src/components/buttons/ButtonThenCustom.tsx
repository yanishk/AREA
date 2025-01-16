import React from "react";

interface ButtonThenCustomProps {
    text: string;
    color: string;
    iconPath: string;
    onEdit: () => void;
    isEdit?: boolean;
    onDelete: () => void;
}

const ButtonThenCustom: React.FC<ButtonThenCustomProps> = ({ text, color, iconPath, onEdit, onDelete, isEdit }) => {
    const buttonColor = color || "#222222";

    return (
        <div className="relative">
            <button
                style={{ backgroundColor: buttonColor }}
                className="w-full h-[125px] flex justify-start items-center text-white rounded-lg font-[900]"
            >
                <p
                    className="ml-[35px] pr-5 text-[2.8rem] sm:text-[5.8rem]"
                >
                    Then
                </p>
                <img
                    alt="icon"
                    src={iconPath}
                    className="size-11"
                />
                <p
                    className="text-[1.4rem] font-bold ml-4">
                        {text}
                </p>
            </button>

            { isEdit && (
                <button
                    onClick={onEdit}
                    className="absolute top-[15px] right-[70px] text-[12px] font-extrabold text-white underline"
                >
                    Edit
                </button>
            )}

            <button
                onClick={onDelete}
                className="absolute top-[15px] right-[15px] text-[12px] text-sm font-extrabold text-white underline"
            >
                Delete
            </button>
        </div>
    );
};

export default ButtonThenCustom;
