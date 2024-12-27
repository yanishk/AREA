import React from "react";

interface MaxWidthWrapperProps {
    children: React.ReactNode;
    maxWidth: string;
}

const MaxWidthWrapper: React.FC<MaxWidthWrapperProps> = ({ children, maxWidth}) => {
    return (
        <div
            style={{ maxWidth }}
            className="w-full mx-auto"
        >
            {children}
        </div>
    );
};

export default MaxWidthWrapper;
