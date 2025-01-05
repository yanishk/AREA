import { faGoogle, faSpotify, faDiscord, faGithub} from '@fortawesome/free-brands-svg-icons';
import { useState } from "react";
import axios from "axios";
import ErrorMessage from "../components/errors/ErrorMessage";
import InputField from "../components/inputs/InputsField";
import ButtonValidation from "../components/buttons/ButtonValidation";
import SocialButton from "../components/buttons/SocialButton";
import CallToActionLink from "../components/CallToActionLink";
import { useNavigate } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faArrowLeft } from '@fortawesome/free-solid-svg-icons';
import { handleOAuthLogin } from '../OAuth';

const RegisterForm = () => {

    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const [errorTrigger, setErrorTrigger] = useState(0);

    const navigate = useNavigate();

    const goToLogin = () => {
        navigate("/login");
    }

    const goBack = () => {
        navigate(-1);
    }

    const sendFormOnClick = async () => {
        if (!email || !password) {
            setEmail("");
            setPassword("");
            setError("Please fill email and password");
            return;
        }

        try {
            const result = await axios.post(`${import.meta.env.VITE_API_URL}register`, {
                email, password
            });

            if (result.data.success) {
                goToLogin();
            } else {
                setEmail("");
                setPassword("");
                setError("The email and password don't match");
                setErrorTrigger(Date.now());
            }
        } catch (error) {
            setEmail("");
            setPassword("");
            setError("The email and password don't match");
            setErrorTrigger(Date.now());
        }
    }

    const handleErrorClose = () => {
        setError("");
        setErrorTrigger(0);
    }

    return (
        <>
            <ErrorMessage message={error} trigger={errorTrigger} onClose={handleErrorClose}/>
            <div className="flex flex-col items-center justify-center min-h-screen bg-[#F9FAFB]">
                <h1
                    className="text-4xl text-[#222222] hover:text-[#333333] font-extrabold mt-10">
                    <a
                        href="/"
                    >
                        AREA
                    </a>
                </h1>
                <div className="relative flex w-full max-w-[480px] flex-col rounded-lg bg-white shadow-sm p-6 mt-12 mb-24">

                    <FontAwesomeIcon icon={faArrowLeft} className="size-8 cursor-pointer" onClick={goBack}></FontAwesomeIcon>
                    <div className="relative mt-5 items-center flex flex-col justify-center">
                        <h1
                            className=" text-[2.7rem] font-[900] text-[#222222]">
                            Sign up
                        </h1>
                    </div>

                    <div className="mt-6">
                        <InputField
                            type="email"
                            placeholder="Email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                        />
                        <InputField
                            type="password"
                            placeholder="Password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />

                        <ButtonValidation onClick={sendFormOnClick} text="Get started" />
                    </div>

                    <div className="flex items-center mt-6">
                        <div className="flex-grow border-t border-[#999999]" />
                            <span
                                className="mx-4 text-[#999999] font-[600]">
                                or
                            </span>
                        <div className="flex-grow border-t border-[#999999]" />
                    </div>

                    <SocialButton
                        icon={faGoogle}
                        bgColor="bg-[#E84A4A]"
                        hoverColor="hover:bg-[#F55B5B]"
                        text="Continue with Google"
                        onClick={() => handleOAuthLogin("Google")}
                    />
                    <SocialButton
                        icon={faSpotify}
                        bgColor="bg-[#20ca5d]"
                        hoverColor="hover:bg-[#20d561]"
                        text="Continue with Spotify"
                        onClick={() => handleOAuthLogin("Spotify")}
                    />
                    <SocialButton
                        icon={faDiscord}
                        bgColor="bg-[#7189DA]"
                        hoverColor="hover:bg-[#7a90da]"
                        text="Continue with Discord"
                        iconPosition="left-[27px]"
                        onClick={() => handleOAuthLogin("Discord")}
                    />
                    <SocialButton
                        icon={faGithub}
                        bgColor="bg-[#222222]"
                        hoverColor="hover:bg-[#333333]"
                        text="Continue with Github"
                        iconPosition="left-[20px]"
                        iconSize="size-11"
                        onClick={() => handleOAuthLogin("Github")}
                    />

                    <div
                        className="flex justify-center mt-8">
                        <CallToActionLink
                            href="/login"
                            mainText="Already on AREA?"
                            underlinedText="Log in here."
                        />
                    </div>

                </div>
            </div>
        </>
    )
}

export default RegisterForm