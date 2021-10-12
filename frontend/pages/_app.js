import 'tailwindcss/tailwind.css'
//import { AuthProvider } from '../context/auth'
import { Provider } from "next-auth/client"

function MyApp({ Component, pageProps: {session, ...pageProps} }) {
  return (
    <>
      <Provider session={session}>
        <Component {...pageProps} />
      </Provider>
    </>
  )
}

export default MyApp
