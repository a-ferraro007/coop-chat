import { createContext, useState, useEffect, useContext } from 'react'
import Router from 'next/router'
import axios from 'axios'
export const AuthContext = createContext()

export const AuthProvider = (props) => {
  const [user, setUser] = useState(null)
  const [isLoading, setIsLoading] = useState(false)

  useEffect(() => {
    console.log(user)
  }, [user])

  const createAccount = async (email, username, password) => {
    try {
      const body = {
        email: email,
        username: username,
        password: password
      }
      const res = await axios.post('http://localhost:8000/create_account', body)
      setUser(res.data)
      Router.push('/')
    } catch (error) {
      console.log(error)
    }
  }

  const login = async (username, password) => {
    try {
      const body = {
        username: username,
        password: password
      }
      const res = await axios.post('https://chat.tony.place/login', body)
      setUser(res.data)
      Router.push('/')
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
        create: createAccount,
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
