import React, {useState, useEffect, useContext, useRef, Component} from 'react';
import ConversationSearch from '../ConversationSearch';
import ConversationListItem from '../ConversationListItem';
import MessageList from '../MessageList';

import Toolbar from '../Toolbar';
import ToolbarButton from '../ToolbarButton';
import axios from 'axios';
import AuthContext from "../../context/AuthContext"
import './ConversationList.css';
import Axios from 'axios';

export default class ConversationList extends Component {
    static contextType = AuthContext
    constructor(props) {
      super(props);
      this.state = {
        conversations: [],
        currentChannel: this.props.currentChannel,
        currentChannelName: this.props.currentChannelName,
      }
      this.updateCurrentChannel = this.updateCurrentChannel.bind(this);
    }

    componentDidMount() {
      const {auth} = this.context
      axios.get('https://api.will-hwang.me/v1/channels', 
      {
        headers: {
          "Authorization": auth.token
        }
      }
      ).then(response => {
          let channels = response.data.map(results => {
            return {
              name: results.name,
              description: results.description,
              channelID: results.id,
              private: results.private
            }
          })
          this.setState({conversations: channels});
      });
    }

  postChannels = (e) => {
    const {auth} = this.context
    e.preventDefault();
    Axios({
      method: 'POST', 
      url: "https://api.will-hwang.me/v1/channels",
      headers: {
          Authorization: auth.token
      },
      data: {
          name: this.name.value,
          description: this.description.value,
          private: this.private.checked
      }
    }).then(result => {
        if (result.status == 200) {
          window.location.reload(false);
        }
    }).catch(e => {
        console.log(e);
    })
   } 
    updateCurrentChannel= (e, channelID, channelName) => {
      e.preventDefault();
      this.props.updateChannel(channelID, channelName);
    }
   render() {
    return (
      <div className="conversation-list">
        <Toolbar
          title="Channels"
          leftItems={[
            <ToolbarButton key="cog" icon="ion-ios-cog" />
          ]}
          rightItems={[
            // <ToolbarButton key="add" icon="ion-ios-add-circle-outline" />
            <svg xmlns="http://www.w3.org/2000/svg" width="30" height="30" viewBox="0 0 8 8"data-toggle="modal" data-target="#channelModal" data-whatever="@mdo" >
              <path d="M3 0v3h-3v2h3v3h2v-3h3v-2h-3v-3h-2z" />
              <title>Add Channel</title>
            </svg>
          ]}
        />
        <ConversationSearch />
        {
          this.state.conversations.map(conversation =>
            <ConversationListItem
              key={conversation.name}
              data={conversation}
              updateCurrentChannel={this.updateCurrentChannel}
            />
          )
        }
        <div class="signup-form">
          <div class="modal fade" id="channelModal" tabindex="-1" role="dialog" aria-labelledby="exampleModalLabel" aria-hidden="true">
              <div class="modal-dialog" role="document">
                  <form class="form-horizontal">
                      <div class="col-xs-8 col-xs-offset-4">
                          <h2>Create Channel</h2>
                      </div>		
                      <div class="form-group">
                          <label class="control-label col-xs-4">Channel Name</label>
                          <div class="col-xs-8">
                              <input type="text" class="form-control" name="channel_name" required="required" ref={node => (this.name = node)}/>
                          </div>        	
                      </div>
                      <div class="form-group">
                          <label class="control-label col-xs-4">Description</label>
                          <div class="col-xs-8">
                              <input type="text" class="form-control" name="channel_description" ref={node => (this.description = node)}
                              />
                          </div>        	
                      </div>
                      <div class="form-group">
                          <label class="control-label col-xs-4" for="defaultChecked">Public</label>
                          <div class="col-xs-8">
                            <input type="radio" class="custom-control-input" id="defaultChecked" value = "public" name="defaultExampleRadios"/>
                          </div>
                      </div>
                      <div class="form-group">
                          <label class="control-label col-xs-4" for="defaultUnchecked">Private</label>
                          <div class="col-xs-8">
                            <input type="radio" class="custom-control-input" id="defaultUnchecked" value = "private" name="defaultExampleRadios" ref={node => (this.private = node)}/>
                          </div>
                      </div>
                      <div id ="update-response" class="text-center">
                          {/* <p style={{color: this.state.color}}>{this.state.message}</p> */}
                      </div>
                      <div class="form-group">
                          <div class="col-xs-8 col-xs-offset-4">
                              <button class="btn btn-primary btn-sm" onClick={e => this.postChannels(e)}>Create</button>
                          </div>  
                          <div class="col-xs-8 col-xs-offset-4">
                          <button type="button" class="btn btn-primary btn-sm" data-dismiss="modal"
                          >
                              Cancel</button>
                      </div> 
                      </div>	
                  </form>
              </div>
          </div>
        </div>
      </div>
    );
   }
}