import React, { useEffect } from 'react';
import { Card } from 'react-hephaestus';
import { connect } from 'react-redux';
import { Link } from 'react-router-dom';

import conf from '../conf';
import { User, USER_ROLE } from '../modules/charon/auth/api';
import * as sessionSelectors from '../modules/session/selector';
import { AppState } from '../modules/store';
import { ROUTE_EVENT_LIST } from './routes';
import './HomePage.scss';

export interface ConnectedHomePageProps {
  user: User | null;
};

const UserListMenu = () => (
  <Link to={ROUTE_EVENT_LIST}>
    <Card cardType="outlined" className="menu">
      <div><i className="fas fa-user-friends" /> Daftar peserta</div>
      <i className="fas fa-chevron-right" />
    </Card>
  </Link>
);

const EventListMenu = () => (
  <Link to={ROUTE_EVENT_LIST}>
    <Card cardType="outlined" className="menu">
      <div><i className="fas fa-book" /> Daftar ujian</div>
      <i className="fas fa-chevron-right" />
    </Card>
  </Link>
);

const HomePage = (props: ConnectedHomePageProps) => {
  const { user } = props;
  useEffect(() => { document.title = conf.appName; }, []);

  return (
    <div className="home-page">
      <Card>
        {!!user
          ? (
            <>
              <h1>Selamat datang, {user.name}!</h1>
              <div className="menus">
                <EventListMenu />
                {user.role === USER_ROLE.ADMIN && <UserListMenu />}
              </div>
            </>)
          : (<>
              <h1>{conf.appName} </h1>
              <p>{conf.appName} adalah software untuk ujian semi-online.</p>
            </>)}
      </Card>
    </div>
  );
};

const mapStateToProps = (state: AppState) => ({
  user: sessionSelectors.getUser(state),
});

export default connect(mapStateToProps)(HomePage);
