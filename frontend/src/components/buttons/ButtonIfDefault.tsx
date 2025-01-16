import React from "react";

interface ButtonIfDefaultProps {
    onClick: () => void;
    text: string;
}

const ButtonIfDefault: React.FC<ButtonIfDefaultProps> = ({ text, onClick }) => {

    return (
        <button
            onClick={onClick}
            className="w-full h-[125px] bg-[#222222] flex items-center text-white text-[2.8rem] sm:text-[5.8rem] rounded-lg font-[900] pr-6 cursor-pointer relative"
        >
            <span className="absolute sm:left-1/2 left transform sm:-translate-x-1/2 sm:pl-0 pl-5">
                {text}
            </span>
            <button
                onClick={onClick}
                className="h-[60px] bg-[#F9FAFB] px-10 text-black text-[20px] rounded-[50px] font-bold ml-auto hover:bg-[#EEEEEE] cursor-pointer"
            >
                Add
            </button>
        </button>
    );
};

export default ButtonIfDefault;


