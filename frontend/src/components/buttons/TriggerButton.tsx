import clsx from "clsx";
import React from "react";

interface TriggerButtonProps {
    onClick: () => void;
    text: string;
    textColor: string;
    bgColor: string;
    hoverColor: string;
}

const TriggerButton: React.FC<TriggerButtonProps> = ({ onClick, text, textColor, bgColor, hoverColor }) => {
    return (
        <button
            onClick={onClick}
            className={clsx(
                "w-full max-w-[420px] h-28 text-[2.2rem] rounded-full mt-5 font-[900]",
                textColor,
                bgColor,
                hoverColor
            )}
        >
            {text}
        </button>
    );
};

export default TriggerButton;
