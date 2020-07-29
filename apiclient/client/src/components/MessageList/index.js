import React, {useEffect, useState, useContext, Component} from 'react';
import Compose from '../Compose';
import Toolbar from '../Toolbar';
import ToolbarButton from '../ToolbarButton';
import Message from '../Message';
import moment from 'moment';
import './MessageList.css';
import axios from 'axios';
import AuthContext from '../../context/AuthContext';
import propTypes from 'prop-types'
import {store} from 'react-notifications-component'
export default class MessageList extends Component {
  static contextType = AuthContext
  static prop
  constructor(props) {
    super(props);
    this.state = {
      currentChannel: this.props.currentChannel,
      currentChannelName: this.props.currentChannelName,
      currentUserID: this.props.currentUserID,
      isChannelCreator: this.props.isChannelCreator,
      isPrivateChannel: this.props.isPrivateChannel,
      messages: [],
      message: "",
      color: "red",
    }
    this.update = this.update.bind(this);
  }
  ws = new WebSocket("wss://localhost:4000/v1/websocket")
  componentDidMount() {
    const auth = this.context;
    this.ws.onopen = () => {
      // on connecting, do nothing but log it to the console
      console.log('connected')
    }

    this.ws.onmessage = event => {
      // on receiving a message, add it to the list of messages
      let data = JSON.parse(event.data).data
      console.log(data);
      if (data.channelID == this.state.currentChannel) {
        let msg = {
          author: data.creatorID,
          message: data.body,
          timestamp: data.createdAt,
          firstName: data.firstname,
          lastName: data.lastname
        }
        this.setState(previousState => ({
          messages: [...previousState.messages, msg]
        }));
      }
      if (data.creatorID != this.props.currentUserID && data.channelID != this.state.currentChannel) {
        store.addNotification({
          title: event.type,
          message: data.body,
          type: "success",
          insert: "top",
          container: "top-right",
          animationIn: ["animated", "fadeIn"],
          animationOut: ["animated", "fadeOut"],
          dismiss: {
            duration: 3000,
            onScreen: false
          }
        });
      }
    }

    this.ws.onclose = () => {
      console.log('disconnected')
      // automatically try to reconnect on connection loss
      this.setState({
        ws: new WebSocket("wss://localhost:4000/v1/websocket"),
      })
    }
    this.scrollToBottom();
  }
  componentDidUpdate() {
    this.scrollToBottom();
  }

  scrollToBottom = () => {
    this.messagesEnd.scrollIntoView({ behavior: "smooth" });
  }

  componentWillReceiveProps(newProps) {
    const {auth} = this.context
    axios.get('https://localhost:4000/v1/channels/'+ newProps.currentChannel, 
    {
      headers: {
        "Authorization": auth.token
      }
    }
    ).then(response => {
        console.log(response);
        let msgs = response.data.map(results => {
          return {
            author: results.creatorID,
            message: results.body,
            timestamp: results.createdAt,
            firstName: results.firstname,
            lastName: results.lastname
          }
        })
        this.setState({messages: msgs});
    });
    this.setState({currentChannel: newProps.currentChannel});
  }

  update = (message) => {
    const {auth} = this.context
    axios.get('https://localhost:4000/v1/channels/'+ this.state.currentChannel, 
    {
      headers: {
        "Authorization": auth.token
      }
    }
    ).then(response => {
        let msgs = response.data.map(results => {
          return {
            author: results.creatorID,
            message: results.body,
            timestamp: results.createdAt,
            firstName: results.firstname,
            lastName: results.lastname
          }
        })
        this.setState({messages: msgs});
      }).then(
        this.ws.send(message)
      );
  }
  deleteChannel = (e) => {
    const {auth} = this.context
    e.preventDefault();
    axios.delete("https://localhost:4000/v1/channels/" + this.state.currentChannel, {
        headers: {
            Authorization: auth.token
        }
    }).then(result => {
        if (result.status == 200) {
            window.location.reload(false);
        }
    }).catch(e => {
        console.log(e);
    });
  }

