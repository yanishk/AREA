import React from "react";
import clsx from "clsx";

interface CancelButtonProps {
    onClick: () => void;
    text: string;
    textColor: string;
    borderColor: string;
    bgColor: string;
}

const CancelButton: React.FC<CancelButtonProps> = ({ onClick, text, textColor, borderColor, bgColor }) => {
    return (
        <button
            onClick={onClick}
            className={clsx(
                "w-full border-4 text-[18px] rounded-full font-[900] py-3",
                textColor,
                borderColor,
                bgColor
            )}
        >
            {text}
        </button>
    );
};

export default CancelButton;
