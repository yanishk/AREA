import React from "react";

interface NavbarButtonProps {
    onClick: () => void;
    text: string;
    isWhiteMod?: boolean;
}

const NavbarButton: React.FC<NavbarButtonProps> = ({ onClick, text, isWhiteMod }) => {
    return (
        <>
            {isWhiteMod ? (
                <button
                    onClick={onClick}
                    className="bg-white text-[#222222] font-[900] px-10 tracking-wider py-4 rounded-full hover:opacity-90"
                >
                        {text}
                </button>
            ) : (
                <button
                    onClick={onClick}
                    className="bg-[#222222] text-white font-bold px-10 tracking-wider py-4 rounded-full hover:bg-[#333333]"
                >
                        {text}
                </button>
            )}
        </>
    );
};

export default NavbarButton;
