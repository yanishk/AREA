import React from "react";

interface NavbarButtonProps {
    onClick: () => void;
    text: string;
}

const NavbarButton: React.FC<NavbarButtonProps> = ({ onClick, text }) => {
    return (
        <button
            onClick={onClick}
            className="bg-[#000000] text-white px-4 py-2 rounded-full hover:bg-[#222222]"
        >
                {text}
        </button>
    );
};

export default NavbarButton;
