import axios from "axios";
import { useEffect, useState } from "react";
import Navbar from "../components/navbar/Navbar";
import MaxWidthWrapper from "../components/MaxWidthWrapper";
import InputField from "../components/inputs/InputsField";
import ButtonValidation from "../components/buttons/ButtonValidation";
import { useNavigate } from "react-router-dom";
import ErrorField from "../components/errors/ErrorField";


const ChangePasswordPage = () => {
    const [userEmail, setUserEmail] = useState("");
    const [password, setPassword] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");
    const [error, setError] = useState("");
    const [passwordError, setPasswordError] = useState("");
    const [newPasswordError, setNewPasswordError] = useState("");
    const [confirmPasswordError, setConfirmPasswordError] = useState("");
    const navigate = useNavigate();

    const getUserEmail = async () => {
        try {
            const result = await axios.get(`${import.meta.env.VITE_API_URL}user`, {
                withCredentials: true,
            });
            setUserEmail(result.data.user.email);
        } catch (error) {
            console.error("Failed to get email : ", error);
        }
    };

    useEffect(() => {
        getUserEmail();
    }, []);

    const goHome = () => {
        navigate("/explore");
    };

    const handleSubmit = async () => {
        if (!password) {
            setPasswordError("Please provide your current password");
            setNewPasswordError("");
            setConfirmPasswordError("");
            setError("");
            return;
        }

        if (!newPassword) {
            setPasswordError("");
            setNewPasswordError("Please provide your new password");
            setConfirmPasswordError("");
            setError("");
            return;
        }

        if (!confirmPassword) {
            setPasswordError("");
            setNewPasswordError("");
            setConfirmPasswordError("Please confirm your new password");
            setError("");
            return;
        }

        if (newPassword !== confirmPassword) {
            setPasswordError("");
            setNewPasswordError("");
            setConfirmPasswordError("");
            setError("The new password and the confirmation do not match");
            return;
        }

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

    const handleCancel = () => {
        navigate(-1);
    }

    return (
        <>
            <Navbar />
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
                            Change password
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
                        <p
                            className="font-extrabold text-[23.5px] text-[#909090]"
                        >
                            {userEmail}
                        </p>

                        <p
                            className="font-extrabold text-[22.5px] mb-[6px] mt-10"
                        >
                            Current password
                        </p>

                        { error &&
                            <ErrorField
                                error={error}
                            />
                        }

                        { passwordError &&
                            <ErrorField
                                error={passwordError}
                            />
                        }

                        <InputField
                            type="password"
                            placeholder=""
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />

                        <p
                            className="font-extrabold text-[22.5px] mb-[6px]"
                        >
                            New password
                        </p>

                        { newPasswordError &&
                            <ErrorField
                                error={newPasswordError}
                            />
                        }

                        <InputField
                            type="password"
                            placeholder=""
                            value={newPassword}
                            onChange={(e) => setNewPassword(e.target.value)}
                        />

                        <p
                            className="font-extrabold text-[22.5px] mb-[6px]"
                        >
                            Confirm new password
                        </p>

                        { confirmPasswordError &&
                            <ErrorField
                                error={confirmPasswordError}
                            />
                        }

                        <InputField
                            type="password"
                            placeholder=""
                            value={confirmPassword}
                            onChange={(e) => setConfirmPassword(e.target.value)}
                        />

                        <div className="flex flex-col text-center justify-center mb-28 max-w-[300px] mx-auto">
                            <ButtonValidation
                                text="Change"
                                onClick={handleSubmit}
                            />
                            <button
                                className="underline mt-7 text-[1.2em] text-[#222222] font-[900] cursor-pointer"
                                onClick={handleCancel}
                            >
                                Cancel
                            </button>
                        </div>
                    </div>
                </div>
            </MaxWidthWrapper>
        </>
    )
}

export default ChangePasswordPage;
