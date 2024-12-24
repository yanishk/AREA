import React from "react";

interface InputFieldProps {
    type: string;
    placeholder?: string;
    value: string;
    onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

const InputField: React.FC<InputFieldProps> = ({ type, placeholder, value, onChange }) => (
    <input
        type={type}
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        className="w-full px-4 h-[4.8rem] mb-[22px] border-[7px] rounded-[15px] border-[#EEEEEE] focus:outline-none focus:border-black placeholder-[#C8C8C8] font-[900] text-[25px] leading-7"
    />
);

export default InputField;
