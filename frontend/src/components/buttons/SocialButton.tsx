import React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { IconDefinition } from "@fortawesome/fontawesome-svg-core";

interface SocialButtonProps {
    icon: IconDefinition;
    bgColor: string;
    hoverColor: string;
    text: string;
    iconPosition?: string;
    iconSize?: string;
    onClick?: () => void;
}

const SocialButton: React.FC<SocialButtonProps> = ({
    icon,
    bgColor,
    hoverColor,
    text,
    iconPosition = "left-6",
    iconSize = "size-8",
    onClick,
}) => {
    return (
        <button
            onClick={onClick}
            className={`w-full h-[75px] ${bgColor} text-white text-[1.6rem] rounded-full hover:${hoverColor} mt-4 font-[900] flex items-center justify-center relative`}
        >
            <FontAwesomeIcon
                icon={icon}
                className={`absolute ${iconPosition} ${iconSize}`}
            />
            <span
                className="flex justify-center items-center w-full ml-4">
                {text}
            </span>
        </button>
    );
};

export default SocialButton;
