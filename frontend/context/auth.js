import { createContext, useState, useEffect, useContext } from 'react'
import Router from 'next/router'
export const AuthContext = createContext()

export const AuthProvider = (props) => {
  const [user, setUser] = useState(null)
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    setUser(false)
  }, [])

  const login = async (email, passsword) => {
    try {
      const body = {
        username: username,
        password: password
      }
      const resp = await axios.post('http://localhost:8000/login', body)
      console.log(resp.data)
      setUser(resp.data)
    } catch (error) {
      console.log(error)
    }
  }

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated: !!user,
        user: user,
        login: login,
        loading: isLoading
      }}
    >
      {props.children}
    </AuthContext.Provider>
  )
}

export const useAuth = () => {
  return useContext(AuthContext)
}

const ProtectedRoute = (Component) => {
  const Wrapper = (props) => {
    const auth = useAuth()
    useEffect(() => {
      if (auth.isAuthenticated !== true) {
        Router.push('/login')
      }
    }, [])
    return <Component {...props} />
  }
  //Wrapper.getInitialProps = async (ctx) => {
  //  console.log('ctx')
  //  return {}
  //}
  return Wrapper
}

export default ProtectedRoute
