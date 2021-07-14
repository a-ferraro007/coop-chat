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



const ProtectedRoute = Component => {
  const Wrapper = props => {
    const auth = useContext(AuthContext)
    return <Component {...props} />

  }

  return Wrapper
}

export default ProtectedRoute