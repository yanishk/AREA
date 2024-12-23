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
import MaxWidthWrapper from "../components/MaxWidthWrapper";

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

    const handleChangePassword = () => {
        navigate("/settings/change_password")
    }

    return (
        <>
            {isModalVisible && (
                <Modal
                    message={"Are you sure you want to delete your account ?"}
                    onClickYes={handleDelete}
                    onClickNo={handleCloseModal}
                />
            )}

            <Navbar />

            <div
                className="flex items-center justify-between px-10 py-5"
            >
                <div>
                    <CancelButton
                        text="Cancel"
                        textColor="text-[#222222]"
                        borderColor="border-[#222222]"
                        bgColor="bg-[#F9FAFB]"
                        onClick={goBack}
                    />
                </div>
            </div>

            <MaxWidthWrapper
                maxWidth="720px"
            >
                <div
                    className="flex items-center justify-center mt-10"
                >
                    <div
                        className="flex flex-col items-center text-center w-full"
                    >
                        <h1
                            className="text-5xl text-[#222222] font-[900]"
                        >
                            Account settings
                        </h1>

                        <div
                            className="flex items-center w-full mt-6"
                        >
                            <div
                                className="flex-grow border-[1.5px] border-[#e7e7e7] w-full"
                            />
                        </div>

                    </div>
                </div>

                <div
                    className="flex flex-col justify-center items-center"
                >
                    <div
                        className="w-full mt-10"
                    >
                        <p
                            className="font-extrabold text-[22.5px] mb-[6px]"
                        >
                            Email
                        </p>
                        <InputField
                            type="text"
                            value={userEmail}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                    </div>

                    <div
                        className="flex flex-col justify-start mt-4 w-full"
                    >
                        <p
                            className="font-extrabold text-[22.5px] mt-[6px]"
                        >
                            Password
                        </p>
                        <input
                            className="text-[1.6em] mt-2 font-mono"
                            type="password"
                            name="password"
                            value="****************"
                            disabled={true}
                        />
                        <div
                            className="flex justify-start mt-2"
                        >
                            <p
                                className="font-[700] text-[20px] text-[#0099FF] hover:text-[#1BA3FF] hover:cursor-pointer"
                                onClick={handleChangePassword}
                            >
                                Change password
                            </p>
                        </div>

                        <div
                            className="flex justify-start mt-6"
                        >
                            <p
                                className="font-[700] text-[20px] text-[#D0011C] hover:text-[#D51932] hover:cursor-pointer"
                                onClick={handleShowModal}
                            >
                                Delete my account
                            </p>
                        </div>
                    </div>


                        {/* <InputField
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
                        /> */}
                        {/* <ButtonValidation
                            text="Submit"
                            onClick={handleSubmit}
                        />
                        <ButtonValidation
                            text="Logout"
                            onClick={handleLogout}
                        /> */}
                </div>
            </MaxWidthWrapper>
        </>
    )
}

export default AccountPage;
