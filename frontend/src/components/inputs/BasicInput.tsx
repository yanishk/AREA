import React from "react";

interface BasicInputProps {
    type: string;
    placeholder?: string;
    value: string;
    onChange?: (event: React.ChangeEvent<HTMLInputElement>) => void;
}

const BasicInput: React.FC<BasicInputProps> = ({ type, placeholder, value, onChange }) => (
    <input
        type={type}
        placeholder={placeholder}
        value={value}
        onChange={onChange}
        className={`text-black w-full px-4 h-[4.8rem] rounded-[15px] placeholder-[#C8C8C8] font-[900] text-[25px] leading-7`}
    />
);

export default BasicInput;
