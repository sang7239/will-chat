import React, {useState, useRef, useContext, Component} from 'react';
import './Compose.css';
import Axios from 'axios';
import AuthContext from '../../context/AuthContext';

export default class Compose extends Component {
  static contextType = AuthContext
  constructor(props) {
    super(props);
  }
  componentDidMount() {
    const auth = this.context;
  }
  handleSubmit = (e) => {
    e.preventDefault();
    const {auth} = this.context;
    if (e.keyCode == 13) {
      Axios({
        method: 'POST', 
        url: "https://api.will-hwang.me/v1/messages",
        headers: {
            Authorization: auth.token
        },
        data: {
            channelId: this.props.currentChannel,
            body: this.message.value,
            firstName: this.props.currentUserFirstName,
            lastName: this.props.currentUserLastName
        }
      }).then(result => {
          if (result.status == 200) {
            this.props.update(JSON.stringify(result.data));
          }
      }).catch(e => {
          console.log(e);
      })
      this.message.value = ""   
     } 
  }
  render() {
    return (
      <div className="compose">
          <input
            type="text"
            ref={node => (this.message = node)}
            className="compose-input"
            placeholder="Type a message"
            onKeyUp={e => this.handleSubmit(e)}
          />
      </div>
    );
  }
}