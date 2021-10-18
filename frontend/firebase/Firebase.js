import firebase from "firebase/compat/app"
import {
  getAuth,
  signInWithEmailAndPassword,
  signOut,
  createUserWithEmailAndPassword,
} from "firebase/auth"

export let app
export let auth
export let signIn
export let createUser
export let signOutUser

if (typeof window !== "undefined" && !firebase.apps.length) {
  const FirebaseCredentials = {
    apiKey: process.env.NEXT_PUBLIC_FIREBASE_PUBLIC_API_KEY,
    authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
    projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
  }

  app = firebase.initializeApp(FirebaseCredentials)
  auth = getAuth(app)
  signIn = signInWithEmailAndPassword
  createUser = createUserWithEmailAndPassword
  signOutUser = signOut
}
