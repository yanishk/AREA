import React from "react";

interface ButtonThenDefaultProps {
    onClick: () => void;
    text: string;
    disabled?: boolean;
}

const ButtonThenDefault: React.FC<ButtonThenDefaultProps> = ({ text, onClick, disabled }) => {

    return (
        <button
        onClick={disabled ? undefined : onClick}
        className={`w-[730px] h-[125px] flex items-center text-white text-[5.8rem] rounded-lg font-[900] pr-6 ${
                disabled
                    ? "cursor-not-allowed bg-[#999999] justify-center"
                    : "cursor-pointer bg-[#222222]"
            }`}
        >
            <span className="ml-14">
                {text}
            </span>
            {!disabled &&
            (
                <button
                    onClick={!disabled ? onClick : undefined}
                    className="h-[60px] bg-[#F9FAFB] px-10 text-black text-[20px] rounded-[50px] font-bold ml-auto hover:bg-[#EEEEEE]"
                >
                    Add
                </button>
                )}
        </button>
    );
};

export default ButtonThenDefault;

