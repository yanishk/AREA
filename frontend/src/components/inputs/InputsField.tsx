import React from "react";

interface InputFieldProps {
    type: string;
    placeholder?: string;
    value: string;
    onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void;
    disableFocus?: boolean,
}

const InputField: React.FC<InputFieldProps> = ({ type, placeholder, value, onChange, disableFocus }) => (
    <input
        type={type}
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        className={`text-black w-full px-4 h-[4.8rem] mb-[22px] border-[7px] rounded-[15px] border-[#EEEEEE] placeholder-[#C8C8C8] font-[900] text-[25px] leading-7 ${
            disableFocus
                ? ""
                : "focus:outline-none focus:border-black"
        }`}
    />
);

export default InputField;
