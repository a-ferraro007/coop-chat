import React, { useEffect, useState } from "react"
import { useAuth } from "../context/auth"
import { useRouter } from "next/router"

function login() {
  const [isCreate, setUseIsCreate] = useState(false)
  const [username, setUsername] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const {
    loading,
    user,
    signInWithEmailAndPassword,
    createUserWithEmailAndPassword,
  } = useAuth()
  const router = useRouter()

  useEffect(() => {
    if (!loading && user) {
      router.push("/")
    }
  }, [loading])

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (isCreate) {
      try {
        const res = await createUserWithEmailAndPassword(username, password)
        router.push("/")
      } catch (error) {
        console.log(error)
      }
    } else {
      try {
        const res = await signInWithEmailAndPassword(username, password)

        router.push("/")
      } catch (error) {
        console.log(error)
      }
    }
  }

  const handleIsCreate = async (e) => {
    e.preventDefault()
    setEmail("")
    setUsername("")
    setConfirmPassword("")
    setPassword("")
    setUseIsCreate(!isCreate)
  }

  return !loading && !user ? (
    <>
      <form
        className="min-h-400 w-665 mx-auto bg-white rounded-60 mt-10 box-shadow  px-28 py-14"
        onSubmit={handleSubmit}
      >
        <h1 className="text-primary-dark text-3xl text-center font-sans font-bold">
          {" "}
          Voice Chat{" "}
        </h1>

        {isCreate ? (
          <>
            <input
              className="w-full rounded-60 box-shadow h-14 mb-2 mt-6 pl-4 text-primary-dark"
              type="email"
              placeholder="email"
              value={email}
              onChange={(e) => {
                setEmail(e.target.value)
              }}
              required={true}
            />
            <input
              className="w-full rounded-60 box-shadow h-14 mb-2 pl-4 text-primary-dark"
              type="text"
              placeholder="username"
              onChange={(e) => {
                setUsername(e.target.value)
              }}
              value={username}
              required={true}
            />
            <input
              className="w-full rounded-60 box-shadow h-14 pl-4 text-primary-dark mb-2"
              type="password"
              placeholder="password"
              onChange={(e) => {
                setPassword(e.target.value)
              }}
              value={password}
              required={true}
            />
            <input
              className="w-full rounded-60 box-shadow h-14 pl-4 text-primary-dark"
              type="password"
              placeholder="confirm password"
              onChange={(e) => {
                setConfirmPassword(e.target.value)
              }}
              value={confirmPassword}
              required={true}
            />
          </>
        ) : (
          <>
            <input
              className="w-full h-14 mb-4 mt-10 pl-4 text-primary-dark border-b-2 border-primary-dark"
              type="text"
              placeholder="username"
              onChange={(e) => {
                setUsername(e.target.value)
              }}
              value={username}
              required={true}
            />
            <input
              className="w-full box-shadow h-14 pl-4 text-primary-dark mb-8 border-b-2 border-primary-dark"
              type="password"
              placeholder="password"
              onChange={(e) => {
                setPassword(e.target.value)
              }}
              value={password}
              required={true}
            />
          </>
        )}
        <button
          className="border-2 border-primary-dark rounded-60 px-4 py-3 w-full mb-2 bg-transparent text-
        primary-dark"
        >
          {!isCreate ? "login" : "create account"}
        </button>
        <button
          className=" text-primary-dark underline text-sm center w-full"
          onClick={handleIsCreate}
        >
          {!isCreate ? "create account" : "login"}
        </button>
      </form>
    </>
  ) : (
    <></>
  )
}

export default login
