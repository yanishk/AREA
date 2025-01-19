import axios from "axios";
import { useEffect, useState } from "react";
import InputField from "../components/inputs/InputsField";
import { useNavigate } from "react-router-dom";
import { useAuth } from "./AuthContext";
import ActionButton from "../components/buttons/ActionButton";
import Modal from "../components/modals/Modal";
import Navbar from "../components/navbar/Navbar";
import MaxWidthWrapper from "../components/MaxWidthWrapper";

const AccountPage = () => {
    const [userEmail, setUserEmail] = useState("");
    const [userConnectionType, setUserConnectionType] = useState("");
    const [isModalVisible, setIsModalVisible] = useState(false);
    const navigate = useNavigate();
    const { logout } = useAuth();

    const getUserEmail = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}user`, {
                withCredentials: true,
            });
            setUserEmail(result.data.user.email);
            setUserConnectionType(result.data.user.connectiontype);
        } catch (error) {
            console.error("Failed to get email : ", error);
        }
    };

    useEffect(() => {
        getUserEmail();
    }, []);

    const handleDelete = async () => {
        try {
            const result = await axios.delete(
                `${import.meta.env.VITE_API_URL}user`, {
                    withCredentials: true,
                }
            );

            if (result.data.success) {
                navigate("/explore", {
                    state: {
                        notification: "Account deleted",
                        notificationTrigger: Date.now()
                    }
                });
            }

            logout();
        } catch (error) {
            console.error("Failed to delete account : ", error);
        }
    }

    const goBack = () => {
        navigate(-1);
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
                className="flex items-center justify-between px-6 lg:py-5 pb-14"
            >
                <div>
                    <ActionButton
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
                            className="text-4xl sm:text-5xl text-[#222222] font-[900]"
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
                    className="flex flex-col justify-center items-center lg:px-0 px-10"
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
                        />
                    </div>

                    <div
                        className="flex flex-col justify-start mt-4 w-full"
                    >
                        {userConnectionType === "basic" && (
                            <>
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
                                    <button
                                        className="font-[700] text-[20px] text-[#0099FF] hover:text-[#1BA3FF] hover:cursor-pointer"
                                        onClick={handleChangePassword}
                                    >
                                        Change password
                                    </button>
                                </div>
                            </>
                        )}

                        <div
                            className="flex justify-start mt-6"
                        >
                            <button
                                className="font-[700] text-[20px] text-[#D0011C] hover:text-[#D51932] hover:cursor-pointer"
                                onClick={handleShowModal}
                            >
                                Delete my account
                            </button>
                        </div>
                    </div>

                </div>
            </MaxWidthWrapper>
        </>
    )
}

export default AccountPage;
