interface ErrorFieldProps {
    error: string;
}

const ErrorField: React.FC<ErrorFieldProps> = ({ error }) => {
    return (
        <input
            value={error}
            className="w-full px-4 h-[4rem] mb-[22px] rounded-[12px] font-[800] text-[17.5px] text-white leading-7 bg-[#D0001B]"
            disabled={true}
        />
    );
};

export default ErrorField;
