import React from 'react';
import { Button, Card } from 'react-hephaestus';

const ParticipationStatusLoadingPage = () => (
  <Card className="participation-status-page">
    <h1 className="title">Status Peserta</h1>
    <div className="participation-status-list">
      <div className="participation-status">
        <div className="user"><span className="skeleton">Nama peserta</span></div>
        <div className="ip-address"><span className="skeleton">192.168.0.1</span></div>
        <div className="login-at"><span className="skeleton">abc</span></div>
        <div className="lock-button"><Button className="skeleton">Lock</Button></div>
      </div>
      <div className="participation-status">
        <div className="user"><span className="skeleton">Nama peserta</span></div>
        <div className="ip-address"><span className="skeleton">192.168.0.1</span></div>
        <div className="login-at"><span className="skeleton">abc</span></div>
        <div className="lock-button"><Button className="skeleton">Lock</Button></div>
      </div>
      <div className="participation-status">
        <div className="user"><span className="skeleton">Nama peserta</span></div>
        <div className="ip-address"><span className="skeleton">192.168.0.1</span></div>
        <div className="login-at"><span className="skeleton">abc</span></div>
        <div className="lock-button"><Button className="skeleton">Lock</Button></div>
      </div>
      <div className="participation-status">
        <div className="user"><span className="skeleton">Nama peserta</span></div>
        <div className="ip-address"><span className="skeleton">192.168.0.1</span></div>
        <div className="login-at"><span className="skeleton">abc</span></div>
        <div className="lock-button"><Button className="skeleton">Lock</Button></div>
      </div>
      <div className="participation-status">
        <div className="user"><span className="skeleton">Nama peserta</span></div>
        <div className="ip-address"><span className="skeleton">192.168.0.1</span></div>
        <div className="login-at"><span className="skeleton">abc</span></div>
        <div className="lock-button"><Button className="skeleton">Lock</Button></div>
      </div>
    </div>
  </Card>
);

export default ParticipationStatusLoadingPage;
