import React from 'react';
import {BrowserRouter as Router, Link, Route} from "react-router-dom";
import PrivateRoute from './PrivateRoute';
import Home from './pages/Home';
import Admin from './pages/Admin';
import {AuthProvider} from "./context/AuthContext";
import ReactNotification from 'react-notifications-component'
import 'react-notifications-component/dist/theme.css'

import Login from './pages/Home';

function App(props) {
  const auth = { token: null, loggedIn: false }
  const channel = {id : null}
  return (
    <AuthProvider value={auth}>
      <ReactNotification />
      <Router>
            <div>
              <Route exact path="/" component={Home}/>
              <Route path="/home" component={Home} />
              <PrivateRoute path="/admin" component={Admin} />
            </div>
      </Router>
    </AuthProvider>
  );
}
export default App;
