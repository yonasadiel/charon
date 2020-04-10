import React from 'react';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import logo from '../../assets/logo.svg';
import profileIcon from '../../assets/profile.svg';
import conf from '../../conf';
import { User } from '../../modules/charon/auth/api';
import * as sessionSelectors from '../../modules/session/selector';
import { AppState } from '../../modules/store';
import { ROUTE_HOME, ROUTE_LOGIN } from '../../pages/routes';
import './Navbar.scss';

export interface NavbarProps {
  user: User | null;
}

const Navbar = (props: NavbarProps) => {
  const { user } = props;
  return (
    <div className="navbar">
      <Link to={ROUTE_HOME} className="app">
        <img className="logo" src={logo} alt="charon-logo"/>
        <h1 className="title">{conf.appName}</h1>
      </Link>
      <div className="divider" />
      <Link to={ROUTE_LOGIN} className="user">
        {!!user && (<img src={profileIcon} alt="profile-icon" />)}
        <span>{!!user ? user.name : "Login"}</span>
      </Link>
    </div>
  );
}

const mapStateToProps = (state: AppState) => ({
  user: sessionSelectors.getUser(state),
});

export default connect(mapStateToProps)(Navbar);
