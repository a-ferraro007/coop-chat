import nookies from "nookies"
import { useState, useEffect } from "react"
import { auth, signIn, createUser, signOutUser } from "../firebase/Firebase"

export default function useFirebaseAuth() {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)

  const authStateChanged = async (authState) => {
    if (!authState) {
      setUser(null)
      setLoading(false)
      nookies.set(undefined, "token", "", { path: "/" })
      return
    }

    setLoading(true)
    const token = await authState.getIdToken()
    setUser(authState)
    nookies.set(undefined, "token", token, { path: "/" })
    setLoading(false)
    return
  }

  const clear = () => {
    setUser(null)
    setLoading(false)
    return true
  }

  const signInWithEmailAndPassword = (email, password) =>
    signIn(auth, email, password)

  const createUserWithEmailAndPassword = async (email, password) => {
    const fb = await createUser(auth, email, password)

    try {
      let headers = new Headers()
      headers.append("Authorization", `Bearer ${fb.user.accessToken}`)
      const body = JSON.stringify({
        uid: fb.user.uid,
      })
      const opts = {
        method: "POST",
        headers,
        body,
      }
      const res = await fetch("http://localhost:8000/create_user", opts)
      const data = await res.json()
    } catch (error) {
      console.log(error)
    }
    return fb
  }

  const signOut = () => signOutUser(auth).then(clear)

  useEffect(() => {
    return auth.onIdTokenChanged(authStateChanged)
  }, [])

  return {
    user,
    loading,
    signInWithEmailAndPassword,
    createUserWithEmailAndPassword,
    signOut,
  }
}
