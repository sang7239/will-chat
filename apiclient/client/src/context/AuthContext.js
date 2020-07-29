import React, {Component} from 'react'

const AuthContext = React.createContext()

class AuthProvider extends Component {
    // Context state
      state = {
        auth: {token: localStorage.getItem("token"), isLoggedIn: localStorage.getItem("token") != null}
      }
  
    // Method to update state
    setAuth = (auth) => {
      this.setState((prevState) => ({ auth }))
      localStorage.setItem("token", auth.token);
    }

    removeAuth = (auth) => {
      this.setState((prevState) => ({auth}))
      localStorage.removeItem("token", auth.token);
    }
    
    render() {
      const { children } = this.props
      const { auth } = this.state
      const { setAuth } = this
      const { removeAuth} = this
      return (
        <AuthContext.Provider
          value={{
            auth,
            setAuth,
            removeAuth,
          }}
        >
          {children}
        </AuthContext.Provider>
      )
    }
  }
  
  export default AuthContext
  
  export { AuthProvider }