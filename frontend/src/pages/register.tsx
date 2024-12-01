import {
  faGoogle,
  faFacebookF,
  faApple,
  faXTwitter,
} from "@fortawesome/free-brands-svg-icons";
import { useState } from "react";
import axios from "axios";
import ErrorMessage from "../components/ErrorMessage";
import InputField from "../components/inputs/InputsField";
import ButtonValidation from "../components/buttons/ButtonValidation";
import SocialButton from "../components/buttons/SocialButton";
import CallToActionLink from "../components/CallToActionLink";

const RegisterForm = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");

  const sendFormOnClick = async () => {
    if (!email || !password) {
      setEmail("");
      setPassword("");
      setError("Please fill email and password");
      return;
    }

    try {
      const result = await axios.post(
        `${import.meta.env.VITE_API_URL}register`,
        {
          email,
          password,
        },
      );

      if (result.data.success) {
        alert("Successful register");
      } else {
        setEmail("");
        setPassword("");
        setError("The email and password don't match");
      }
    } catch (error) {
      setEmail("");
      setPassword("");
      setError("The email and password don't match");
    }
  };

  const handleGoogleRegister = () => {
    return;
  };

  const handleFacebookRegister = () => {
    return;
  };

  const handleTwitterRegister = () => {
    return;
  };

  const handleAppleRegister = () => {
    return;
  };

  return (
    <>
      <ErrorMessage message={error} />
      <div className="flex items-center justify-center min-h-screen bg-[#F9FAFB]">
        <div className="relative flex w-full max-w-[480px] flex-col rounded-lg bg-white shadow-sm p-6 mt-24 mb-24">
          <div className="relative mt-5 items-center flex flex-col justify-center">
            <h1 className=" text-[2.7rem] font-[900] text-black">Sign up</h1>
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
            <span className="mx-4 text-[#999999] font-[600]">or</span>
            <div className="flex-grow border-t border-[#999999]" />
          </div>

          <SocialButton
            icon={faGoogle}
            bgColor="bg-[#E84A4A]"
            hoverColor="bg-[#F55B5B]"
            text="Continue with Google"
            onClick={handleGoogleRegister}
          />
          <SocialButton
            icon={faFacebookF}
            bgColor="bg-[#3B579D]"
            hoverColor="bg-[#4C6AAE]"
            text="Continue with Facebook"
            onClick={handleFacebookRegister}
          />
          <SocialButton
            icon={faXTwitter}
            bgColor="bg-[#A9A9A9]"
            hoverColor="bg-[#B8B8B8]"
            text="Continue with X"
            iconPosition="left-[27px]"
            onClick={handleTwitterRegister}
          />
          <SocialButton
            icon={faApple}
            bgColor="bg-[#222222]"
            hoverColor="bg-[#333333]"
            text="Continue with Apple"
            iconPosition="left-[20px]"
            iconSize="size-11"
            onClick={handleAppleRegister}
          />

          <div className="flex justify-center mt-8">
            <CallToActionLink
              href="/login"
              mainText="Already on AREA?"
              underlinedText="Log in here."
            />
          </div>
        </div>
      </div>
    </>
  );
};

export default RegisterForm;
