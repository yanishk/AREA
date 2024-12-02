import React from "react";

interface ErrorMessageProps {
    message?: string;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({ message }) => {
    if (!message) return null;

    return (
        <div className="flex justify-center text-center">
            <p
                className="text-white bg-[#d0011b] fixed top-0 left-0 z-[9999] w-full p-5 text-[1.5em] font-[700] text-lg"
                role="alert"
                aria-live="assertive"
            >
                {message}
            </p>
        </div>
    );
};

export default ErrorMessage;
