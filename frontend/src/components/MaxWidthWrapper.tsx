import React from "react";

interface MaxWidthWrapperProps {
    children: React.ReactNode;
    maxWidth?: string;
}

const MaxWidthWrapper: React.FC<MaxWidthWrapperProps> = ({ children, maxWidth = "780px" }) => {
    return (
        <div
            className={`w-full mx-auto max-w-[${maxWidth}]`}
        >
            {children}
        </div>
    );
};

export default MaxWidthWrapper;