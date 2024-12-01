import React from "react";

interface NavbarLinkProps {
    href: string;
    text: string;
}

const NavbarLink: React.FC<NavbarLinkProps> = ({ href, text }) => {
    return (
        <a
            href={href}
            className="text-[#222222] hover:text-[#444444] font-bold text-[20px] ml-3"
        >
            {text}
        </a>
    );
};

export default NavbarLink;
