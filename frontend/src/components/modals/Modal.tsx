import TriggerButton from "../buttons/TriggerButton";

interface ModalProps {
    message: string;
    onClickYes: () => void;
    onClickNo: () => void;
}

const Modal: React.FC<ModalProps> = ({message, onClickYes, onClickNo}) => {
    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
            <div className="bg-[#F9FAFB] rounded-lg shadow p-6">
                <h3 className="mb-5 text-lg font-[900] text-black dark:text-gray-400">
                    {message}
                </h3>
                <TriggerButton
                    text="Yes"
                    textColor="text-[#D0001D]"
                    bgColor="bg-[#F9FAFB]"
                    hoverColor="hover:bg-[#D0001D] hover:text-white"
                    onClick={onClickYes}
                />
                <TriggerButton
                    text="No"
                    textColor="text-white"
                    bgColor="bg-[#222222]"
                    hoverColor="hover:bg-[#333333]"
                    onClick={onClickNo}
                />
            </div>
        </div>
    );
};

export default Modal;