  editChannel = (e) => {
    const {auth} = this.context
    e.preventDefault();
    axios({
        method: 'PATCH', 
        url: "https://localhost:4000/v1/channels/" + this.state.currentChannel,
        headers: {
            Authorization: auth.token
        },
        data: {
            name: this.channelName.value,
            description: this.description.value
        }
    }).then(result => {
        if (result.status == 200) {
            window.location.reload(false);
        }
    }).catch(e => {
        console.log(e);
    })
  }

  resetMessage = () => {
    this.setState({message:"", color:""});
  }

  addUserToChannel = (e) => {
    const {auth} = this.context
    e.preventDefault();
    axios.get('https://localhost:4000/v1/users/'+ this.userName.value, 
    {
      headers: {
        "Authorization": auth.token
      }
    }
    ).then(response => {
      axios({
        method: 'LINK', 
        url: "https://localhost:4000/v1/channels/" + this.state.currentChannel,
        headers: {
            Authorization: auth.token,
            Link: response.data
        },
      }).then(result => {
          if (result.status == 200) {
              this.setState({message: result.data, color: "green"});
          }
      }).catch(e => {
        this.setState({message: e.message, color: "red"});
      })
    });
  }

  removeUserFromChannel = (e) => {
    const {auth} = this.context
    e.preventDefault();
    axios.get('https://localhost:4000/v1/users/'+ this.removeName.value, 
    {
      headers: {
        "Authorization": auth.token
      }
    }
    ).then(response => {
      axios({
        method: 'UNLINK', 
        url: "https://localhost:4000/v1/channels/" + this.state.currentChannel,
        headers: {
            Authorization: auth.token,
            Link: response.data
        },
      }).then(result => {
          console.log(result);
          if (result.status == 200) {
            this.setState({message: result.data, color: "green"});
          }
      }).catch(e => {
          this.setState({message: e.message, color: "red"});
          console.log(e);
      })
    });
  }

