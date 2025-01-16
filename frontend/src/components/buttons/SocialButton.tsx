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
    iconPosition = "left-[20px]",
    iconSize = "size-11",
    onClick,
}) => {
    return (
        <button
            onClick={onClick}
            className={`w-full h-[75px] ${bgColor} text-white text-[1.4rem] sm:text-[1.6rem] rounded-full ${hoverColor} mt-4 font-[900] flex items-center justify-center relative`}
        >
            <FontAwesomeIcon
                icon={icon}
                className={`ml-5 ${iconPosition} ${iconSize}`}
            />
            <span
                className="flex justify-center items-center w-full pl-2 pr-2">
                {text}
            </span>
        </button>
    );
};

export default SocialButton;
