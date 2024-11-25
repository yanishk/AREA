import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faGoogle, faFacebookF, faApple, faXTwitter } from '@fortawesome/free-brands-svg-icons';
import React, { useState } from "react";

export const RegisterForm = () => {

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const sendFormOnClick = async () => {

    try {
      const result = await fetch('http://localhost:8080/login', {
        method: 'POST',
        headers: {'Content-Type': 'application/json',},
        body: JSON.stringify({
          email: email,
          password: password,
        })
      })

      if (result.ok) {
        alert("Successful login");
      } else {
        alert("Login failed");
      }
    } catch (error) {
      console.log(error)
      alert("ERROR");
    }
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-[#F9FAFB]">
      <div className="relative flex w-full max-w-[480px] flex-col rounded-lg bg-white shadow-sm p-6 mt-24 mb-24">

        <div className="relative mt-5 items-center flex flex-col justify-center">
          <h1
            className=" text-[2.7rem] font-[900] text-black">
              Sign up
          </h1>
        </div>

        <div className="mt-6">
          <input
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full px-4 h-[4.8rem] mb-3 border-[7px] rounded-[15px] border-[#EEEEEE] focus:outline-none focus:border-black placeholder-[#C8C8C8] font-[900] text-[25px] leading-7"
          />
          <input
            type="password"
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full px-4 h-[4.8rem] mt-2.5 border-[7px] rounded-[15px] border-[#EEEEEE] focus:outline-none focus:border-black placeholder-[#C8C8C8] font-[900] text-[25px] leading-7"
          />

          <button
            onClick={sendFormOnClick}
            className="w-full h-28 bg-[#222222] text-white text-[2.2rem] rounded-full hover:bg-[#333333] mt-8 font-[900]">
              Get started
          </button>
        </div>

        <div className="flex items-center mt-6">
          <div className="flex-grow border-t border-[#999999]" />
            <span
              className="mx-4 text-[#999999] font-[600]">
                or
            </span>
          <div className="flex-grow border-t border-[#999999]" />
        </div>

        <button
          className="w-full h-[75px] bg-[#E84A4A] text-white text-[1.6rem] rounded-full hover:bg-[#F55B5B] mt-6 font-[900] flex items-center justify-center relative">
          <FontAwesomeIcon
            icon={faGoogle}
            className="absolute left-6 size-8"/>
          <span className="flex justify-center items-center w-full ml-4">Continue with Google</span>
        </button>
        <button
          className="w-full h-[75px] bg-[#3B579D] text-white text-[1.6rem] rounded-full hover:bg-[#4C6AAE] mt-4 font-[900] flex items-center justify-center relative">
          <FontAwesomeIcon
            icon={faFacebookF}
            className="absolute left-6 size-8"/>
          <span className="flex justify-center items-center w-full ml-4">Continue with Facebook</span>
        </button>
        <button
          className="w-full h-[75px] bg-[#A9A9A9] text-white text-[1.6rem] rounded-full hover:bg-[#B8B8B8] mt-4 font-[900] flex items-center justify-center relative">
          <FontAwesomeIcon
            icon={faXTwitter}
            className="absolute left-[27px] size-8"/>
          <span className="flex justify-center items-center w-full">Continue with X</span>
        </button>
        <button
          className="w-full h-[75px] bg-[#222222] text-white text-[1.6rem] rounded-full hover:bg-[#333333] mt-4 font-[900] flex items-center justify-center relative">
          <FontAwesomeIcon
            icon={faApple}
            className="absolute left-[20px] size-11"/>
          <span className="flex justify-center items-center w-full mt-1 ml-4">Continue with Apple</span>
        </button>

        <div
          className="flex justify-center mt-6">
          <a
            href="#"
            className="text-[19px] mt-3 mb-6 font-[900] text-black ">
              Already on AREA?
              <span
                className="underline">
                  Log in here.
              </span>
          </a>
        </div>


      </div>
    </div>
  );
};

export default RegisterForm;
