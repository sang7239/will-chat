import React, {useState, Component} from "react";
import {Link, Redirect} from 'react-router-dom';
import axios from 'axios';
import logo from '../img/willchat.png'; //https://logomakr.com/6hIK2k
import AuthContext, { AuthConsumer } from "../context/AuthContext";
class Home extends Component {
    static contextType = AuthContext
    constructor(props) {
        super(props);
        this.state = {
            message: "",
            color: "",
            loginFailed: "",
        }
    }
    componentDidMount() {
        const auth = this.context;
    }
    handleSignUp = (event) => {
        event.preventDefault();
        axios.post("https://localhost:4000/v1/users",  {
            email: this.email.value,
            password: this.password.value,
            passwordConf: this.passwordConf.value,
            userName: this.userName.value,
            firstName: this.firstName.value,
            lastName: this.lastName.value
        }).then(result => {
            this.email.value = "";
            this.password.value = "";
            this.passwordConf.value = "";
            this.userName.value = "";
            this.firstName.value = "";
            this.lastName.value = "";
            if (result.status == 200) {
                this.setState({message: "User Created! Please Login", color: "green"});
            }
        }).catch(e => {
            this.setState({message: e.response.data, color: "red"});
            console.log(e);
        })
    }
    handleLogIn = (event, auth, setAuth) => {
        event.preventDefault();
        axios.post("https://localhost:4000/v1/sessions",  {
            email: this.loginemail.value,
            password: this.loginpass.value
        }).then(result => {
            if (result.status == 200) {
                const newauth = {token: result.headers.authorization, loggedIn: true};
                setAuth(newauth);
            }
        })
        .catch(e => {
            this.setState({loginFailed: e.response.data, color: "red"});
            console.log(e);
        })
    }
    clickClose = () => {
        this.inputElement.click();
      }
    render(){
        const { auth, setAuth } = this.context
        if (auth.loggedIn) {
            this.clickClose();
            return <Redirect to="/admin" />
        }
        return (
            <div>                
                <div class="signup-form">
                <img class="col-xs-12" src={logo} alt="logo"></img>
                <form onSubmit={this.handleSignUp} class="form-horizontal">
                    <div class="col-xs-8 col-xs-offset-4">
                        <h2>Sign Up</h2>
                    </div>		
                    <div class="form-group">
                        <label class="control-label col-xs-4">First Name</label>
                        <div class="col-xs-8">
                            <input type="text" class="form-control" name="firstname" required="required" 
                            ref={node => (this.firstName = node)}/>
                        </div>        	
                    </div>
                    <div class="form-group">
                        <label class="control-label col-xs-4">Last Name</label>
                        <div class="col-xs-8">
                            <input type="text" class="form-control" name="lastname" required="required" 
                            ref={node => (this.lastName = node)}/>
                        </div>        	
                    </div>
                    <div class="form-group">
                        <label class="control-label col-xs-4">Username</label>
                        <div class="col-xs-8">
                            <input type="text" class="form-control" name="username" required="required" 
                            ref={node => (this.userName = node)}/>
                        </div>        	
                    </div>
                    <div class="form-group">
                        <label class="control-label col-xs-4">Email Address</label>
                        <div class="col-xs-8">
                            <input type="email" class="form-control" name="email" required="required" 
                            ref={node => (this.email = node)}/>
                        </div>        	
                    </div>
                    <div class="form-group">
                        <label class="control-label col-xs-4">Password</label>
                        <div class="col-xs-8">
                            <input type="password" class="form-control" name="password" required="required" 
                            ref={node => (this.password = node)}/>
                        </div>        	
                    </div>
                    <div class="form-group">
                        <label class="control-label col-xs-4">Confirm Password</label>
                        <div class="col-xs-8">
                            <input type="password" class="form-control" name="confirm_password" required="required" 
                            ref={node => (this.passwordConf = node)}/>
                        </div>        	
                    </div>
                    <div id ="signup-response" class="text-center" ref={node => (this.response = node)}>
                        <p style={{color: this.state.color}}>{this.state.message}</p>
                    </div>
                    <div class="form-group">
                        <div class="col-xs-8 col-xs-offset-4">
                            <button type="submit" class="btn btn-primary btn-lg">Sign Up</button>
                        </div>  
                        <div class="col-xs-8 col-xs-offset-4">
                            <button type="button" class="btn btn-primary btn-lg" data-toggle="modal" data-target="#exampleModal" data-whatever="@mdo">Log In</button>
                        </div>  
                    </div>		      
                </form>
                <div class="modal fade" id="exampleModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel" aria-hidden="true">
                    <div class="modal-dialog" role="document">
                        <form onSubmit={(e) => {this.handleLogIn(e, auth, setAuth)}} class="form-horizontal">
                            <div class="col-xs-8 col-xs-offset-4">
                                <h2>Log In</h2>
                            </div>		
                            <div class="form-group">
                                <label class="control-label col-xs-4">Email Address</label>
                                <div class="col-xs-8">
                                    <input type="email" class="form-control" name="login_email" required="required" 
                                    ref={node => (this.loginemail = node)}/>
                                </div>        	
                            </div>
                            <div class="form-group">
                                <label class="control-label col-xs-4">Password</label>
                                <div class="col-xs-8">
                                    <input type="password" class="form-control" name="login_password" required="required" 
                                    ref={node => (this.loginpass = node)}/>
                                </div>        	
                            </div>
                            <div id ="login-response" class="text-center">
                                <p style={{color: this.state.color}}>{this.state.loginFailed}</p>
                            </div>
                            <div class="form-group">
                                <div class="col-xs-8 col-xs-offset-4">
                                    <button type="submit" class="btn btn-primary btn-sm">Log In</button>
                                </div>  
                                <div class="col-xs-8 col-xs-offset-4">
                                    <button type="button" class="btn btn-primary btn-sm" data-dismiss="modal"
                                    ref={input => this.inputElement = input}>Close</button>
                                </div>  
                            </div>	
                        </form>
                    </div>
                </div>
                </div>
            </div>
        )
    }
}
export default Home;