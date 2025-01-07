import clsx from "clsx";
import React from "react";

interface NavbarLinkProps {
    href: string;
    text: string;
    textColor: string;
    hoverColor: string;
}

const NavbarLink: React.FC<NavbarLinkProps> = ({ href, text, textColor, hoverColor }) => {
    return (
        <a
            href={href}
            // className="text-[#222222] hover:text-[#444444] font-bold text-[20px] ml-3"
            className={clsx(
                "font-bold text-[20px] ml-3",
                `${textColor}`,
                `${hoverColor}`
            )}
        >
            {text}
        </a>
    );
};

export default NavbarLink;
