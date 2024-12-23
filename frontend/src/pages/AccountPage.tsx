import axios from "axios";
import { useEffect, useState } from "react";
import InputField from "../components/inputs/InputsField";
import ButtonValidation from "../components/buttons/ButtonValidation";
import TriggerButton from "../components/buttons/TriggerButton";
import { useNavigate } from "react-router-dom";
import { useAuth } from "./AuthContext";
import CancelButton from "../components/buttons/CancelButton";
import Modal from "../components/modals/Modal";
import Navbar from "../components/navbar/Navbar";

const AccountPage = () => {
    const [userEmail, setUserEmail] = useState("");
    const [password, setPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [isModalVisible, setIsModalVisible] = useState(false);
    const navigate = useNavigate();
    const { logout } = useAuth();

    const getUserEmail = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}user/email`, {
                withCredentials: true,
            });
            setUserEmail(result.data.email);
        } catch (error) {
            console.error("Failed to get email : ", error);
        }
    };

    useEffect(() => {
        getUserEmail();
    }, []);

    const handleSubmit = async () => {
        if (!password || !newPassword)
            return;
        try {
            const result = await axios.put(
                `${import.meta.env.VITE_API_URL}user/modify-password`,
                {
                    oldpassword: password,
                    password: newPassword,
                },
                {
                    withCredentials: true,
                }
            );

            if (result.data.success) {
                goHome();
            }
        } catch (error) {
            console.error("Failed to modify password : ", error);
        }
    }

    const handleLogout = () => {
        logout();
        goHome();
    }

    const handleDelete = async () => {
        try {
            await axios.delete(
                `${import.meta.env.VITE_API_URL}user/delete-account`, {
                    withCredentials: true,
                }
            );
            handleLogout();
        } catch (error) {
            console.error("Failed to delete account : ", error);
        }
    }

    const goBack = () => {
        navigate(-1);
    };

    const goHome = () => {
        navigate("/");
    };

    const handleShowModal = () => {
        setIsModalVisible(true);
    };

    const handleCloseModal = () => {
        setIsModalVisible(false);
    };

    return (
        <>
            {isModalVisible && (
                <Modal
                    message={"Do you want to delete you'r account ?"}
                    onClickYes={handleDelete}
                    onClickNo={handleCloseModal}
                />
            )}
            <Navbar />
            <div className="flex items-center justify-between px-10 py-5">
                <div>
                    <CancelButton
                        text="Cancel"
                        textColor="text-[#222222]"
                        borderColor="border-[#222222]"
                        bgColor="bg-[#F9FAFB]"
                        onClick={goHome}
                    />
                </div>
            </div>
            <div className="flex flex-col justify-center items-center">
                <p className="text-black text-xl font-[800]">
                    {userEmail}
                </p>
                <div className="max-w-96 mt-10">
                    <InputField
                        type="password"
                        placeholder="Current password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                    />
                    <InputField
                        type="password"
                        placeholder="New password"
                        value={newPassword}
                        onChange={(e) => setNewPassword(e.target.value)}
                    />
                    <ButtonValidation
                        text="Submit"
                        onClick={handleSubmit}
                    />
                    <ButtonValidation
                        text="Logout"
                        onClick={handleLogout}
                        />
                    <TriggerButton
                        text="Delete my account"
                        textColor="text-[#D0001D]"
                        bgColor="bg-[#F9FAFB]"
                        hoverColor="hover:bg-[#D0001D] hover:text-white"
                        onClick={handleShowModal}
                        />
                </div>
            </div>
        </>
    )
}

export default AccountPage;
