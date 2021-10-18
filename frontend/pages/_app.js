import "tailwindcss/tailwind.css"
//import { Provider } from "next-auth/client"
//import { UserProvider } from "@auth0/nextjs-auth0"
import { AuthProvider } from "../context/auth"

function MyApp({ Component, pageProps }) {
  return (
    <>
      <AuthProvider>
        <Component {...pageProps} />
      </AuthProvider>
    </>
  )
}

export default MyApp
