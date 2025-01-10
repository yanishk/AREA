import React from "react";

interface SelectFieldProps {
    label: string;
    value: string | undefined;
    options: string[];
    onChange: (value: string) => void;
}

const SelectField: React.FC<SelectFieldProps> = ({ label, value, options, onChange }) => {
    return (
        <div className="mb-6">
            <label className="text-white block text-xl font-bold mb-2">{label}</label>
            <select
                value={value || ""}
                onChange={(e) => onChange(e.target.value)}
                className="w-full text-black h-[4.8rem] px-4 border-[3px] rounded-[10px] border-[#EEEEEE] placeholder-[#C8C8C8] font-[900] text-[1.6rem] leading-7"
            >
                {options.map((option) => (
                    <option key={option} value={option}>
                        {option}
                    </option>
                ))}
            </select>
        </div>
    );
};

export default SelectField;
