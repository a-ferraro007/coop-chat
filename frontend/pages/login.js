import React, { useState } from "react";
import { signIn } from "next-auth/client";
import axios from "axios";
import { useAuth } from "../context/auth";

function login() {
  return (
    <div>
      {" "}
      <p>
        <a
          href="/api/auth/signin"
          onClick={(e) => {
            e.preventDefault();
            signIn();
          }}
        >
          You must be signed in to view this page
        </a>
      </p>
    </div>
  );
}

export default login;

//const [isCreate, setUseIsCreate] = useState(false)
//const [username, setUsername] = useState('')
//const [email, setEmail] = useState('')
//const [password, setPassword] = useState('')
//const [confirmPassword, setConfirmPassword] = useState('')
//const { login, user, create } = useAuth()

//const submitClick = async (e) => {
//  e.preventDefault()
//  if (isCreate) {
//    create(email, username, password)
//  } else {
//    login(username, password)
//  }
//}

//const createClick = (e) => {
//  e.preventDefault()
//  setEmail('')
//  setUsername('')
//  setConfirmPassword('')
//  setPassword('')
//  setUseIsCreate(!isCreate)
//}

//return (
//  <>
//    <form
//      className="min-h-400 w-665 mx-auto bg-white rounded-60 mt-10 box-shadow  px-28 py-14"
//      onSubmit={submitClick}
//    >
//      <h1 className="text-primary-dark text-3xl text-center font-sans font-bold">
//        {' '}
//        Voice Chat{' '}
//      </h1>

//      {isCreate ? (
//        <>
//          <input
//            className="w-full rounded-60 box-shadow h-14 mb-2 mt-6 pl-4 text-primary-dark"
//            type="email"
//            placeholder="email"
//            value={email}
//            onChange={(e) => {
//              setEmail(e.target.value)
//            }}
//            required={true}
//          />
//          <input
//            className="w-full rounded-60 box-shadow h-14 mb-2 pl-4 text-primary-dark"
//            type="text"
//            placeholder="username"
//            onChange={(e) => {
//              setUsername(e.target.value)
//            }}
//            value={username}
//            required={true}
//          />
//          <input
//            className="w-full rounded-60 box-shadow h-14 pl-4 text-primary-dark mb-2"
//            type="password"
//            placeholder="password"
//            onChange={(e) => {
//              setPassword(e.target.value)
//            }}
//            value={password}
//            required={true}
//          />
//          <input
//            className="w-full rounded-60 box-shadow h-14 pl-4 text-primary-dark"
//            type="password"
//            placeholder="confirm password"
//            onChange={(e) => {
//              setConfirmPassword(e.target.value)
//            }}
//            value={confirmPassword}
//            required={true}
//          />
//        </>
//      ) : (
//        <>
//          <input
//            className="w-full h-14 mb-4 mt-10 pl-4 text-primary-dark border-b-2 border-primary-dark"
//            type="text"
//            placeholder="username"
//            onChange={(e) => {
//              setUsername(e.target.value)
//            }}
//            value={username}
//            required={true}
//          />
//          <input
//            className="w-full box-shadow h-14 pl-4 text-primary-dark mb-8 border-b-2 border-primary-dark"
//            type="password"
//            placeholder="password"
//            onChange={(e) => {
//              setPassword(e.target.value)
//            }}
//            value={password}
//            required={true}
//          />
//        </>
//      )}
//      <button
//        className="border-2 border-primary-dark rounded-60 px-4 py-3 w-full mb-2 bg-transparent text-
//      primary-dark"
//      >
//        {!isCreate ? 'login' : 'create account'}
//      </button>
//      <button
//        className=" text-primary-dark underline text-sm center w-full"
//        onClick={createClick}
//      >
//        {!isCreate ? 'create account' : 'login'}
//      </button>
//    </form>
//  </>
//)
