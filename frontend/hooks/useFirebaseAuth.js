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
  }

  const clear = () => {
    setUser(null)
    setLoading(false)
    return true
  }

  const signInWithEmailAndPassword = (email, password) =>
    signIn(auth, email, password)

  const createUserWithEmailAndPassword = (email, password) =>
    createUser(auth, email, password)

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
