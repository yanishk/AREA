import React from "react";

interface ButtonThenDefaultProps {
    onClick: () => void;
    text: string;
    isDisabled?: boolean;
}

const ButtonThenDefault: React.FC<ButtonThenDefaultProps> = ({ text, onClick, isDisabled }) => {

    return (
        <button
        onClick={isDisabled ? undefined : onClick}
        className={`w-full h-[125px] flex items-center text-white text-[2.8rem] lg:text-[5.8rem] md:text-[4.8rem] rounded-lg font-[900] pr-6 ${
                isDisabled
                    ? "cursor-not-allowed bg-[#999999] justify-center"
                    : "cursor-pointer bg-[#222222]"
            }`}
        >
            <span className="pl-5">
                {text}
            </span>
            {!isDisabled &&
            (
                <button
                    onClick={!isDisabled ? onClick : undefined}
                    className="h-[60px] bg-[#F9FAFB] px-10 text-black text-[20px] rounded-[50px] font-bold ml-auto hover:bg-[#EEEEEE]"
                >
                    Add
                </button>
                )}
        </button>
    );
};

export default ButtonThenDefault;

