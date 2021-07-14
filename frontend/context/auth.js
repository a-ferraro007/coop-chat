//import router
//create authcontext
//create authprovider component
//

import { createContext, useState, useEffect, useContext } from 'react'
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

export const useAuth = () => useContext(AuthContext)

export const ProtectedRoute = (props) => {
  const auth = useAuth()

  if (auth.loading || auth.isAuthenticated === false) {
    return (
      <>
        <div>...Loading </div>
      </>
    )
  } else {
    Router.push('/login')
  }

  return props.children
}
