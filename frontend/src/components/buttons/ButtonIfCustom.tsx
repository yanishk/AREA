import React from "react";

interface ButtonIfCustomProps {
    text: string;
}

const ButtonIfCustom: React.FC<ButtonIfCustomProps> = ({text }) => {
    return (
        <button
            className="w-full h-28 bg-[#222222] text-white text-[2.2rem] rounded-full hover:bg-[#333333] mt-5 font-[900]"
        >
            {text}
        </button>
    );
};

export default ButtonIfCustom;
