import React from "react";

interface NavbarLinkProps {
    href: string;
    text: string;
}

const NavbarLink: React.FC<NavbarLinkProps> = ({ href, text }) => {
    return (
        <a
            href={href}
            className="text-gray-700 hover:text-black"
        >
            {text}
        </a>
    );
};

export default NavbarLink;
