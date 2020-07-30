import React, { Component } from 'react';
import logo from '../../img/willchat.png';
import AuthContext from "../../context/AuthContext";
import Axios from "axios";
import './Header.css';
export default class Header extends Component {
    static contextType = AuthContext
    constructor(props) {
        super(props);
        this.state = {
            currentUserID: "",
            firstname: "",
            lastname: "",
            message: "",
            color: ""
        }
    }
    componentDidMount() {
        const auth = this.context;
    }
    handleLogOut = (event, auth, removeAuth) => {
        event.preventDefault();
        Axios.delete("https://api.will-hwang.me/v1/sessions/mine", {
            headers: {
                Authorization: auth.token
            }
        }).then(result => {
            if (result.status == 200) {
                const newauth = {token: null, loggedIn: false};
                removeAuth(newauth);
            }
        }).catch(e => {
            console.log(e);
        });
    }

    getUserInfo = (event, auth) => {
        event.preventDefault();
        this.first.value = "";
        this.last.value = "";
        this.setState({firstname: "", lastname: "", message: "", color: ""});
        Axios.get("https://api.will-hwang.me/v1/users/me", {
            headers: {
                Authorization: auth.token
            }
        }).then(result => {
            if (result.status == 200) {
                this.setState({firstname: result.data.firstName, lastname: result.data.lastName});
            }
        }).catch(e => {
            console.log(e);
        })
    }

    updateUserInfo = (event, auth) => {
        event.preventDefault();
        Axios({
            method: 'PATCH', 
            url: "https://api.will-hwang.me/v1/users/me",
            headers: {
                Authorization: auth.token
            },
            data: {
                firstName: this.first.value,
                lastName: this.last.value
            }
        }).then(result => {
            if (result.status == 200) {
                this.setState({message: "User Updated Successfully!", color: "green"});
            }
        }).catch(e => {
            console.log(e);
        })
    }
    render() {
        const { auth, removeAuth } = this.context
        return (
            <div id="header">
                <div class="signup-form">
                    <div id ="logo">
                        <img src={logo}></img>
                    </div>
                    <div class="modal fade" id="exampleModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel" aria-hidden="true">
                        <div class="modal-dialog" role="document">
                            <form class="form-horizontal">
                                <div class="col-xs-8 col-xs-offset-4">
                                    <h2>Edit Profile</h2>
                                </div>		
                                <div class="form-group">
                                    <label class="control-label col-xs-4">First Name</label>
                                    <div class="col-xs-8">
                                        <input type="text" class="form-control" name="edit_first" required="required" placeholder={this.state.firstname}
                                        ref={node => (this.first = node)}
                                        />
                                    </div>        	
                                </div>
                                <div class="form-group">
                                    <label class="control-label col-xs-4">Last Name</label>
                                    <div class="col-xs-8">
                                        <input type="text" class="form-control" name="edit_last" required="required" placeholder={this.state.lastname}
                                        ref={node => (this.last = node)}
                                        />
                                    </div>        	
                                </div>
                                <div id ="update-response" class="text-center">
                                    <p style={{color: this.state.color}}>{this.state.message}</p>
                                </div>
                                <div class="form-group">
                                    <div class="col-xs-8 col-xs-offset-4">
                                        <button class="btn btn-primary btn-sm" onClick={(e) => {this.updateUserInfo(e, auth)}}>Save</button>
                                    </div>  
                                    <div class="col-xs-8 col-xs-offset-4">
                                        <button type="button" class="btn btn-primary btn-sm" data-dismiss="modal"
                                        >
                                            Close</button>
                                    </div>  
                                </div>	
                            </form>
                        </div>
                    </div>
                </div>
                <div>
                    <div>
                        <svg xmlns="http://www.w3.org/2000/svg" width="50" height="50" viewBox="0 0 8 8" data-toggle="modal" data-target="#exampleModal" data-whatever="@mdo" onClick={(e) => {this.getUserInfo(e, auth)}}>
                            <path d="M4 0c-1.1 0-2 1.12-2 2.5s.9 2.5 2 2.5 2-1.12 2-2.5-.9-2.5-2-2.5zm-2.09 5c-1.06.05-1.91.92-1.91 2v1h8v-1c0-1.08-.84-1.95-1.91-2-.54.61-1.28 1-2.09 1-.81 0-1.55-.39-2.09-1z" />
                            <title>Edit Profile</title>
                        </svg>
                    </div>
                    <div>
                        <svg xmlns="http://www.w3.org/2000/svg" width="50" height="50" viewBox="0 0 8 8" onClick={(e) => {this.handleLogOut(e, auth, removeAuth)}}>
                            <path d="M3 0v4h1v-4h-1zm-1.28 1.44l-.38.31c-.81.64-1.34 1.64-1.34 2.75 0 1.93 1.57 3.5 3.5 3.5s3.5-1.57 3.5-3.5c0-1.11-.53-2.11-1.34-2.75l-.38-.31-.63.78.38.31c.58.46.97 1.17.97 1.97 0 1.39-1.11 2.5-2.5 2.5s-2.5-1.11-2.5-2.5c0-.8.36-1.51.94-1.97l.41-.31-.63-.78z"/>
                            <title>Log Out</title>
                        </svg>
                    </div>
                </div>
            </div>
        );
    }
}