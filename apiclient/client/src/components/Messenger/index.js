import React, { Component } from 'react';
import ConversationList from '../ConversationList';
import MessageList from '../MessageList';
import Header from '../Header';
import AuthContext from "../../context/AuthContext";
import './Messenger.css';
import axios from "axios";

export default class Messenger extends Component{
    static contextType = AuthContext
    constructor(props) {
      super(props);
      this.state = {
        currentChannelName: this.props.currentChannelName,
        currentChannel: this.props.currentChannel,
        currentUserID: this.props.currentUserID,
        isChannelCreator: this.props.isChannelCreator,
        isPrivateChannel: false
      }
    }
    componentWillMount() {
      const auth = this.context
    }
    updateChannel = (channel, channelName) => {
      const {auth} = this.context;
      this.setState({currentChannel: channel, currentChannelName: channelName});
      axios.get("https://localhost:4000/v1/channelInfo/" + channel, {
            headers: {
                Authorization: auth.token
            }
        }).then(result => {
          console.log(result.data.creatorID);
            if (result.status == 200) {
              if (result.data.creatorID == this.props.currentUserID) {
                this.setState({isChannelCreator: true});
              } else {
                this.setState({isChannelCreator: false});
              }
              this.setState({isPrivateChannel: result.data.private});
            } 
        }).catch(e => {
            console.log(e);
        })
    }

    render() {
      return (
        <div className="messenger">
          {/* <Toolbar
            title="Messenger"
            leftItems={[
              <ToolbarButton key="cog" icon="ion-ios-cog" />
            ]}
            rightItems={[
              <ToolbarButton key="add" icon="ion-ios-add-circle-outline" />
            ]}
          /> */}
  
          {/* <Toolbar
            title="Conversation Title"
            rightItems={[
              <ToolbarButton key="info" icon="ion-ios-information-circle-outline" />,
              <ToolbarButton key="video" icon="ion-ios-videocam" />,
              <ToolbarButton key="phone" icon="ion-ios-call" />
            ]}
          /> */}
          <div className="scrollable sidebar">
            <ConversationList currentChannel={this.state.currentChannel} 
                              updateChannel={this.updateChannel}
                              updateChannelCreator={this.updateChannelCreator}/>
          </div>
  
          <div className="scrollable content">
            <Header updateUser={this.updateUser}/>
            <MessageList currentChannel={this.state.currentChannel} 
                          currentChannelName={this.state.currentChannelName}
                          currentUserID={this.props.currentUserID}
                          isChannelCreator={this.state.isChannelCreator}
                          isPrivateChannel={this.state.isPrivateChannel}
                          currentUserFirstName={this.props.currentUserFirstName}
                          currentUserLastName={this.props.currentUserLastName}/>
          </div>
        </div>
      );
    }
}