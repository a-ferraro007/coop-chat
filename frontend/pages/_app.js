import 'tailwindcss/tailwind.css'
import { AuthProvider, ProtectedRoute } from '../context/auth'

//import ProtectedRoute from '../protectedRoute'

function MyApp({ Component, pageProps }) {
  return (
    <>
      <AuthProvider>
        <ProtectedRoute>
          <Component {...pageProps} />
        </ProtectedRoute>
      </AuthProvider>
    </>
  )
}

export default MyApp
