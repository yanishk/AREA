import React from "react";
import { fieldType } from "../../utils/FieldType";

interface Field {
    id: string;
    label: string;
    type: fieldType;
    options?: string[];
}

interface DynamicFieldsProps {
    fields: Field[];
    values: Record<string, any>;
    onChange: (id: string, value: any) => void;
}

const DynamicFields: React.FC<DynamicFieldsProps> = ({ fields, values, onChange }) => {
    return (
        <div className="w-full max-w-[520px] mx-auto">
            {fields.map((field) => (
                <div key={field.id} className="mb-6">
                    <label className="text-white block text-xl font-bold mb-2">{field.label}</label>
                    {field.type === fieldType.SELECT ? (
                        <select
                            value={values[field.id] || ""}
                            onChange={(e) => onChange(field.id, e.target.value)}
                            className="w-full text-black h-[4.8rem] px-4 border-[3px] rounded-[10px] border-[#EEEEEE] placeholder-[#C8C8C8] font-[900] text-[1.6rem] leading-7"
                        >
                            {field.options?.map((option) => (
                                <option key={option} value={option}>
                                    {option}
                                </option>
                            ))}
                        </select>
                    ) : (
                        <input
                            type={field.type === fieldType.INPUT ? "input" : ""}
                            value={values[field.id] || ""}
                            onChange={(e) => onChange(field.id, e.target.value)}
                            className="w-full text-black h-[4.8rem] px-4 border-[3px] rounded-[10px] border-[#EEEEEE] focus:outline-none focus:border-black placeholder-[#C8C8C8] font-[900] text-[1.6rem] leading-7"
                        />
                    )}
                </div>
            ))}
        </div>
    );
};

export default DynamicFields;
