import React, { useEffect, useState } from "react";

interface ErrorMessageProps {
    message?: string;
    trigger?: number;
}

const ErrorMessage: React.FC<ErrorMessageProps> = ({ message, trigger }) => {
    const [isVisible, setIsVisible] = useState(false);

    useEffect(() => {
        if (message) {
            setIsVisible(true);

            const timer = setTimeout(() => {
                setIsVisible(false);
            }, 3000);

            return () => clearTimeout(timer);
        }
    }, [message, trigger]);

    if (!message) return null;

    return (
        <div
            className={`fixed top-0 left-0 w-full p-5 bg-[#D0011B] text-white text-center font-bold text-[30px] z-50 transition-transform duration-500 ${
                isVisible ? "translate-y-0" : "-translate-y-full"
            }`}
        >
            {message}
        </div>
    );
};

export default ErrorMessage;
