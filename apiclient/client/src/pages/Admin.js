import React, {useState, Component} from "react";
import {Link, Redirect} from 'react-router-dom';
import AuthContext from "../context/AuthContext";
import Axios from "axios";
import Messenger from "../components/Messenger";
import Header from "../components/Header";

export default class Admin extends Component {
    static contextType = AuthContext
    constructor(props) {
        super(props);
        this.state = {
            currentChannelName:"",
            currentChannel: "",
            currentUserID: "",
            currentUserFirstName: "",
            currentUserLastName: "",
            isChannelCreator: false
        }
    }
    componentDidMount() {
        const {auth} = this.context;
        Axios.get("https://api.will-hwang.me/v1/users/me", {
            headers: {
                Authorization: auth.token
            }
        }).then(result => {
            if (result.status == 200) {
                this.setState({currentUserID: result.data.id, currentUserFirstName: result.data.firstName, currentUserLastName: result.data.lastName});
            }
        }).catch(e => {
            console.log(e);
        })
    }
    render() {
        const { auth } = this.context
        if (!auth.token) {
            return <Redirect to="/home" />
        }
        return (
            <div>
                <Messenger currentChannel={this.state.currentChannel} 
                            currentChannelName={this.state.currentChannelName}
                            currentUserID={this.state.currentUserID}
                            isChannelCreator={this.state.isChannelCreator}
                            currentUserFirstName={this.state.currentUserFirstName}
                            currentUserLastName={this.state.currentUserLastName}/>
            </div>
        )
    }
}
