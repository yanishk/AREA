import React from "react";

interface ButtonIfDefaultProps {
    onClick: () => void;
    text: string;
}

const ButtonIfDefault: React.FC<ButtonIfDefaultProps> = ({ text, onClick }) => {

    const handleKeyDown = (event: React.KeyboardEvent<HTMLDivElement>) => {
        if (event.key === "Enter" || event.key === " ") {
            onClick();
        }
    };

    return (
        <div
            onClick={onClick}
            onKeyDown={handleKeyDown}
            role="button"
            tabIndex={0}
            className="w-[730px] h-[125px] bg-[#222222] flex items-center text-white text-[5.8rem] rounded-lg font-[900] pr-6 cursor-pointer relative"
        >
            <span className="absolute left-1/2 transform -translate-x-1/2">
                {text}
            </span>
            <button
                onClick={onClick}
                className="h-[60px] bg-[#F9FAFB] px-10 text-black text-[20px] rounded-[50px] font-bold ml-auto hover:bg-[#EEEEEE] cursor-pointer"
            >
                Add
            </button>
        </div>
    );
};

export default ButtonIfDefault;


