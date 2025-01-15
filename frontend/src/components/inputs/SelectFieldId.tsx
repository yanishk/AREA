import React from 'react';

interface SelectFieldIdProps {
    label: string;
    value: string | undefined;
    paramMap: Record<string, string>;
    onChange: (id: string) => void;
}

const SelectFieldId: React.FC<SelectFieldIdProps> = ({ label, value, paramMap, onChange }) => {
    return (
        <div className="mb-6">
            <label className="text-white block text-xl font-bold mb-2">{label}</label>
            <select
                value={value || ""}
                onChange={(e) => onChange(e.target.value)}
                className="w-full text-black h-[4.8rem] px-4 border-[3px] rounded-[10px] border-[#EEEEEE] placeholder-[#C8C8C8] font-[900] text-[1.6rem] leading-7"
            >
                {Object.entries(paramMap).map(([name, id]) => (
                    <option key={id} value={id}>
                        {name}
                    </option>
                ))}
            </select>
        </div>
    );
};

export default SelectFieldId;
