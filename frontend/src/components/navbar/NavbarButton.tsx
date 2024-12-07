import React from "react";

interface NavbarButtonProps {
    onClick: () => void;
    text: string;
}

const NavbarButton: React.FC<NavbarButtonProps> = ({ onClick, text }) => {
    return (
        <button
            onClick={onClick}
            className="bg-[#222222] text-white font-bold px-10 tracking-wider py-4 rounded-full hover:bg-[#333333]"
        >
                {text}
        </button>
    );
};

export default NavbarButton;
