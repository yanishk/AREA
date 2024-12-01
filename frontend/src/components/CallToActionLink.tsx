import React from "react";

interface RegisterLinkProps {
    href: string;
    mainText?: string;
    underlinedText: string;
}

const CallToActionLink: React.FC<RegisterLinkProps> = ({ href, mainText, underlinedText }) => {
    return (
        <a
            href={href}
            className="text-[19px] mb-6 font-[900] text-black"
        >
            {mainText}
            <span className="underline ml-1">{underlinedText}</span>
        </a>
    );
};

export default CallToActionLink;
