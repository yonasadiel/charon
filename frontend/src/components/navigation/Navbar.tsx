import React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import logo from '../../assets/logo.svg';
import profileIcon from '../../assets/profile.svg';
import conf from '../../conf';
import { User } from '../../modules/charon/auth/api';
import * as charonAuthActions from '../../modules/charon/auth/action';
import * as sessionSelectors from '../../modules/session/selector';
import { AppState } from '../../modules/store';
import { ROUTE_HOME, ROUTE_LOGIN } from '../../pages/routes';
import './Navbar.scss';

export interface ConnectedNavbarProps {
  user: User | null;
  logoutAction: () => void,
}

const Navbar = (props: ConnectedNavbarProps) => {
  const { logoutAction, user } = props;
  const handleLogoutClick = () => {
    logoutAction();
  };
  return (
    <div className="navbar">
      <Link to={ROUTE_HOME} className="app">
        <img className="logo" src={logo} alt="charon-logo"/>
        <h1 className="title">{conf.appName}</h1>
      </Link>
      <div className="divider" />
      {!!user ? (
        <div className="user" onClick={handleLogoutClick}>
          <img src={profileIcon} alt="profile-icon" />
          <span>{user.name}</span>
        </div>
      ) : (
        <Link to={ROUTE_LOGIN}>
          <div className="user">Login</div>
        </Link>
      )}
    </div>
  );
}

const mapStateToProps = (state: AppState) => ({
  user: sessionSelectors.getUser(state),
});

const mapDispatchToProps = {
  logoutAction: charonAuthActions.logout,
};

export default connect(mapStateToProps, mapDispatchToProps)(Navbar);
