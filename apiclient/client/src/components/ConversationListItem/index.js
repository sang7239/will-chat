import React, {useState, useEffect, useContext, Component} from 'react';
import shave from 'shave';

import './ConversationListItem.css';

export default class ConversationListItem extends Component {
  // useEffect(() => {
  //   shave('.conversation-snippet', 20);
  // })
  constructor(props) {
    super(props);
    this.state = {
      name: this.props.data.name,
      description: this.props.data.description,
      channelID: this.props.data.channelID,
      private: this.props.data.private
    }
  }
  render() {
    const isPrivate = this.props.data.private;
    return (
      <div>
        <div className="conversation-list-item" onClick={(e) => this.props.updateCurrentChannel(e, this.state.channelID, this.state.name)}>
          <div className="conversation-info">
            <h1 className="conversation-title">{ this.state.name }
            {
              isPrivate ?
              <svg width="1.8em" height="1.8em" viewBox="0 0 16 16" class="bi bi-lock" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
                <path fill-rule="evenodd" d="M11.5 8h-7a1 1 0 0 0-1 1v5a1 1 0 0 0 1 1h7a1 1 0 0 0 1-1V9a1 1 0 0 0-1-1zm-7-1a2 2 0 0 0-2 2v5a2 2 0 0 0 2 2h7a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-7zm0-3a3.5 3.5 0 1 1 7 0v3h-1V4a2.5 2.5 0 0 0-5 0v3h-1V4z"/>
              </svg> : null
            }
            </h1>
            <p className="conversation-snippet">{ this.state.description }
            </p>
          </div>
        </div>
      </div>
    );
  }
}