  renderMessages = (messages) => {
    let MY_USER_ID = this.props.currentUserID;
    let i = 0;
    let messageCount = messages.length;
    let tempMessages = [];

    while (i < messageCount) {
      let previous = messages[i - 1];
      let current = messages[i];
      let next = messages[i + 1];
      let isMine = current.author === MY_USER_ID;
      let currentMoment = moment(current.timestamp);
      let prevBySameAuthor = false;
      let nextBySameAuthor = false;
      let startsSequence = true;
      let endsSequence = true;
      let showTimestamp = true;

      if (previous) {
        let previousMoment = moment(previous.timestamp);
        let previousDuration = moment.duration(currentMoment.diff(previousMoment));
        prevBySameAuthor = previous.author === current.author;
        
        if (prevBySameAuthor && previousDuration.as('hours') < 1) {
          startsSequence = false;
        }

        if (previousDuration.as('hours') < 1) {
          showTimestamp = false;
        }
      }

      if (next) {
        let nextMoment = moment(next.timestamp);
        let nextDuration = moment.duration(nextMoment.diff(currentMoment));
        nextBySameAuthor = next.author === current.author;

        if (nextBySameAuthor && nextDuration.as('hours') < 1) {
          endsSequence = false;
        }
      }

      tempMessages.push(
        <Message
          key={i}
          isMine={isMine}
          startsSequence={startsSequence}
          endsSequence={endsSequence}
          showTimestamp={showTimestamp}
          data={current}
        />
      );

      // Proceed to the next message.
      i += 1;
    }
    return tempMessages;
  }
  render() {
    const channelExists = this.state.currentChannel != ""
    const isCreator = this.props.isChannelCreator
    const isPrivateChannel = this.props.isPrivateChannel
    return(        
      <div id="channel-header">
        <div className="message-list" >
          <Toolbar
            title={this.props.currentChannelName}
            rightItems={[
              <ToolbarButton key="info" icon="ion-ios-information-circle-outline" />,
              <ToolbarButton key="video" icon="ion-ios-videocam" />,
              <ToolbarButton key="phone" icon="ion-ios-call" />
            ]}
          />
          {
            isCreator ? 
            <div id="controls">
                {isPrivateChannel ?
                  <div>
                    <svg xmlns="http://www.w3.org/2000/svg" width="50" height="50" viewBox="0 0 8 8" data-toggle="modal" data-target="#AddToChannel">
                      <path d="M3 0v3h-3v2h3v3h2v-3h3v-2h-3v-3h-2z" />
                      <title>Add User to Channel</title>
                    </svg> 
                    <svg xmlns="http://www.w3.org/2000/svg" width="50" height="50" viewBox="0 0 8 8" data-toggle="modal" data-target="#removeFromChannel">
                      <path d="M0 0v2h8v-2h-8z" transform="translate(0 3)" />
                      <title>Remove User from Channel</title>
                    </svg>
                    </div>
                : null
                }
                <svg xmlns="http://www.w3.org/2000/svg" width="50" height="50" viewBox="0 0 8 8" data-toggle="modal" data-target="#EditChannel" data-whatever="@mdo">
                  <path d="M3.5 0l-.5 1.19c-.1.03-.19.08-.28.13l-1.19-.5-.72.72.5 1.19c-.05.1-.09.18-.13.28l-1.19.5v1l1.19.5c.04.1.08.18.13.28l-.5 1.19.72.72 1.19-.5c.09.04.18.09.28.13l.5 1.19h1l.5-1.19c.09-.04.19-.08.28-.13l1.19.5.72-.72-.5-1.19c.04-.09.09-.19.13-.28l1.19-.5v-1l-1.19-.5c-.03-.09-.08-.19-.13-.28l.5-1.19-.72-.72-1.19.5c-.09-.04-.19-.09-.28-.13l-.5-1.19h-1zm.5 2.5c.83 0 1.5.67 1.5 1.5s-.67 1.5-1.5 1.5-1.5-.67-1.5-1.5.67-1.5 1.5-1.5z"/>
                  <title>Edit Channel</title>
                </svg>
                <svg xmlns="http://www.w3.org/2000/svg" width="50" height="50" viewBox="0 0 8 8" data-toggle="modal" data-target="#deleteChannel" data-whatever="@mdo">
                  <path d="M3 0c-.55 0-1 .45-1 1h-1c-.55 0-1 .45-1 1h7c0-.55-.45-1-1-1h-1c0-.55-.45-1-1-1h-1zm-2 3v4.81c0 .11.08.19.19.19h4.63c.11 0 .19-.08.19-.19v-4.81h-1v3.5c0 .28-.22.5-.5.5s-.5-.22-.5-.5v-3.5h-1v3.5c0 .28-.22.5-.5.5s-.5-.22-.5-.5v-3.5h-1z" />            
                  <title>Delete Channel</title>
                </svg>
            </div>
           
            : null
          }
          {channelExists 
            ? <Compose currentChannel={this.state.currentChannel} 
                        update={message => {this.update(message)}} 
                        currentUserFirstName={this.props.currentUserFirstName} 
                        currentUserLastName={this.props.currentUserLastName}
                        />
            : <h1>Create or Join Channel</h1>
          }
        </div>
        <div>
        <div class="signup-form">
            <div class="modal fade" id="EditChannel" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel" aria-hidden="true">
                <div class="modal-dialog" role="document">
                    <form class="form-horizontal">
                        <div class="col-xs-8 col-xs-offset-4">
                            <h2>Edit Channel</h2>
                        </div>		
                        <div class="form-group">
                            <label class="control-label col-xs-4">Channel Name</label>
                            <div class="col-xs-8">
                                <input type="text" class="form-control" name="edit_channel_name" required="required" placeholder={this.state.firstname}
                                ref={node => (this.channelName = node)}
                                />
                            </div>        	
                        </div>
                        <div class="form-group">
                            <label class="control-label col-xs-4">Description</label>
                            <div class="col-xs-8">
                                <input type="text" class="form-control" name="edit_channel_description" placeholder={this.state.lastname}
                                ref={node => (this.description = node)}
                                />
                            </div>        	
                        </div>
                        {/* <div id ="update-response" class="text-center">
                            <p style={{color: this.state.color}}>{this.state.message}</p>
                        </div> */}
                        <div class="form-group">
                            <div class="col-xs-8 col-xs-offset-4">
                                <button class="btn btn-primary btn-sm" onClick={(e) => {this.editChannel(e)}}>Save</button>
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
        <div class="signup-form">
            <div class="modal fade" id="deleteChannel" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel" aria-hidden="true">
                <div class="modal-dialog" role="document">
                    <form class="form-horizontal">
                        <h2>Are you sure you want to delete this channel?</h2>
                        <div class="form-group">
                            <div class="col-xs-8 col-xs-offset-4">
                                <button class="btn btn-primary btn-sm" onClick={(e) => {this.deleteChannel(e)}}>YES</button>
                            </div>  
                            <div class="col-xs-8 col-xs-offset-4">
                                <button type="button" class="btn btn-primary btn-sm" data-dismiss="modal"
                                >
                                    NO</button>
                            </div>  
                        </div>	
                    </form>
                </div>
            </div>
        </div>
        <div class="signup-form">
            <div class="modal fade" id="AddToChannel" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel" aria-hidden="true">
                <div class="modal-dialog" role="document">
                    <form class="form-horizontal">
                        <div class="col-xs-8 col-xs-offset-4">
                            <h2>Add User to Channel</h2>
                        </div>		
                        <div class="form-group">
                            <label class="control-label col-xs-4">Username</label>
                            <div class="col-xs-8">
                                <input type="text" class="form-control" name="edit_channel_name" required="required" placeholder={this.state.firstname}
                                ref={node => (this.userName = node)}
                                />
                            </div>        	
                        </div>
                        <div id ="add-response" class="text-center">
                            <p style={{color: this.state.color}}>{this.state.message}</p>
                        </div>
                        <div class="form-group">
                            <div class="col-xs-8 col-xs-offset-4">
                                <button class="btn btn-primary btn-sm" onClick={(e) => {this.addUserToChannel(e)}}>Add</button>
                            </div>  
                            <div class="col-xs-8 col-xs-offset-4">
                                <button type="button" class="btn btn-primary btn-sm" data-dismiss="modal" onClick={this.resetMessage}
                                >
                                    Close</button>
                            </div>  
                        </div>	
                    </form>
                </div>
            </div>
        </div>
        <div class="signup-form">
            <div class="modal fade" id="removeFromChannel" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel" aria-hidden="true">
                <div class="modal-dialog" role="document">
                    <form class="form-horizontal">
                        <div class="col-xs-8 col-xs-offset-4">
                            <h2>Remove User From Channel</h2>
                        </div>		
                        <div class="form-group">
                            <label class="control-label col-xs-4">Username</label>
                            <div class="col-xs-8">
                                <input type="text" class="form-control" name="edit_channel_name" required="required" placeholder={this.state.firstname}
                                ref={node => (this.removeName = node)}
                                />
                            </div>        	
                        </div>
                        <div id ="remove-response" class="text-center">
                            <p style={{color: this.state.color}}>{this.state.message}</p>
                        </div>
                        <div class="form-group">
                            <div class="col-xs-8 col-xs-offset-4">
                                <button class="btn btn-primary btn-sm" onClick={(e) => {this.removeUserFromChannel(e)}}>Remove</button>
                            </div>  
                            <div class="col-xs-8 col-xs-offset-4">
                                <button type="button" class="btn btn-primary btn-sm" data-dismiss="modal" onClick={this.resetMessage}
                                >
                                    Close</button>
                            </div>  
                        </div>	
                    </form>
                </div>
            </div>
        </div>
        </div>
        <div className="message-list-container">
            {
              this.renderMessages(this.state.messages)
            }
          </div>
        <div style={{ float:"left", clear: "both" }}
          ref={(el) => { this.messagesEnd = el; }}>
        </div>
      </div>
    );
  }
}
