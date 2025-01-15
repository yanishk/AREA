import React from "react";

interface NavbarButtonProps {
    onClick: () => void;
    text: string;
    isWhiteMode?: boolean;
}

const NavbarButton: React.FC<NavbarButtonProps> = ({ onClick, text, isWhiteMode }) => {
    return (
        <>
            {isWhiteMode ? (
                <button
                    onClick={onClick}
                    className="bg-white text-[#222222] font-[900] px-10 tracking-wider py-4 rounded-full hover:opacity-90 ml-2"
                >
                        {text}
                </button>
            ) : (
                <button
                    onClick={onClick}
                    className="bg-[#222222] text-white font-bold px-10 tracking-wider py-4 rounded-full hover:bg-[#333333] ml-2"
                >
                        {text}
                </button>
            )}
        </>
    );
};

export default NavbarButton;
