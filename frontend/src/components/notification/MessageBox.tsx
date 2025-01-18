import React, { useEffect, useState } from "react";

interface ErrorMessageProps {
    message: string;
    type: string;
    trigger: number;
    timeout: number;
    onClose: () => void;
}

const MessageBox: React.FC<ErrorMessageProps> = ({ message, trigger, onClose, timeout, type }) => {
    const [isVisible, setIsVisible] = useState(false);

    useEffect(() => {
        if (message) {
            setIsVisible(true);

            const timer = setTimeout(() => {
                setIsVisible(false);
                onClose?.();
            }, timeout);

            return () => clearTimeout(timer);
        }
    }, [message, trigger]);

    if (!message) return null;

    return type === "error" ? (
            <div
                className={`fixed top-0 left-0 w-full p-5 bg-[#D0011B] text-white text-center font-bold text-[30px] z-50 transition-transform duration-500 ${
                    isVisible ? "translate-y-0" : "-translate-y-full"
                }`}
            >
                {message}
            </div>
        ) : (
            <div
                className={`fixed top-0 left-0 w-full p-5 bg-[#2CBE60] text-white text-center font-bold text-[30px] z-50 transition-transform duration-500 ${
                    isVisible ? "translate-y-0" : "-translate-y-full"
                }`}
            >
                {message}
            </div>
        );
};

export default MessageBox;
